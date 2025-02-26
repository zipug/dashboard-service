package models

import (
	"strings"
)

type PermissionAction string

const (
	Create PermissionAction = "create"
	Read   PermissionAction = "read"
	Update PermissionAction = "update"
	Delete PermissionAction = "delete"
)

type Permission struct {
	Name   string
	Action PermissionAction
}

type PermissionData struct {
	Id          int64
	Name        string
	Description string
	Create      bool
	Read        bool
	Update      bool
	Delete      bool
}

func ToPermission(perm string) Permission {
	perm_and_action := strings.Split(perm, ":")
	p := Permission{Name: perm_and_action[0], Action: PermissionAction(perm_and_action[1])}
	return p
}

type Role struct {
	Id          int64
	Name        string
	Description string
	Permissions []PermissionData
	IsCustom    bool
}
