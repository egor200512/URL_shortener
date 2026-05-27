package jwt

import (
	"context"
	"errors"
	"testing"

	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAuthIntersepter(t *testing.T) {
	const (
		token  = "access-token"
		userID = "user-1"
	)

	tests := []struct {
		name       string
		ctx        context.Context
		authClient desc.AuthServiceClient
		wantCode   codes.Code
		wantUserID string
	}{
		{
			name:       "no authorization metadata",
			ctx:        context.Background(),
			authClient: &fakeAuthClient{},
			wantCode:   codes.Unauthenticated,
		},
		{
			name:       "empty token",
			ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "")),
			authClient: &fakeAuthClient{},
			wantCode:   codes.Unauthenticated,
		},
		{
			name:       "auth service error",
			ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", TOKEN_PREFIX+token)),
			authClient: &fakeAuthClient{verifyErr: errors.New("invalid token")},
			wantCode:   codes.Unauthenticated,
		},
		{
			name:       "auth service returned user id",
			ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", TOKEN_PREFIX+token)),
			authClient: &fakeAuthClient{userID: userID},
			wantCode:   codes.OK,
			wantUserID: userID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := AuthIntersepter(tt.authClient)

			resp, err := interceptor(tt.ctx, "request", &grpc.UnaryServerInfo{}, func(ctx context.Context, req any) (any, error) {
				require.Equal(t, "request", req)
				require.Equal(t, tt.wantUserID, ctx.Value(ClaimsCtxKey))
				return "response", nil
			})

			require.Equal(t, tt.wantCode, status.Code(err))
			if tt.wantCode == codes.OK {
				require.NoError(t, err)
				require.Equal(t, "response", resp)
			}
		})
	}
}

type fakeAuthClient struct {
	desc.AuthServiceClient
	userID    string
	verifyErr error
}

func (f *fakeAuthClient) Register(context.Context, *desc.RegisterRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) Login(context.Context, *desc.LoginRequest, ...grpc.CallOption) (*desc.LoginResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) VerifyToken(_ context.Context, req *desc.VerifyTokenRequest, _ ...grpc.CallOption) (*desc.VerifyTokenResponse, error) {
	if f.verifyErr != nil {
		return nil, f.verifyErr
	}
	if req.GetAccessToken() == "" {
		return nil, errors.New("empty token")
	}
	return &desc.VerifyTokenResponse{UserId: f.userID}, nil
}
