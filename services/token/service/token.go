package service

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
	proto "ivpn.net/auth/services/token/proto"
)

type HSMClient interface {
	Token(input string, ttlMinutes int) (*model.HSMToken, error)
}

type Server struct {
	HSMClient HSMClient
	proto.UnimplementedTokenServer
	Cfg *config.Config
}

func New(hsm HSMClient, cfg config.Config) *Server {
	return &Server{
		HSMClient: hsm,
		Cfg:       &cfg,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting token service on %s:%s", s.Cfg.Host, s.Cfg.Port)

	lis, err := net.Listen("tcp", s.Cfg.Host+":"+s.Cfg.Port)
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
	token, err := s.generateToken(req.Input, int(req.TtlMinutes))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &proto.Response{
		Token: token.Token,
	}, nil
}

func (s *Server) generateToken(input string, ttlMinutes int) (*model.HSMToken, error) {
	return s.HSMClient.Token(input, ttlMinutes)
}
