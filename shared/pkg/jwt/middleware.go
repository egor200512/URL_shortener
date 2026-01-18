package jwt

// import (
// 	"context"
// 	"strings"

// 	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/status"
// )

// const (
// 	TOKEN_PREFIX = "Bearer "
// )

// type ClaimsCtx string

// var ClaimsCtxKey ClaimsCtx = "ClaimsCtx"

// func AuthIntersepter(authClient desc.AuthServiceClient) grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req any,
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (resp any, err error) {

// 		md, ok := metadata.FromIncomingContext(ctx)
// 		if !ok {
// 			return nil, status.Errorf(codes.Unauthenticated, "No metadata: %v", err)
// 		}

// 		tk := md["authorization"]
// 		if len(tk) == 0 {
// 			return nil, status.Errorf(codes.Unauthenticated, "No token: %v", err)
// 		}

// 		if !strings.HasPrefix(tk[0], TOKEN_PREFIX) {
// 			return nil, status.Error(codes.Unauthenticated, "Unexpected token structure: %v")
// 		}

// 		token := strings.TrimPrefix(tk[0], TOKEN_PREFIX)

// 		// cl, err := VerifyToken(token, []byte(os.Getenv(configs.JWT_SECRET_KEY)))
// 		// if err != nil || cl == nil {
// 		// 	return nil, status.Errorf(codes.Unauthenticated, "No access: %v", err)
// 		// }

// 		newCtx := context.WithValue(ctx, ClaimsCtxKey, cl)

// 		return handler(newCtx, req)
// 	}
// }
