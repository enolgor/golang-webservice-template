package timezone

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

var TimeZones []string = []string{}

const zoneFile string = "/usr/share/zoneinfo/zone1970.tab"

func init() {
	file, err := os.Open(zoneFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) >= 3 {
			TimeZones = append(TimeZones, fields[2])
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	slices.Sort(TimeZones)
}
