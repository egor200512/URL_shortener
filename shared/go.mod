module github.com/egor200512/URL_shortener/shared

go 1.25.5

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/egor200512/URL_shortener/services/links v0.0.0-20260211124845-9dd7721167b1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.5
	github.com/lpernett/godotenv v0.0.0-20230527005122-0de1d4c5ef5e
	github.com/redis/go-redis/v9 v9.17.2
	google.golang.org/genproto/googleapis/api v0.0.0-20260120221211-b8f7ae30c516
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
)
