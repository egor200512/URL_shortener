module github.com/egor200512/URL_shortener/services/links

go 1.25.5

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.4
	google.golang.org/genproto/googleapis/api v0.0.0-20260114163908-3f89685c29c3
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251222181119-0a764e51fe1b // indirect
)

replace github.com/egor200512/URL_shortener/shared => ../../shared
