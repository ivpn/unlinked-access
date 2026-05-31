package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"ivpn.net/auth/services/preauth/config"
	proto "ivpn.net/auth/services/proto"
)

type TokenClient struct {
	Client proto.TokenClient
}

func New(cfg config.TokenServerConfig) (*TokenClient, error) {
	conn, err := connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to token server: %w", err)
	}

	client := proto.NewTokenClient(conn)

	return &TokenClient{
		Client: client,
	}, nil
}

func (c *TokenClient) GenerateToken(input string) (string, error) {
	req := &proto.Request{
		Input: input,
	}

	resp, err := c.Client.Generate(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return resp.Token, nil
}

func connect(cfg config.TokenServerConfig) (*grpc.ClientConn, error) {
	address := cfg.Host + ":" + cfg.Port

	var creds grpc.DialOption
	if cfg.TLSEnabled {
		tlsCfg, err := buildClientTLS(cfg)
		if err != nil {
			return nil, fmt.Errorf("tls config: %w", err)
		}
		creds = grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
	} else {
		if os.Getenv("PREAUTH_ALLOW_INSECURE") != "true" {
			return nil, errors.New("TLS is disabled but PREAUTH_ALLOW_INSECURE is not set to 'true'; refusing insecure connection")
		}
		log.Println("WARNING: gRPC connection to token server is unencrypted (TLS disabled)")
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(address, creds)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	return conn, nil
}

func buildClientTLS(cfg config.TokenServerConfig) (*tls.Config, error) {
	caPEM, err := os.ReadFile(cfg.TLSCACertFile)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("failed to parse CA certificate")
	}

	cert, err := tls.LoadX509KeyPair(cfg.TLSCertFile, cfg.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("load client cert/key: %w", err)
	}

	return &tls.Config{
		RootCAs:      caPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}, nil
}
