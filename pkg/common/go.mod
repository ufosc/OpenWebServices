module github.com/ufosc/OpenWebServices/pkg/common

go 1.20

replace github.com/ufosc/OpenWebServices/pkg/authapi => ../authapi

replace github.com/ufosc/OpenWebServices/pkg/authdb => ../authdb

replace github.com/ufosc/OpenWebServices/pkg/websmtp => ../websmtp

require (
	github.com/google/uuid v1.6.0
	github.com/wagslane/go-password-validator v0.3.0
	golang.org/x/crypto v0.17.0
)
