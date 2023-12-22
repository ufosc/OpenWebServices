module github.com/ufosc/OpenWebServices/pkg/common

go 1.20

replace github.com/ufosc/OpenWebServices/pkg/authapi => ../authapi

replace github.com/ufosc/OpenWebServices/pkg/authdb => ../authdb

replace github.com/ufosc/OpenWebServices/pkg/websmtp => ../websmtp

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/ufosc/OpenWebServices/pkg/authdb v0.0.0-00010101000000-000000000000
	github.com/wagslane/go-password-validator v0.3.0
	golang.org/x/crypto v0.17.0
)

require (
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.13.1 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/text v0.14.0 // indirect
)
