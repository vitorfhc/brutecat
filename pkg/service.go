package brutecat

import (
	"context"
	"fmt"
)

type Authenticator interface {
	fmt.Stringer
	Authenticate(context.Context, Credentials) (bool, error)
	TestConnection(context.Context) error
}
