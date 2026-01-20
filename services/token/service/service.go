package service

import (
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	proto "ivpn.net/auth/services/proto"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
)

type Signer interface {
	Generate(input string) (*model.HSMToken, error)
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

	lis, err := net.Listen("tcp", ":"+s.Cfg.Port)
	if err != nil {
		log.Println(err)
		return err
	}

	srv := grpc.NewServer()
	proto.RegisterTokenServer(srv, s)
	reflection.Register(srv)

	err = srv.Serve(lis)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Server) Generate(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	token, err := s.generateToken(req.Input)
	if err != nil {
		if strings.Contains(err.Error(), "Status: 401") || strings.Contains(err.Error(), "Status: 403") {
			log.Println("Re-authenticating Signer client session...")
			err = s.Signer.Authenticate()
			if err == nil {
				token, err = s.generateToken(req.Input)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				return &proto.Response{
					Token: token.Token,
				}, nil
			}
		}

		log.Println(err)
		return nil, err
	}

	return &proto.Response{
		Token: token.Token,
	}, nil
}

func (s *Server) generateToken(input string) (*model.HSMToken, error) {
	return s.Signer.Generate(input)
}
