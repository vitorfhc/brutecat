package authenticators

import (
	"context"
	"fmt"
	"strings"

	"github.com/jlaffaye/ftp"
	brutecat "github.com/vitorfhc/brutecat/pkg"
)

type FTP struct {
	Host string
	Port uint16
}

func (a *FTP) Authenticate(ctx context.Context, creds brutecat.Credentials) (bool, error) {
	addr := fmt.Sprintf("%s:%d", a.Host, a.Port)

	options := ftp.DialWithContext(ctx)
	conn, err := ftp.Dial(addr, options)
	if err != nil {
		return false, err
	}
	defer conn.Quit()

	err = conn.Login(creds.Username, creds.Password)
	if err != nil && strings.Contains(err.Error(), "530") {
		return false, nil
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
