package controller

import (
	"context"
	"log"
	"net"

	"github.com/Mark-Grigorev/auth/internal/gen/proto"
	"github.com/Mark-Grigorev/auth/internal/logic"
	"github.com/Mark-Grigorev/auth/internal/model"
	"google.golang.org/grpc"
)

type Controller struct {
	listener net.Listener
	grpc     *grpc.Server
	config   model.AppConfig
	proto.UnimplementedAuthServiceServer
	logic *logic.Logic
}

func New(cfg model.AppConfig, logic *logic.Logic) *Controller {
	listener, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		log.Fatal(err.Error())
	}
	grpc := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpc, proto.UnimplementedAuthServiceServer{})
	return &Controller{
		listener: listener,
		grpc:     grpc,
		config:   cfg,
		logic:    logic,
	}
}

func (c *Controller) Start() {
	if err := c.grpc.Serve(c.listener); err != nil {
		log.Fatal(err)
	}
}

func (c *Controller) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	userID, err := c.logic.Register(ctx, &model.UserRegistrationData{
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		Login:      req.Login,
		Password:   req.Password,
	})
	if err != nil {
		return &proto.RegisterResponse{}, err
	}
	return &proto.RegisterResponse{
		UserId: userID,
	}, nil
}

func (c *Controller) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := c.logic.Authorization(ctx, req.Login, req.Password)
	if err != nil {
		return &proto.LoginResponse{}, err
	}
	return &proto.LoginResponse{
		Token: token,
	}, nil
}

func (c *Controller) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	valid, err := c.logic.ValidateToken(ctx, req.Token)
	if err != nil {
		return &proto.ValidateTokenResponse{}, err
	}
	return &proto.ValidateTokenResponse{
		Valid: valid,
	}, nil
}
