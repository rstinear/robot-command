package client

import (
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	r "github.com/rstinear/workflow/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// New creates a new instance of the conker client for controlling robots remotely
func New(token string, address string, tls bool) r.CommandClient {

	creds := JWTCreds{
		Token: token,
	}

	dialOption := grpc.WithInsecure()
	if tls {
		tcreds := credentials.NewClientTLSFromCert(nil, "")
		dialOption = grpc.WithTransportCredentials(tcreds)
	}

	keeps := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		PermitWithoutStream: true,
		Time:                time.Second * 5,
		Timeout:             time.Second * 5,
	})

	conn, err := grpc.Dial(address, keeps, grpc.WithPerRPCCredentials(creds), dialOption)
	if err != nil {
		log.Fatalf("couldn't connect: %v", err)
	}

	rc := r.NewCommandClient(conn)

	return rc
}

// JWTCreds is a GRPC compliant credentials object for JWT tokens
type JWTCreds struct {
	Token string
}

// GetRequestMetadata returns a metadata formatted map for the JWTCreds type
func (j JWTCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {

	if j.Token == "" {
		return nil, errors.New("Token cannot be empty")
	}

	token := j.Token

	return map[string]string{
		"authorization": token,
	}, nil
}

// RequireTransportSecurity returns whether nor not these credentials require TLS
func (j JWTCreds) RequireTransportSecurity() bool {
	return false //TODO change this to true and configure TLS
}
