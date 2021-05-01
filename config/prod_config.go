// +build !dev,!test

// The build tag means this file is included when the prod tag is used.
//
// eg: go run -tags prod api/server.gp
//
// https://golang.org/pkg/go/build/

package config

const (
	API_PORT = "80"

	DB_HOST = "localhost"
	DB_PORT = "6379"
)
