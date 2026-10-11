package apperrors

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden: access denied")
	ErrBadRequest          = errors.New("invalid input data")
	ErrTokenExpired        = errors.New("token has expired")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvitationExpired   = errors.New("invitation has expired")
	ErrInvitationAccepted  = errors.New("invitation already accepted")
	ErrUserInactive        = errors.New("user account is inactive")
	ErrOrganizationInactive = errors.New("organization is inactive")
	ErrInternal            = errors.New("internal server error")
)
