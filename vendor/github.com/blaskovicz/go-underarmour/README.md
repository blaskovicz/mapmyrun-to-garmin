# go-underarmour
> Golang library for interacting with the [Under Armour](https://developer.underarmour.com) API.

## Install

```
$ go get github.com/blaskovicz/go-underarmour
```

## Use

```go
import (
  ua "github.com/blaskovicz/go-underarmour"
)

// initialize a default client
// requires setting UNDERARMOUR_COOKIE_AUTH_TOKEN env var to your
// auth-token cookie value
client := ua.New()

// then fetch your profile info
u, err := client.ReadUser("self")
if err != nil {
  panic(err)
}

fmt.Printf("Authenticated as user %s\n", u.Username)

// more to come, including using access and refresh token.
```

## Test

```
$ go test ./...
```
