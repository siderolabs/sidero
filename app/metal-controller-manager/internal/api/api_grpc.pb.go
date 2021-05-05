// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentClient interface {
	CreateServer(ctx context.Context, in *CreateServerRequest, opts ...grpc.CallOption) (*CreateServerResponse, error)
	MarkServerAsWiped(ctx context.Context, in *MarkServerAsWipedRequest, opts ...grpc.CallOption) (*MarkServerAsWipedResponse, error)
	ReconcileServerAddresses(ctx context.Context, in *ReconcileServerAddressesRequest, opts ...grpc.CallOption) (*ReconcileServerAddressesResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	UpdateBMCInfo(ctx context.Context, in *UpdateBMCInfoRequest, opts ...grpc.CallOption) (*UpdateBMCInfoResponse, error)
}

type agentClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentClient(cc grpc.ClientConnInterface) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) CreateServer(ctx context.Context, in *CreateServerRequest, opts ...grpc.CallOption) (*CreateServerResponse, error) {
	out := new(CreateServerResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/CreateServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) MarkServerAsWiped(ctx context.Context, in *MarkServerAsWipedRequest, opts ...grpc.CallOption) (*MarkServerAsWipedResponse, error) {
	out := new(MarkServerAsWipedResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/MarkServerAsWiped", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) ReconcileServerAddresses(ctx context.Context, in *ReconcileServerAddressesRequest, opts ...grpc.CallOption) (*ReconcileServerAddressesResponse, error) {
	out := new(ReconcileServerAddressesResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/ReconcileServerAddresses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) UpdateBMCInfo(ctx context.Context, in *UpdateBMCInfoRequest, opts ...grpc.CallOption) (*UpdateBMCInfoResponse, error) {
	out := new(UpdateBMCInfoResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/UpdateBMCInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServer is the server API for Agent service.
// All implementations must embed UnimplementedAgentServer
// for forward compatibility
type AgentServer interface {
	CreateServer(context.Context, *CreateServerRequest) (*CreateServerResponse, error)
	MarkServerAsWiped(context.Context, *MarkServerAsWipedRequest) (*MarkServerAsWipedResponse, error)
	ReconcileServerAddresses(context.Context, *ReconcileServerAddressesRequest) (*ReconcileServerAddressesResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	UpdateBMCInfo(context.Context, *UpdateBMCInfoRequest) (*UpdateBMCInfoResponse, error)
	mustEmbedUnimplementedAgentServer()
}

// UnimplementedAgentServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct{}

func (UnimplementedAgentServer) CreateServer(context.Context, *CreateServerRequest) (*CreateServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateServer not implemented")
}

func (UnimplementedAgentServer) MarkServerAsWiped(context.Context, *MarkServerAsWipedRequest) (*MarkServerAsWipedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkServerAsWiped not implemented")
}

func (UnimplementedAgentServer) ReconcileServerAddresses(context.Context, *ReconcileServerAddressesRequest) (*ReconcileServerAddressesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReconcileServerAddresses not implemented")
}

func (UnimplementedAgentServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}

func (UnimplementedAgentServer) UpdateBMCInfo(context.Context, *UpdateBMCInfoRequest) (*UpdateBMCInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBMCInfo not implemented")
}
func (UnimplementedAgentServer) mustEmbedUnimplementedAgentServer() {}

// UnsafeAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServer will
// result in compilation errors.
type UnsafeAgentServer interface {
	mustEmbedUnimplementedAgentServer()
}

func RegisterAgentServer(s grpc.ServiceRegistrar, srv AgentServer) {
	s.RegisterService(&Agent_ServiceDesc, srv)
}

func _Agent_CreateServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).CreateServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/CreateServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).CreateServer(ctx, req.(*CreateServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_MarkServerAsWiped_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkServerAsWipedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).MarkServerAsWiped(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/MarkServerAsWiped",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).MarkServerAsWiped(ctx, req.(*MarkServerAsWipedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_ReconcileServerAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReconcileServerAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).ReconcileServerAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/ReconcileServerAddresses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).ReconcileServerAddresses(ctx, req.(*ReconcileServerAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_UpdateBMCInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBMCInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).UpdateBMCInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/UpdateBMCInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).UpdateBMCInfo(ctx, req.(*UpdateBMCInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Agent_ServiceDesc is the grpc.ServiceDesc for Agent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Agent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateServer",
			Handler:    _Agent_CreateServer_Handler,
		},
		{
			MethodName: "MarkServerAsWiped",
			Handler:    _Agent_MarkServerAsWiped_Handler,
		},
		{
			MethodName: "ReconcileServerAddresses",
			Handler:    _Agent_ReconcileServerAddresses_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _Agent_Heartbeat_Handler,
		},
		{
			MethodName: "UpdateBMCInfo",
			Handler:    _Agent_UpdateBMCInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
