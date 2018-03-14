package kms

import (
	"flag"
	"fmt"
	"net"
	"net/url"

	"github.com/enj/kms/api/v1beta1"
	"github.com/enj/kms/pkg/kms"

	"google.golang.org/grpc"
)

const (
	// unix domain socket is the only supported protocol
	unixProtocol = "unix"
)

var (
	endpointFlag = flag.String("endpoint", "", `the address to listen on, for example "unix:///var/run/kms-provider.sock"`)
)

func Execute() error {
	flag.Parse()

	endpoint, err := parseEndpoint(endpointFlag)
	if err != nil {
		return err
	}
	listener, err := net.Listen(unixProtocol, endpoint)
	if err != nil {
		return err
	}

	clevisKMS, err := kms.NewClevisKMS()
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption // TODO see if we need any options
	grpcServer := grpc.NewServer(opts...)
	v1beta1.RegisterKeyManagementServiceServer(grpcServer, clevisKMS)

	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func parseEndpoint(endpointFlag *string) (string, error) {
	if endpointFlag == nil {
		return "", fmt.Errorf("no endpoint provided")
	}

	endpoint := *endpointFlag

	if len(endpoint) == 0 {
		return "", fmt.Errorf("cannot use empty string as endpoint")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint %q: %v", endpoint, err)
	}

	if u.Scheme != unixProtocol {
		return "", fmt.Errorf("unsupported scheme %q, must be %q", u.Scheme, unixProtocol)
	}

	return u.Path, nil
}