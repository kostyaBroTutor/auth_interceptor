package interceptor

import (
	"github.com/kostyaBroTutor/auth_interceptor/pkg/slice"
	"github.com/kostyaBroTutor/auth_interceptor/proto"
)

type TokenInfo struct {
	UserID      string
	Roles       []proto.Role
	Permissions []proto.Permission
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
