package failure

import "github.com/redhajuanda/komon/fail"

var (
	ErrPathNotFound      = &fail.Failure{Code: "404000", Message: "There is no path that matches the given URI", HTTPStatus: 404}
	ErrNoteNotFound      = &fail.Failure{Code: "404001", Message: "Note not found", HTTPStatus: 404}
	ErrNamespaceNotFound = &fail.Failure{Code: "404002", Message: "Namespace not found", HTTPStatus: 404}
	ErrSecretFileNotFound = &fail.Failure{Code: "404003", Message: "Secret file not found", HTTPStatus: 404}
	ErrSecretNotFound    = &fail.Failure{Code: "404004", Message: "Secret not found", HTTPStatus: 404}

	ErrNoteAlreadyExists      = &fail.Failure{Code: "409001", Message: "Note already exists", HTTPStatus: 409}
	ErrSecretFileAlreadyExists = &fail.Failure{Code: "409002", Message: "Secret file already exists", HTTPStatus: 409}
	ErrSecretKeyAlreadyExists  = &fail.Failure{Code: "409003", Message: "Secret key already exists in this file", HTTPStatus: 409}
)