package interceptor

import (
	"encoding/json"
	"fmt"

	"github.com/kostyaBroTutor/auth_interceptor/pkg/slice"
	"github.com/kostyaBroTutor/auth_interceptor/proto"
)

type TokenInfo struct {
	UserID      string
	Roles       []proto.Role
	Permissions []proto.Permission
}

type tokenInfo struct {
	UserID      string   `json:"user_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func (t TokenInfo) String() string {
	ti := tokenInfo{
		UserID:      t.UserID,
		Roles:       make([]string, len(t.Roles)),
		Permissions: make([]string, len(t.Permissions)),
	}

	for i, role := range t.Roles {
		ti.Roles[i] = role.String()
	}

	for i, permission := range t.Permissions {
		ti.Permissions[i] = permission.String()
	}

	j, err := json.Marshal(ti)
	if err != nil {
		return fmt.Sprintf("can't marshal token info: %s", err.Error())
	}

	return string(j)
}

// HasAccessToMethodWithWithAuthOptions checks if the user has access to the method.
// If options have required roles, then the user must have at least one of them.
// If options have required permissions, then the user must have all of them.
func (t TokenInfo) HasAccessToMethodWithWithAuthOptions(
	options *proto.MethodAuthOptions,
) bool {
	if options == nil {
		return true
	}

	if len(options.GetRequiredRoles()) > 0 {
		if !t.hasAtLeastOneRole(options.GetRequiredRoles()...) {
			return false
		}
	}

	if len(options.GetRequiredPermissions()) > 0 {
		if !t.hasAllPermissions(options.GetRequiredPermissions()...) {
			return false
		}
	}

	return true
}

func (t TokenInfo) hasAtLeastOneRole(requiredRoles ...proto.Role) bool {
	tokenRolesSet := slice.ToSet(t.Roles)

	for _, requiredRole := range requiredRoles {
		if _, hasRole := tokenRolesSet[requiredRole]; hasRole {
			return true
		}
	}

	return false
}

func (t TokenInfo) hasAllPermissions(requiredPermissions ...proto.Permission) bool {
	tokenPermissionsSet := slice.ToSet(t.Permissions)

	for _, requiredPermission := range requiredPermissions {
		if _, hasPermission := tokenPermissionsSet[requiredPermission]; !hasPermission {
			return false
		}
	}

	return true
}
