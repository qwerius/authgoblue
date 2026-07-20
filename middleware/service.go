package middleware

import (
	"authgoblue/ctx"
	"authgoblue/permission"
	"authgoblue/revoke"
	"authgoblue/role"
	"authgoblue/session"
	"authgoblue/token"
)

type Service struct {
	token      *token.Service
	context    *ctx.Service
	role       *role.Service
	permission *permission.Service
	revoke     *revoke.Service
	session    *session.Service

	header     string
	prefix     string
	cookie     bool
	cookieName string
}

func NewService(
	tokenService *token.Service,
	contextService *ctx.Service,
	roleService *role.Service,
	permissionService *permission.Service,
	revokeService *revoke.Service,
	sessionService *session.Service,

	header string,
	prefix string,
	cookie bool,
	cookieName string,
) *Service {

	return &Service{

		token: tokenService,

		context: contextService,

		role: roleService,

		permission: permissionService,

		revoke: revokeService,

		session: sessionService,

		header: header,

		prefix: prefix,

		cookie: cookie,

		cookieName: cookieName,
	}
}
