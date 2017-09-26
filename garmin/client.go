package garmin

import (
	"fmt"
	"os"
)

const (
	DefaultRootURI = "https://api.ua.com"
	envVarPrefix   = "GARMIN"
)

type Client struct {
	rootURI         string
	cookieAuthToken string
	accessToken     string // TODO
}

func New() *Client {
	rootURI := os.Getenv(envVarPrefix + "_ROOT_URI")
	if rootURI == "" {
		rootURI = DefaultRootURI
	}
	cookieAuthToken := os.Getenv(envVarPrefix + "_COOKIE_AUTH_TOKEN")
	return &Client{rootURI: rootURI, cookieAuthToken: cookieAuthToken}
}

func (c *Client) uri(path string, pathArgs ...interface{}) string {
	return fmt.Sprintf("%s%s", c.rootURI, fmt.Sprintf(path, pathArgs...))
}
