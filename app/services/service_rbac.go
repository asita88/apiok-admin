package services

import (
	"apiok-admin/app/models"
	"apiok-admin/app/utils"
	"strings"
)

const (
	RoleAdmin    = "admin"
	RoleViewer   = "viewer"
	RoleOperator = "operator"
	permWildcard = "*"
)

func PermissionsFromRole(role string) []string {
	switch role {
	case RoleAdmin, "":
		return []string{permWildcard}
	case RoleViewer:
		return []string{"dashboard", "service", "router", "upstream", "ssl", "global-plugin"}
	case RoleOperator:
		out := PermissionsFromRole(RoleViewer)
		out = append(out, "log", "access-log")
		return out
	default:
		return []string{permWildcard}
	}
}

func RolesFromRole(role string) []string {
	r := role
	if r == "" {
		r = RoleAdmin
	}
	return []string{r}
}

func UserByTokenSubject(subject string) models.Users {
	m := models.Users{}
	u := m.UserInfoByName(subject)
	if u.ResID != "" {
		return u
	}
	u = m.UserInfoByEmail(subject)
	if u.ResID != "" {
		return u
	}
	return models.Users{Role: RoleAdmin}
}

func PermissionsForSubject(subject string) []string {
	u := UserByTokenSubject(subject)
	return PermissionsFromRole(u.Role)
}

func hasPermission(granted []string, code string) bool {
	for _, g := range granted {
		if g == permWildcard {
			return true
		}
		if g == code {
			return true
		}
	}
	return false
}

func hasAnyPermission(granted []string, codes []string) bool {
	for _, c := range codes {
		if hasPermission(granted, c) {
			return true
		}
	}
	return false
}

type rbacRule struct {
	prefix string
	anyOf  []string
}

var rbacRules = []rbacRule{
	{"/admin/service/plugin/config", []string{"service"}},
	{"/admin/router/plugin/config", []string{"router"}},
	{"/admin/global/plugin/config", []string{"global-plugin"}},
	{"/admin/plugin/", []string{"service", "router", "global-plugin"}},
	{"/admin/service/", []string{"service"}},
	{"/admin/router/", []string{"router"}},
	{"/admin/upstream/", []string{"upstream"}},
	{"/admin/certificate/", []string{"ssl"}},
	{"/admin/letsencrypt/", []string{"ssl"}},
	{"/admin/cluster-node/", []string{"upstream"}},
	{"/admin/log/access/aggregation", []string{"dashboard", "access-log"}},
	{"/admin/log/access/", []string{"access-log"}},
	{"/admin/log/list", []string{"log"}},
	{"/admin/user/list", []string{"user"}},
	{"/admin/user/info/", []string{"user"}},
	{"/admin/user/add", []string{"user"}},
	{"/admin/user/update/", []string{"user"}},
	{"/admin/user/delete/", []string{"user"}},
}

func pathAlwaysAllowed(path string) bool {
	return strings.HasPrefix(path, "/admin/user/logout") ||
		strings.HasPrefix(path, "/admin/user/change-password")
}

func RequiredPermissionsForPath(path string) (skipRBAC bool, anyOf []string) {
	if pathAlwaysAllowed(path) {
		return true, nil
	}
	for _, rule := range rbacRules {
		if strings.HasPrefix(path, rule.prefix) {
			return false, rule.anyOf
		}
	}
	return false, nil
}

func RBACAllowed(token string, path string) bool {
	if pathAlwaysAllowed(path) {
		return true
	}
	subject, err := utils.ParseToken(token)
	if err != nil {
		return false
	}
	granted := PermissionsForSubject(subject)
	skip, required := RequiredPermissionsForPath(path)
	if skip {
		return true
	}
	if len(required) == 0 {
		return true
	}
	return hasAnyPermission(granted, required)
}
