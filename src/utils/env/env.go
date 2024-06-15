package env

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"

	"github.com/joho/godotenv"
)

type loadFunc func() (string, string, error)

var (
	errNotPresent   error = errors.New("not present")
	errNotParseable error = errors.New("can't be parsed")
	errInvalid      error = errors.New("invalid value")
)

var Bool func(string) (bool, error) = strconv.ParseBool
var String func(string) (string, error) = func(s string) (string, error) { return s, nil }

func ParseEnv(log *slog.Logger, level *slog.LevelVar, loadEnvFuncs ...loadFunc) {
	if err := godotenv.Load(); err != nil {
		log.Error("failed to load environment file")
		os.Exit(-1)
	}
	setLogLevel(level)
	log.Info("loaded environment file")
	isErr := false
	for _, f := range loadEnvFuncs {
		if key, strvalue, err := f(); err != nil {
			log.Error("failed to load environment variable", "key", key)
			log.Debug("failed to load environment variable", "key", key, "value", strvalue, "error", err.Error())
			isErr = true
		} else {
			log.Debug("loaded environment variable", "key", key, "value", strvalue)
		}
	}
	if isErr {
		log.Error("failed to load environment variables, exiting")
		os.Exit(-1)
	}
	log.Info("environment variables loaded successfully")
}

func getEnvOrElseErr(key string) (string, error) {
	if value, ok := os.LookupEnv(key); !ok {
		return "", errNotPresent
	} else {
		return value, nil
	}
}

func Get[T comparable](key string, target *T, strToType func(string) (T, error), defaultValue T) loadFunc {
	return func() (string, string, error) {
		strvalue := os.Getenv(key)
		if value, err := strToType(os.Getenv(key)); err != nil {
			*target = defaultValue
		} else {
			*target = value
		}
		return key, strvalue, nil
	}
}

func MustGet[T comparable](key string, target *T, strToType func(string) (T, error)) loadFunc {
	return func() (string, string, error) {
		if strvalue, err := getEnvOrElseErr(key); err != nil {
			return key, strvalue, err
		} else if value, err := strToType(strvalue); err != nil {
			return key, strvalue, errNotParseable
		} else {
			*target = value
			return key, strvalue, nil
		}
	}
}

func MustGetValid[T comparable](key string, target *T, strToType func(string) (T, error), validValues ...T) loadFunc {
	return func() (string, string, error) {
		if strvalue, err := getEnvOrElseErr(key); err != nil {
			return key, strvalue, err
		} else if value, err := strToType(strvalue); err != nil {
			return key, strvalue, errNotParseable
		} else if !slices.Contains(validValues, value) {
			return key, strvalue, errInvalid
		} else {
			*target = value
			return key, strvalue, nil
		}
	}
}

func setLogLevel(level *slog.LevelVar) {
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "DEBUG":
		level.Set(slog.LevelDebug)
	case "INFO":
		level.Set(slog.LevelInfo)
	case "WARN":
		level.Set(slog.LevelWarn)
	case "ERROR":
		level.Set(slog.LevelError)
	case "":
		level.Set(slog.LevelInfo)
	default:
		panic(fmt.Errorf("Unrecognized log level: " + logLevel))
	}
}
