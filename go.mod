module RESTCryptoServer

go 1.23.0

toolchain go1.24.6

require github.com/go-chi/chi/v5 v5.2.3

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.41.0
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)
