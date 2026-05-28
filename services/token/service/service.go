package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	proto "ivpn.net/auth/services/proto"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
)

const maxInputBytes = 4096

// ErrAuthRequired is returned when the HSM session needs re-authentication.
var ErrAuthRequired = errors.New("hsm auth required")

type Signer interface {
	Generate(ctx context.Context, input string) (*model.HSMToken, error)
	Authenticate() error
}

type Server struct {
	Signer Signer
	proto.UnimplementedTokenServer
	Cfg *config.Config
}

func New(signer Signer, cfg config.Config) *Server {
	return &Server{
		Signer: signer,
		Cfg:    &cfg,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting token service on %s:%s", s.Cfg.Host, s.Cfg.Port)

	lis, err := net.Listen("tcp", s.Cfg.Host+":"+s.Cfg.Port)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	var opts []grpc.ServerOption
	if s.Cfg.TLSEnabled {
		tlsCfg, err := buildServerTLS(s.Cfg)
		if err != nil {
			return fmt.Errorf("tls config: %w", err)
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsCfg)))
		log.Println("Token service TLS (mTLS) enabled")
	} else {
		log.Println("WARNING: Token service is running without TLS — do not use in production")
	}

	srv := grpc.NewServer(opts...)
	proto.RegisterTokenServer(srv, s)

	if s.Cfg.Debug {
		reflection.Register(srv)
		log.Println("WARNING: gRPC reflection enabled — disable in production (TOKEN_DEBUG=false)")
	}

	if err = srv.Serve(lis); err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}

func buildServerTLS(cfg *config.Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.TLSCertFile, cfg.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("load server cert/key: %w", err)
	}

	caPEM, err := os.ReadFile(cfg.TLSCAFile)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("failed to parse CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

func (s *Server) Generate(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	if len(req.Input) > maxInputBytes {
		return nil, fmt.Errorf("input exceeds maximum allowed size of %d bytes", maxInputBytes)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	token, err := s.generateToken(reqCtx, req.Input)
	if err != nil {
		if isAuthError(err) {
			log.Println("Re-authenticating Signer session...")
			if authErr := s.Signer.Authenticate(); authErr == nil {
				token, err = s.generateToken(reqCtx, req.Input)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				return &proto.Response{Token: token.Token}, nil
			}
		}
		log.Println(err)
		return nil, err
	}

	return &proto.Response{Token: token.Token}, nil
}

// isAuthError detects HSM session expiry responses.
func isAuthError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "Status: 401") || strings.Contains(msg, "Status: 403")
}

func (s *Server) generateToken(ctx context.Context, input string) (*model.HSMToken, error) {
	return s.Signer.Generate(ctx, input)
}
