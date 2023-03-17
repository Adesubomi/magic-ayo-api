package lightning

import (
	"fmt"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
	"os"
	"os/user"
	"path"
)

func newClient(lnConfig configPkg.LightningConfig) (*LNClient, error) {
	lndClient := new(LNClient)

	usr, err := user.Current()
	if err != nil {
		fmt.Println("    [✗] Cannot get current user:", err)
		return nil, err
	}

	fmt.Println("    [-] Home directory::", usr.HomeDir)

	tlsCertPath := path.Join(usr.HomeDir, lnConfig.TlsCert)
	macaroonPath := path.Join(usr.HomeDir, lnConfig.Macaroon)
	tlsCredentials, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		fmt.Println("    [✗] Cannot get node tls credentials:", err)
	}
	fmt.Println("    [-] TLS credentials loaded::", lnConfig.TlsCert)

	macaroonBytes, err := os.ReadFile(macaroonPath)
	if err != nil {
		fmt.Println("    [✗] Cannot read macaroon files:", err)
		return nil, err
	}
	fmt.Println("    [-] Macaroons loaded::", lnConfig.Macaroon)

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
		fmt.Println("    [✗] Cannot unmarshal macaroon", err)
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			newCredentials(macaroonBytes)),
	}

	url := fmt.Sprintf(
		"%v:%v",
		lnConfig.Url,
		lnConfig.Port,
	)

	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}

	logPkg.PrintlnGreen("    [✔] gRPC Connection opened @" + url)
	lndClient.Connection = conn

	// Create LightningClient using the LND gRPC connection
	client := lnrpc.NewLightningClient(conn)
	lndClient.Client = client

	return lndClient, nil
}

func NewSenderLnClient(config *configPkg.Config) (*LNClient, error) {
	fmt.Println("")
	fmt.Println("  [...] Connecting to LND over gRPC - Sender")
	return newClient(config.LnSendNode)
}

func NewRecipientLnClient(config *configPkg.Config) (*LNClient, error) {
	fmt.Println("")
	fmt.Println("  [...] Connecting to LND over gRPC - Recipient")
	return newClient(config.LnReceiveNode)
}
