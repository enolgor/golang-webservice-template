package models

import (
	"database/sql/driver"
	"errors"
	"slices"
	"strings"
)

type Role string

type Roles []Role

const (
	USER     Role = "user"
	INTERNAL Role = "internal"
	ADMIN    Role = "admin"
)

var AllRoles []Role = []Role{
	USER, INTERNAL, ADMIN,
}

func (roles Roles) AnyOf(allowed ...Role) bool {
	return slices.ContainsFunc(roles, func(r Role) bool {
		return slices.Contains(allowed, r)
	})
}

func (roles Roles) Value() (driver.Value, error) {
	strRoles := make([]string, len(roles))
	for i, role := range roles {
		if !slices.Contains(AllRoles, role) {
			return nil, errors.New("unsupported role value")
		}
		strRoles[i] = string(role)
	}
	return strings.Join(strRoles, ","), nil
}

func (roles *Roles) Scan(v any) error {
	if value, ok := v.(string); !ok {
		return errors.New("unable to scan value")
	} else {
		strRoles := strings.Split(value, ",")
		*roles = []Role{}
		for _, strRole := range strRoles {
			role := Role(strRole)
			if !slices.Contains(AllRoles, role) {
				return errors.New("unsupported role value")
			}
			(*roles) = append((*roles), role)
		}
	}
	return nil
}
