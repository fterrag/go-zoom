# go-zoom

The `zoom` packages provides a lightweight [Zoom API](https://marketplace.zoom.us/docs/api-reference/introduction/) client. Coverage of endpoints is minimal, but [users.go](zoom/users.go) and [meetings.go](zoom/meetings.go) should act as good examples for implementing support for additional endpoints.

This package is built to be used with [Server-to-Server OAuth](https://marketplace.zoom.us/docs/guides/build/server-to-server-oauth-app/) apps.

## Example Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fterrag/go-zoom/zoom"
)

func main() {
	ctx := context.Background()

	httpClient := &http.Client{}
	client := zoom.NewClient(httpClient, os.Getenv("ZOOM_ACCOUNT_ID"), os.Getenv("ZOOM_CLIENT_ID"), os.Getenv("ZOOM_CLIENT_SECRET"))

	res, _, err := client.Users.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d users\n\n", len(res.Users))

	for _, user := range res.Users {
		fmt.Printf("ID: %s\nDisplay Name: %s\nEmail: %s\n\n", user.ID, user.DisplayName, user.Email)
	}
}
```
