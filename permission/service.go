package permission

import (
	"authgoblue/ctx"
)

type Service struct {
	context *ctx.Service
}

func NewService(
	contextService *ctx.Service,
) *Service {

	return &Service{
		context: contextService,
	}
}
