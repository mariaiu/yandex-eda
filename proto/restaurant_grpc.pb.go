// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: restaurant.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// YandexEdaClient is the client API for YandexEda service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type YandexEdaClient interface {
	GetRestaurants(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetRestaurantsResponse, error)
	GetRestaurant(ctx context.Context, in *GetRestaurantRequest, opts ...grpc.CallOption) (*GetRestaurantResponse, error)
	ParseRestaurants(ctx context.Context, in *ParseRestaurantsRequest, opts ...grpc.CallOption) (*ParseRestaurantsResponse, error)
}

type yandexEdaClient struct {
	cc grpc.ClientConnInterface
}

func NewYandexEdaClient(cc grpc.ClientConnInterface) YandexEdaClient {
	return &yandexEdaClient{cc}
}

func (c *yandexEdaClient) GetRestaurants(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetRestaurantsResponse, error) {
	out := new(GetRestaurantsResponse)
	err := c.cc.Invoke(ctx, "/yandexEda.YandexEda/GetRestaurants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yandexEdaClient) GetRestaurant(ctx context.Context, in *GetRestaurantRequest, opts ...grpc.CallOption) (*GetRestaurantResponse, error) {
	out := new(GetRestaurantResponse)
	err := c.cc.Invoke(ctx, "/yandexEda.YandexEda/GetRestaurant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yandexEdaClient) ParseRestaurants(ctx context.Context, in *ParseRestaurantsRequest, opts ...grpc.CallOption) (*ParseRestaurantsResponse, error) {
	out := new(ParseRestaurantsResponse)
	err := c.cc.Invoke(ctx, "/yandexEda.YandexEda/ParseRestaurants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// YandexEdaServer is the server API for YandexEda service.
// All implementations must embed UnimplementedYandexEdaServer
// for forward compatibility
type YandexEdaServer interface {
	GetRestaurants(context.Context, *emptypb.Empty) (*GetRestaurantsResponse, error)
	GetRestaurant(context.Context, *GetRestaurantRequest) (*GetRestaurantResponse, error)
	ParseRestaurants(context.Context, *ParseRestaurantsRequest) (*ParseRestaurantsResponse, error)
	mustEmbedUnimplementedYandexEdaServer()
}

// UnimplementedYandexEdaServer must be embedded to have forward compatible implementations.
type UnimplementedYandexEdaServer struct {
}

func (UnimplementedYandexEdaServer) GetRestaurants(context.Context, *emptypb.Empty) (*GetRestaurantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRestaurants not implemented")
}
func (UnimplementedYandexEdaServer) GetRestaurant(context.Context, *GetRestaurantRequest) (*GetRestaurantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRestaurant not implemented")
}
func (UnimplementedYandexEdaServer) ParseRestaurants(context.Context, *ParseRestaurantsRequest) (*ParseRestaurantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseRestaurants not implemented")
}
func (UnimplementedYandexEdaServer) mustEmbedUnimplementedYandexEdaServer() {}

// UnsafeYandexEdaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to YandexEdaServer will
// result in compilation errors.
type UnsafeYandexEdaServer interface {
	mustEmbedUnimplementedYandexEdaServer()
}

func RegisterYandexEdaServer(s grpc.ServiceRegistrar, srv YandexEdaServer) {
	s.RegisterService(&YandexEda_ServiceDesc, srv)
}

func _YandexEda_GetRestaurants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YandexEdaServer).GetRestaurants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yandexEda.YandexEda/GetRestaurants",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YandexEdaServer).GetRestaurants(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _YandexEda_GetRestaurant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRestaurantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YandexEdaServer).GetRestaurant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yandexEda.YandexEda/GetRestaurant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YandexEdaServer).GetRestaurant(ctx, req.(*GetRestaurantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YandexEda_ParseRestaurants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseRestaurantsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YandexEdaServer).ParseRestaurants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yandexEda.YandexEda/ParseRestaurants",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YandexEdaServer).ParseRestaurants(ctx, req.(*ParseRestaurantsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// YandexEda_ServiceDesc is the grpc.ServiceDesc for YandexEda service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var YandexEda_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "yandexEda.YandexEda",
	HandlerType: (*YandexEdaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRestaurants",
			Handler:    _YandexEda_GetRestaurants_Handler,
		},
		{
			MethodName: "GetRestaurant",
			Handler:    _YandexEda_GetRestaurant_Handler,
		},
		{
			MethodName: "ParseRestaurants",
			Handler:    _YandexEda_ParseRestaurants_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "restaurant.proto",
}
