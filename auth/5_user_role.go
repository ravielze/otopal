package auth

import (
	"strings"
)

type Role int16

const (
	ROLE_ADMIN Role = iota + 1
	ROLE_CUSTOMER
	ROLE_TECHNICIAN
	//insert new role here
)

// lowercase role name => id (iota+1)
var roles = map[string]Role{
	"admin":      1,
	"customer":   2,
	"technician": 3,
	//insert new role here
}

func (r Role) IsExist() bool {
	//change 2 to the last number
	return int16(r) >= 1 && int16(r) <= 3
}

func (r Role) String() string {
	//add role name here
	return [...]string{"Admin", "Customer", "Technician"}[r-1]
}

func GetRole(role string) Role {
	result := roles[strings.ToLower(role)]
	if int16(result) == 0 {
		result = ROLE_CUSTOMER
	}
	return result
}

func (r Role) Equal(other Role) bool {
	return (int16(r) == int16(other))
}

func (r Role) IsRestricted() bool {
	return r.Equal(ROLE_ADMIN)
}
