// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.13.0
// source: cntlrgrpc.proto

package cntlrgrpc

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

// CntlrGrpcClient is the client API for CntlrGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CntlrGrpcClient interface {
	CntlrGetParamReq(ctx context.Context, in *CntlrGetParamReqData, opts ...grpc.CallOption) (*CntlrReqResult, error)
	CntlrSetParamReq(ctx context.Context, in *CntlrSetParamReqData, opts ...grpc.CallOption) (*CntlrSetParamResData, error)
	CntlrGetInstancesReq(ctx context.Context, in *CntlrGetInstancesReqData, opts ...grpc.CallOption) (*CntlrReqResult, error)
	CntlrAddInstanceReq(ctx context.Context, in *CntlrAddInstanceReqData, opts ...grpc.CallOption) (*CntlrAddInstanceResData, error)
	CntlrOperateReq(ctx context.Context, in *CntlrOperateReqData, opts ...grpc.CallOption) (*CntlrOperateResData, error)
	CntlrGetDatamodelReq(ctx context.Context, in *CntlrGetDatamodelReqData, opts ...grpc.CallOption) (*CntlrReqResult, error)
	CntlrDeleteInstanceReq(ctx context.Context, in *CntlrDeleteInstanceReqData, opts ...grpc.CallOption) (*CntlrReqResult, error)
	CntlrGetAgentMsgs(ctx context.Context, in *CntlrGetAgentMsgsData, opts ...grpc.CallOption) (*CntlrReqResult, error)
	CntlrGetInfo(ctx context.Context, in *None, opts ...grpc.CallOption) (*CntlrInfoData, error)
	CntlrStream(ctx context.Context, in *CntlrGetParamReqData, opts ...grpc.CallOption) (CntlrGrpc_CntlrStreamClient, error)
}

type cntlrGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewCntlrGrpcClient(cc grpc.ClientConnInterface) CntlrGrpcClient {
	return &cntlrGrpcClient{cc}
}

func (c *cntlrGrpcClient) CntlrGetParamReq(ctx context.Context, in *CntlrGetParamReqData, opts ...grpc.CallOption) (*CntlrReqResult, error) {
	out := new(CntlrReqResult)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrGetParamReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrSetParamReq(ctx context.Context, in *CntlrSetParamReqData, opts ...grpc.CallOption) (*CntlrSetParamResData, error) {
	out := new(CntlrSetParamResData)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrSetParamReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrGetInstancesReq(ctx context.Context, in *CntlrGetInstancesReqData, opts ...grpc.CallOption) (*CntlrReqResult, error) {
	out := new(CntlrReqResult)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrGetInstancesReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrAddInstanceReq(ctx context.Context, in *CntlrAddInstanceReqData, opts ...grpc.CallOption) (*CntlrAddInstanceResData, error) {
	out := new(CntlrAddInstanceResData)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrAddInstanceReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrOperateReq(ctx context.Context, in *CntlrOperateReqData, opts ...grpc.CallOption) (*CntlrOperateResData, error) {
	out := new(CntlrOperateResData)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrOperateReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrGetDatamodelReq(ctx context.Context, in *CntlrGetDatamodelReqData, opts ...grpc.CallOption) (*CntlrReqResult, error) {
	out := new(CntlrReqResult)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrGetDatamodelReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrDeleteInstanceReq(ctx context.Context, in *CntlrDeleteInstanceReqData, opts ...grpc.CallOption) (*CntlrReqResult, error) {
	out := new(CntlrReqResult)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrDeleteInstanceReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrGetAgentMsgs(ctx context.Context, in *CntlrGetAgentMsgsData, opts ...grpc.CallOption) (*CntlrReqResult, error) {
	out := new(CntlrReqResult)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrGetAgentMsgs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrGetInfo(ctx context.Context, in *None, opts ...grpc.CallOption) (*CntlrInfoData, error) {
	out := new(CntlrInfoData)
	err := c.cc.Invoke(ctx, "/cntrlgrpc.CntlrGrpc/CntlrGetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cntlrGrpcClient) CntlrStream(ctx context.Context, in *CntlrGetParamReqData, opts ...grpc.CallOption) (CntlrGrpc_CntlrStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &CntlrGrpc_ServiceDesc.Streams[0], "/cntrlgrpc.CntlrGrpc/CntlrStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &cntlrGrpcCntlrStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CntlrGrpc_CntlrStreamClient interface {
	Recv() (*CntlrReqResult, error)
	grpc.ClientStream
}

type cntlrGrpcCntlrStreamClient struct {
	grpc.ClientStream
}

func (x *cntlrGrpcCntlrStreamClient) Recv() (*CntlrReqResult, error) {
	m := new(CntlrReqResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CntlrGrpcServer is the server API for CntlrGrpc service.
// All implementations must embed UnimplementedCntlrGrpcServer
// for forward compatibility
type CntlrGrpcServer interface {
	CntlrGetParamReq(context.Context, *CntlrGetParamReqData) (*CntlrReqResult, error)
	CntlrSetParamReq(context.Context, *CntlrSetParamReqData) (*CntlrSetParamResData, error)
	CntlrGetInstancesReq(context.Context, *CntlrGetInstancesReqData) (*CntlrReqResult, error)
	CntlrAddInstanceReq(context.Context, *CntlrAddInstanceReqData) (*CntlrAddInstanceResData, error)
	CntlrOperateReq(context.Context, *CntlrOperateReqData) (*CntlrOperateResData, error)
	CntlrGetDatamodelReq(context.Context, *CntlrGetDatamodelReqData) (*CntlrReqResult, error)
	CntlrDeleteInstanceReq(context.Context, *CntlrDeleteInstanceReqData) (*CntlrReqResult, error)
	CntlrGetAgentMsgs(context.Context, *CntlrGetAgentMsgsData) (*CntlrReqResult, error)
	CntlrGetInfo(context.Context, *None) (*CntlrInfoData, error)
	CntlrStream(*CntlrGetParamReqData, CntlrGrpc_CntlrStreamServer) error
	mustEmbedUnimplementedCntlrGrpcServer()
}

// UnimplementedCntlrGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedCntlrGrpcServer struct {
}

func (UnimplementedCntlrGrpcServer) CntlrGetParamReq(context.Context, *CntlrGetParamReqData) (*CntlrReqResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrGetParamReq not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrSetParamReq(context.Context, *CntlrSetParamReqData) (*CntlrSetParamResData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrSetParamReq not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrGetInstancesReq(context.Context, *CntlrGetInstancesReqData) (*CntlrReqResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrGetInstancesReq not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrAddInstanceReq(context.Context, *CntlrAddInstanceReqData) (*CntlrAddInstanceResData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrAddInstanceReq not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrOperateReq(context.Context, *CntlrOperateReqData) (*CntlrOperateResData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrOperateReq not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrGetDatamodelReq(context.Context, *CntlrGetDatamodelReqData) (*CntlrReqResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrGetDatamodelReq not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrDeleteInstanceReq(context.Context, *CntlrDeleteInstanceReqData) (*CntlrReqResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrDeleteInstanceReq not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrGetAgentMsgs(context.Context, *CntlrGetAgentMsgsData) (*CntlrReqResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrGetAgentMsgs not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrGetInfo(context.Context, *None) (*CntlrInfoData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CntlrGetInfo not implemented")
}
func (UnimplementedCntlrGrpcServer) CntlrStream(*CntlrGetParamReqData, CntlrGrpc_CntlrStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CntlrStream not implemented")
}
func (UnimplementedCntlrGrpcServer) mustEmbedUnimplementedCntlrGrpcServer() {}

// UnsafeCntlrGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CntlrGrpcServer will
// result in compilation errors.
type UnsafeCntlrGrpcServer interface {
	mustEmbedUnimplementedCntlrGrpcServer()
}

func RegisterCntlrGrpcServer(s grpc.ServiceRegistrar, srv CntlrGrpcServer) {
	s.RegisterService(&CntlrGrpc_ServiceDesc, srv)
}

func _CntlrGrpc_CntlrGetParamReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrGetParamReqData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrGetParamReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrGetParamReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrGetParamReq(ctx, req.(*CntlrGetParamReqData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrSetParamReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrSetParamReqData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrSetParamReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrSetParamReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrSetParamReq(ctx, req.(*CntlrSetParamReqData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrGetInstancesReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrGetInstancesReqData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrGetInstancesReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrGetInstancesReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrGetInstancesReq(ctx, req.(*CntlrGetInstancesReqData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrAddInstanceReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrAddInstanceReqData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrAddInstanceReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrAddInstanceReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrAddInstanceReq(ctx, req.(*CntlrAddInstanceReqData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrOperateReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrOperateReqData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrOperateReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrOperateReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrOperateReq(ctx, req.(*CntlrOperateReqData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrGetDatamodelReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrGetDatamodelReqData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrGetDatamodelReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrGetDatamodelReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrGetDatamodelReq(ctx, req.(*CntlrGetDatamodelReqData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrDeleteInstanceReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrDeleteInstanceReqData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrDeleteInstanceReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrDeleteInstanceReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrDeleteInstanceReq(ctx, req.(*CntlrDeleteInstanceReqData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrGetAgentMsgs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CntlrGetAgentMsgsData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrGetAgentMsgs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrGetAgentMsgs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrGetAgentMsgs(ctx, req.(*CntlrGetAgentMsgsData))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrGetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CntlrGrpcServer).CntlrGetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cntrlgrpc.CntlrGrpc/CntlrGetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CntlrGrpcServer).CntlrGetInfo(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

func _CntlrGrpc_CntlrStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CntlrGetParamReqData)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CntlrGrpcServer).CntlrStream(m, &cntlrGrpcCntlrStreamServer{stream})
}

type CntlrGrpc_CntlrStreamServer interface {
	Send(*CntlrReqResult) error
	grpc.ServerStream
}

type cntlrGrpcCntlrStreamServer struct {
	grpc.ServerStream
}

func (x *cntlrGrpcCntlrStreamServer) Send(m *CntlrReqResult) error {
	return x.ServerStream.SendMsg(m)
}

// CntlrGrpc_ServiceDesc is the grpc.ServiceDesc for CntlrGrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CntlrGrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cntrlgrpc.CntlrGrpc",
	HandlerType: (*CntlrGrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CntlrGetParamReq",
			Handler:    _CntlrGrpc_CntlrGetParamReq_Handler,
		},
		{
			MethodName: "CntlrSetParamReq",
			Handler:    _CntlrGrpc_CntlrSetParamReq_Handler,
		},
		{
			MethodName: "CntlrGetInstancesReq",
			Handler:    _CntlrGrpc_CntlrGetInstancesReq_Handler,
		},
		{
			MethodName: "CntlrAddInstanceReq",
			Handler:    _CntlrGrpc_CntlrAddInstanceReq_Handler,
		},
		{
			MethodName: "CntlrOperateReq",
			Handler:    _CntlrGrpc_CntlrOperateReq_Handler,
		},
		{
			MethodName: "CntlrGetDatamodelReq",
			Handler:    _CntlrGrpc_CntlrGetDatamodelReq_Handler,
		},
		{
			MethodName: "CntlrDeleteInstanceReq",
			Handler:    _CntlrGrpc_CntlrDeleteInstanceReq_Handler,
		},
		{
			MethodName: "CntlrGetAgentMsgs",
			Handler:    _CntlrGrpc_CntlrGetAgentMsgs_Handler,
		},
		{
			MethodName: "CntlrGetInfo",
			Handler:    _CntlrGrpc_CntlrGetInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CntlrStream",
			Handler:       _CntlrGrpc_CntlrStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "cntlrgrpc.proto",
}
