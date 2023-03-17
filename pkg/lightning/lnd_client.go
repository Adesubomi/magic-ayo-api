package lightning

import (
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

type LNClient struct {
	Connection *grpc.ClientConn
	Client     lnrpc.LightningClient
}

type LNClients struct {
	Send    *LNClient
	Receive *LNClient
}
