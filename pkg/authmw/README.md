# authmw
[![Go Reference](https://pkg.go.dev/badge/github.com/ufosc/OpenWebServices/pkg/authmw.svg)](https://pkg.go.dev/github.com/ufosc/OpenWebServices/pkg/authmw)

authmw implements Golang Gin middleware for authenticating users & clients via JWT and OAuth2 bearer tokens.

## Install
```bash
go get github.com/ufosc/OpenWebServices/pkg/authmw
```

## Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authmw"
)

func main() {
	r := gin.Default()

	// Authenticates a user and expects them to have realm1, realm2
	// privileges.
	r.GET("/some/user/route", authmw.AuthenticateUser("secret", database,
	"scope1", "scope2"), func(c *gin.Context) {
		// Some middleware here.
	})

	// Authenticates a client.
	r.GET("/some/user/route", authmw.AuthenticateClient("secret", database),
		func(c *gin.Context) {
		// Some middleware here.
	})

	// Authenticates an OAuth2 access token (bearer schema). Expects
	// assosciated client to have scope1, scope2, realm1, realm2 access.
	r.GET("/some/user/route", authmw.AuthenticateBearer(database,
		[]string{"realm1", "realm2"}, []string{"scope1", scope2"}),
		func(c *gin.Context) {
		// Some middleware here.
	})

	r.Run()
}
```

## License

[GNU AFFERO GENERAL PUBLIC LICENSE](https://github.com/ufosc/OpenWebServices/blob/main/pkg/authmw/LICENSE)

Copyright (C) 2023 Open Source Club
