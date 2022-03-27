package middleware

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/db"
	"github.com/labstack/echo/v4"
)

type RbacConfig struct {
	Users db.UserStore
	Auth  auth.Auth
}

type RBAC struct {
	config RbacConfig
}

func NewRBAC(config RbacConfig) RBAC {
	return RBAC{config}
}

// RBAC Middleware
func (r *RBAC) Middleware(resource string, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId := r.config.Auth.GetUserId(c)
			user, err := r.config.Users.Find(userId)
			if err != nil {
				return echo.ErrUnauthorized
			}
			permissions, err := r.config.Users.Permissions(&user)
			if err != nil {
				return err
			}
			for _, p := range permissions {
				matchResource, matchAction := false, false
				if p.Resource == resource || p.Resource == "*" {
					matchResource = true
				}
				if p.Action == action || p.Action == "*" {
					matchAction = true
				}
				if matchResource && matchAction {
					return next(c)
				}
			}
			return echo.ErrForbidden
		}
	}
}

// Alias for (*RBAC).Middleware
func (r *RBAC) M(resource string, action string) echo.MiddlewareFunc {
	return r.Middleware(resource, action)
}
