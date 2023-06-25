package authenticators

import (
	"context"
	"fmt"
	"strings"

	brutecat "github.com/vitorfhc/brutecat/pkg"
	"github.com/vitorfhc/brutecat/pkg/utils"
	"github.com/vitorfhc/ftp"
)

type FTP struct {
	Host string
	Port uint16
}

func (a *FTP) Authenticate(ctx context.Context, creds brutecat.Credentials) (bool, error) {
	addr := fmt.Sprintf("%s:%d", a.Host, a.Port)
	var conn *ftp.ServerConn

	// We have to do this because the FTP library is not handling the context
	// properly. It's probably not watching for the context to be canceled.
	err := utils.RunWithContext(ctx, func() error {
		var err error
		conn, err = ftp.Dial(addr)
		return err
	})
	if err != nil {
		return false, err
	}
	defer conn.Quit()

	err = utils.RunWithContext(ctx, func() error {
		return conn.Login(creds.Username, creds.Password, false)
	})
	if err != nil && strings.Contains(err.Error(), "530") {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (a *FTP) String() string {
	return "FTP"
}

func (a *FTP) TestConnection(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", a.Host, a.Port)

	options := ftp.DialWithContext(ctx)
	conn, err := ftp.Dial(addr, options)
	if err != nil {
		return err
	}
	defer conn.Quit()

	return nil
}
