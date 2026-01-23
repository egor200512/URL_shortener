module github.com/egor200512/URL_shortener/services/auth

go 1.25.5

require (
	github.com/egor200512/URL_shortener/shared v0.0.0
	github.com/georgysavva/scany v1.2.3
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.7.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.5
	github.com/jackc/pgx/v4 v4.18.3
	github.com/lpernett/godotenv v0.0.0-20230527005122-0de1d4c5ef5e
	github.com/stretchr/testify v1.11.1
	golang.org/x/crypto v0.46.0
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/egor200512/URL_shortener/shared => ../../shared
