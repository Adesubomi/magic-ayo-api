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

func NewLightningClient(config *configPkg.Config) (*LNDClient, error) {
	fmt.Println("")
	fmt.Println("  [...] Connecting to LND over gRPC")
	lndClient := new(LNDClient)

	usr, err := user.Current()
	if err != nil {
		fmt.Println("    [✗] Cannot get current user:", err)
		return nil, err
	}

	fmt.Println("    [-] Home directory::", usr.HomeDir)

	tlsCertPath := path.Join(usr.HomeDir, config.Lightning.TlsCert)
	macaroonPath := path.Join(usr.HomeDir, config.Lightning.Macaroon)
	tlsCredentials, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		fmt.Println("    [✗] Cannot get node tls credentials:", err)
	}
	fmt.Println("    [-] TLS credentials loaded::", config.Lightning.TlsCert)

	macaroonBytes, err := os.ReadFile(macaroonPath)
	if err != nil {
		fmt.Println("    [✗] Cannot read macaroon files:", err)
		return nil, err
	}
	fmt.Println("    [-] Macaroons loaded::", config.Lightning.Macaroon)

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
		config.Lightning.Url,
		config.Lightning.Port,
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
