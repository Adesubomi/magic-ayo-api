package lightning

import (
	"context"
	"encoding/hex"
	"google.golang.org/grpc/credentials"
)

type RpcCredential map[string]string

var _ credentials.PerRPCCredentials = RpcCredential{}

func (m RpcCredential) RequireTransportSecurity() bool { return true }
func (m RpcCredential) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return m, nil
}

func newCredentials(macaroonBytes []byte) RpcCredential {
	c := make(map[string]string)
	c["macaroon"] = hex.EncodeToString(macaroonBytes)
	return c
}
