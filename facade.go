package authgoblue

import (
	"context"

	"github.com/qwerius/authgoblue/login"
)

// SignIn adalah facade sederhana untuk proses login.
//
// Developer cukup memakai:
// agb.SignIn(...)
//
// tanpa harus mengakses:
// agb.Login.Login(...)
func (a *AuthGoBlue) SignIn(
	ctx context.Context,
	req login.Request,
) (*login.Result, error) {

	return a.Login.Login(
		ctx,
		req,
	)
}
