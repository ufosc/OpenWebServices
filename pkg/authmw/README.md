# authmw
[![Go Reference](https://pkg.go.dev/badge/github.com/ufosc/OpenWebServices/pkg/authmw.svg)](https://pkg.go.dev/github.com/ufosc/OpenWebServices/pkg/authmw)

authmw implements Golang Gin middleware for adding authentication to routes.

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
    "net/http"
)

func main() {
	r := gin.Default()

    // Create an authentication middleware function.
    // You'll need a database object here!
    mw := authmw.X(db, authmw.Config{
        Scope:  []string{"clients.read"},
        Realms: []string{"public"},
    })

    // Apply middleware to route.
    r.Get("/my/route", mw, func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "succesfully authenticated!",
        })
    })

	r.Run()
}
```

## License

[GNU AFFERO GENERAL PUBLIC LICENSE](https://github.com/ufosc/OpenWebServices/blob/main/pkg/authmw/LICENSE)

Copyright (C) 2023 Open Source Club
