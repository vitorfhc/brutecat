package brutecat

import (
	"context"
	"fmt"
	"sync"

	"github.com/vitorfhc/brutecat/pkg/utils"
)

type Engine struct {
	// Concurrent execution configuration
	threads uint16

	// Brute force configuration
	authenticator Authenticator
	users         []string
	passwords     []string

	// Runner's channels
	credentialsToTry chan Credentials

	// Runner's callbacks
	OnErrorCallback   func(error)
	OnSuccessCallback func(Credentials)
	OnFailCallback    func(Credentials)
}

type Credentials struct {
	Username string
	Password string
}

func NewEngine(users, passwords []string, threads uint16, authenticator Authenticator) (*Engine, error) {
	if len(users) == 0 {
		errMsg := "users to brute force cannot be empty"
		return nil, fmt.Errorf(errMsg)
	}

	if len(passwords) == 0 {
		errMsg := "passwords to brute force cannot be empty"
		return nil, fmt.Errorf(errMsg)
	}

	if threads == 0 {
		errMsg := "threads must be greater than 0"
		return nil, fmt.Errorf(errMsg)
	}

	if authenticator == nil {
		errMsg := "authenticator cannot be nil"
		return nil, fmt.Errorf(errMsg)
	}

	return &Engine{
		threads:           threads,
		authenticator:     authenticator,
		users:             users,
		passwords:         passwords,
		credentialsToTry:  make(chan Credentials, threads*2),
		OnErrorCallback:   func(err error) {},
		OnSuccessCallback: func(creds Credentials) {},
		OnFailCallback:    func(creds Credentials) {},
	}, nil
}

func NewEngineWithFiles(usersFile, passwordsFile string, threads uint16, authenticator Authenticator) (*Engine, error) {
	users, err := utils.FileLinesToSlice(usersFile)
	if err != nil {
		return nil, err
	}

	passwords, err := utils.FileLinesToSlice(passwordsFile)
	if err != nil {
		return nil, err
	}

	return NewEngine(users, passwords, threads, authenticator)
}

func (b *Engine) credsWorker(ctx context.Context) {
	usersIndex := 0
	passwordsIndex := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if passwordsIndex >= len(b.passwords) {
				passwordsIndex = 0
				usersIndex++
			}

			if usersIndex >= len(b.users) {
				close(b.credentialsToTry)
				return
			}

			creds := Credentials{
				Username: b.users[usersIndex],
				Password: b.passwords[passwordsIndex],
			}
			b.credentialsToTry <- creds
			passwordsIndex++
		}
	}
}

func (b *Engine) worker(ctx context.Context) {
	for {
		select {
		case creds, ok := <-b.credentialsToTry:
			if !ok {
				return
			}

			valid, err := b.authenticator.Authenticate(ctx, creds)
			if err != nil {
				b.OnErrorCallback(err)
			} else if valid {
				b.OnSuccessCallback(creds)
			} else {
				b.OnFailCallback(creds)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *Engine) Run(ctx context.Context) {
	err := b.authenticator.TestConnection(ctx)
	if err != nil {
		b.OnErrorCallback(err)
		return
	}

	wgWorkers := sync.WaitGroup{}

	wgWorkers.Add(1)
	go func() {
		defer wgWorkers.Done()
		b.credsWorker(ctx)
	}()

	for i := uint16(1); i <= b.threads; i++ {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			b.worker(ctx)
		}()
	}

	wgWorkers.Wait()
}
