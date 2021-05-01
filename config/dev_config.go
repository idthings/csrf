// +build dev test

//
// The build tag means this file is included when the prod tag is used.
//
// eg: go run api/server.gp
//
// https://golang.org/pkg/go/build/

package config

const (
	API_PORT = "8080"

	DB_HOST = "127.0.0.1"
	DB_PORT = "6379"
)
