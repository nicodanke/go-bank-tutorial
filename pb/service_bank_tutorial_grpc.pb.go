// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: service_bank_tutorial.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BankTutorialClient is the client API for BankTutorial service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BankTutorialClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

type bankTutorialClient struct {
	cc grpc.ClientConnInterface
}

func NewBankTutorialClient(cc grpc.ClientConnInterface) BankTutorialClient {
	return &bankTutorialClient{cc}
}

func (c *bankTutorialClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/pb.BankTutorial/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankTutorialClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/pb.BankTutorial/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BankTutorialServer is the server API for BankTutorial service.
// All implementations must embed UnimplementedBankTutorialServer
// for forward compatibility
type BankTutorialServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	mustEmbedUnimplementedBankTutorialServer()
}

// UnimplementedBankTutorialServer must be embedded to have forward compatible implementations.
type UnimplementedBankTutorialServer struct {
}

func (UnimplementedBankTutorialServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedBankTutorialServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedBankTutorialServer) mustEmbedUnimplementedBankTutorialServer() {}

// UnsafeBankTutorialServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BankTutorialServer will
// result in compilation errors.
type UnsafeBankTutorialServer interface {
	mustEmbedUnimplementedBankTutorialServer()
}

func RegisterBankTutorialServer(s grpc.ServiceRegistrar, srv BankTutorialServer) {
	s.RegisterService(&BankTutorial_ServiceDesc, srv)
}

func _BankTutorial_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankTutorialServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankTutorial/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankTutorialServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BankTutorial_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankTutorialServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankTutorial/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankTutorialServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BankTutorial_ServiceDesc is the grpc.ServiceDesc for BankTutorial service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BankTutorial_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BankTutorial",
	HandlerType: (*BankTutorialServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _BankTutorial_Login_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _BankTutorial_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_bank_tutorial.proto",
}
