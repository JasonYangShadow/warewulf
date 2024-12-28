package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type User struct {
	Name string `json:"name" yaml:"name"`
	Pass string `json:"pass" yaml:"pass"`
	Role string `json:"role" yaml:"role"`
}

type Authentication struct {
	Users   []User          `json:"users" yaml:"users"`
	conf    string          `json:"-" yaml:"-"`
	userMap map[string]User `json:"-" yaml:"-"`
}

func NewAuthentication() *Authentication {
	auth := new(Authentication)
	auth.userMap = make(map[string]User)
	return auth
}

func (auth *Authentication) ParseFromRaw(data []byte) error {
	err := yaml.Unmarshal(data, auth)
	if err != nil {
		return err
	}
	if len(auth.Users) == 0 {
		return fmt.Errorf("no record users")
	}
	for _, user := range auth.Users {
		if _, ok := auth.userMap[user.Name]; ok {
			return fmt.Errorf("duplicated user names")
		}
		auth.userMap[user.Name] = user
	}
	return nil
}

func (auth *Authentication) Read(confFileName string) error {
	if data, err := os.ReadFile(confFileName); err != nil {
		return err
	} else {
		auth.conf = confFileName
		if err := auth.ParseFromRaw(data); err != nil {
			return err
		}
	}
	return nil
}

var (
	UnauthorizedError  = errors.New("Unauthorized")
	PasswordWrongError = errors.New("Wrong password")
)

func (auth *Authentication) Authenticate(name, pass string) (*User, error) {
	if user, ok := auth.userMap[name]; !ok {
		return nil, fmt.Errorf("the user is not found")
	} else {
		if pass != user.Pass {
			return nil, PasswordWrongError
		}
		return &user, nil
	}
}

func (auth *Authentication) CheckAccess(user *User, requiredRole string) error {
	if target, ok := auth.userMap[user.Name]; !ok {
		return fmt.Errorf("the user is not found")
	} else {
		if requiredRole != target.Role && target.Role != "admin" {
			return UnauthorizedError
		}
	}
	return nil
}
