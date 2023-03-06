package lightning

import (
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

type LNDClient struct {
	Connection *grpc.ClientConn
	Client     lnrpc.LightningClient
}
