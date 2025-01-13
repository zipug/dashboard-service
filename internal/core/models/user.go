package models

import "errors"

type (
	Id        int64
	Email     string
	Password  string
	Name      string
	LastName  string
	AvatarUrl string
)

type State string

const (
	Registered State = "registered"
	Verified   State = "verified"
	Banned     State = "banned"
	Deleted    State = "deleted"
)

var (
	ErrEmailEmpty          = errors.New("email is empty")
	ErrPasswordEmpty       = errors.New("password is empty")
	ErrRepeatPasswordEmpty = errors.New("repeat password is empty")
	ErrNameEmpty           = errors.New("name is empty")
	ErrLastNameEmpty       = errors.New("lastname is empty")
	ErrPasswordMismatch    = errors.New("passwords do not match")
	ErrInvalidEmail        = errors.New("invalid email")
)

func (email Email) isValid() error {
	if email == "" {
		return ErrEmailEmpty
	}
	return nil
}

func (password Password) isValid() error {
	if password == "" {
		return ErrPasswordEmpty
	}
	return nil
}

func (name Name) isValid() error {
	if name == "" {
		return ErrNameEmpty
	}
	return nil
}

func (lastName LastName) isValid() error {
	if lastName == "" {
		return ErrLastNameEmpty
	}
	return nil
}

func (password Password) compare(passwordToCompare Password) error {
	if password != passwordToCompare {
		return ErrPasswordMismatch
	}
	return nil
}

type User struct {
	Id             Id
	State          State
	Email          Email
	Password       Password
	RepeatPassword Password
	Name           Name
	LastName       LastName
	AvatarUrl      AvatarUrl
}

func (user User) IsValidForLogin() (bool, []error) {
	var errs []error
	if err := user.Email.isValid(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return false, errs
	}
	return true, nil
}

func (user User) IsValid() (bool, []error) {
	var errs []error
	if err := user.Email.isValid(); err != nil {
		errs = append(errs, err)
	}
	if err := user.Password.isValid(); err != nil {
		errs = append(errs, err)
	}
	if err := user.Name.isValid(); err != nil {
		errs = append(errs, err)
	}
	if err := user.LastName.isValid(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return false, errs
	}
	return true, nil
}
