package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"ivpn.net/auth/services/generator/config"
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

func (c *TokenClient) GenerateToken(accountID string) (string, error) {
	req := &proto.Request{
		Input: accountID,
	}

	resp, err := c.Client.Generate(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return resp.Token, nil
}

func connect(cfg config.TokenServerConfig) (*grpc.ClientConn, error) {
	address := cfg.Host + ":" + cfg.Port
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	return conn, nil
}
