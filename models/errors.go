package models

import "strings"

const (
	ErrFollowSelf       modelError = "models: cannot follow yourself"
	ErrFollowExists     modelError = "models: you have followed this user already"
	ErrTagExists        modelError = "models: tag name already exists"
	ErrTagNoSpecialChar modelError = "models: no special characters allowed"
	ErrTaggingExists    modelError = "models: tag associated with this tweet already exists"
	ErrLikeExists       modelError = "models: you have liked this tweet already"
	ErrRetweetExists    modelError = "models: you have retweeted this tweet already"
	// ErrFollowSelf modelError = "models: "
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound modelError = "models: resource not found"
	// ErrPasswordIncorrect is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrPasswordIncorrect modelError = "models: incorrect password provided"

	ErrUsernameRequired modelError = "models: username is required"
	ErrUsernameNoLetter modelError = "models: username must begin with a letter"

	ErrNameRequired modelError = "models: name is required"
	// ErrEmailRequired is returned when an email address is
	// not provided when creating a user
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid is returned when an email address provided
	// does not match any of our requirements
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken is returned when an update or create is attempted
	// with an email address that is already in use.
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrPasswordRequired is returned when a create is attempted
	// without a user password provided.
	ErrPasswordRequired modelError = "models: password is required"
	// ErrPasswordTooShort is returned when an update or create is
	// attempted with a user password that is less than 8 characters.
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	ErrPostRequired     modelError = "models: post is required"

	ErrTokenInvalid modelError = "models: token provided is not valid"

	// ErrIDInvalid is returned when an invalid ID is provided
	// to a method like Delete.
	ErrIDInvalid privateError = "models: ID provided was invalid"
	// ErrRememberRequired is returned when a create or update
	// is attempted without a user remember token hash
	ErrRememberRequired privateError = "models: remember token is required"
	// ErrRememberTooShort is returned when a remember token is
	// not at least 32 bytes
	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"
	ErrUserIDRequired   privateError = "models: user ID is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}

// ErrUsernameMin      modelError = "models: length of username must be greater than: "
// ErrUsernameMax      modelError = "models: length of username must be less than: "
