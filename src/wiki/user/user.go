package user

import (
	"bytes"
	"crypto/sha512"
	"io/ioutil"
	"regexp"
	"strings"
)

type User interface {
	UserName() string
	Validate(string) error
	Add(string) error
}

type userError string

func (e userError) Error() string {
	return string(e)
}

const INVALID_DATABASE_ERROR userError = "The database identifier is not valid"
const USER_DOES_NOT_EXIST userError = "The user does not exist"
const INCORRECT_PASSWORD userError = "The password is incorrect"
const USER_EXISTS userError = "The user exists"

type txtUser struct {
	file           string
	userName       string
	passwordSha512 []byte
}

func New(file, userName string) (u *txtUser) {
	return &txtUser{file: file, userName: userName}
}

func (u *txtUser) UserName() string {
	return u.userName
}

var validPath = regexp.MustCompile("^[a-zA-z0-9/]+.db$")

type passwordFunc func(*txtUser, []byte, string) error

func (u *txtUser) passwordTool(f passwordFunc, password string) error {
	if !validPath.MatchString(u.file) {
		return INVALID_DATABASE_ERROR
	}
	data, err := ioutil.ReadFile(u.file)
	if err != nil {
		return err
	}
	return f(u, data, password)
}

func validate(u *txtUser, data []byte, password string) error {
	d := string(data)
	i := strings.Index(d, u.userName)
	if i < 0 {
		return USER_DOES_NOT_EXIST
	}
	//if u.passwordSha512 != nil { TODO
		u.passwordSha512 = []byte(strings.Split(d[i:], "\ufb4f")[1])
	//}
	passwordSha512 := sha512.Sum512([]byte(password))
	if !bytes.Equal(u.passwordSha512, passwordSha512[:]) {
		return INCORRECT_PASSWORD
	}
	return nil
}

func (u *txtUser) Validate(password string) error {
	return u.passwordTool(validate, password)
}

func add(u *txtUser, data []byte, password string) error {
	d := string(data)
	i := strings.Index(d, u.userName)
	if i >= 0 {
		return USER_EXISTS
	}
	passwordSha512 := sha512.Sum512([]byte(password))
	u.passwordSha512 = passwordSha512[:]
	err := ioutil.WriteFile(u.file, append(data, []byte(u.userName+"\ufb4f"+string(u.passwordSha512)+"\ufb4f")...), 060000)
	if err != nil {
		return err
	}
	return nil
}

func (u *txtUser) Add(password string) error {
	return u.passwordTool(add, password)
}
