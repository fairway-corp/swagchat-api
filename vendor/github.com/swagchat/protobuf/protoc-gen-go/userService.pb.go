// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: userService.proto

package protoc_gen_go

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UserService service

type UserServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error)
	RetrieveUsers(ctx context.Context, in *RetrieveUsersRequest, opts ...grpc.CallOption) (*UsersResponse, error)
	RetrieveUser(ctx context.Context, in *RetrieveUserRequest, opts ...grpc.CallOption) (*User, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	RetrieveUserRooms(ctx context.Context, in *RetrieveUserRoomsRequest, opts ...grpc.CallOption) (*UserRoomsResponse, error)
	RetrieveContacts(ctx context.Context, in *RetrieveContactsRequest, opts ...grpc.CallOption) (*UsersResponse, error)
	RetrieveProfile(ctx context.Context, in *RetrieveProfileRequest, opts ...grpc.CallOption) (*User, error)
	// rpc RetrieveDeviceUsers (RetrieveDeviceUsersRequest) returns (DeviceUsersResponse) {}
	RetrieveRoleUsers(ctx context.Context, in *RetrieveRoleUsersRequest, opts ...grpc.CallOption) (*RoleUsersResponse, error)
}

type userServiceClient struct {
	cc *grpc.ClientConn
}

func NewUserServiceClient(cc *grpc.ClientConn) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/CreateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RetrieveUsers(ctx context.Context, in *RetrieveUsersRequest, opts ...grpc.CallOption) (*UsersResponse, error) {
	out := new(UsersResponse)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/RetrieveUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RetrieveUser(ctx context.Context, in *RetrieveUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/RetrieveUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/UpdateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/DeleteUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RetrieveUserRooms(ctx context.Context, in *RetrieveUserRoomsRequest, opts ...grpc.CallOption) (*UserRoomsResponse, error) {
	out := new(UserRoomsResponse)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/RetrieveUserRooms", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RetrieveContacts(ctx context.Context, in *RetrieveContactsRequest, opts ...grpc.CallOption) (*UsersResponse, error) {
	out := new(UsersResponse)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/RetrieveContacts", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RetrieveProfile(ctx context.Context, in *RetrieveProfileRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/RetrieveProfile", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RetrieveRoleUsers(ctx context.Context, in *RetrieveRoleUsersRequest, opts ...grpc.CallOption) (*RoleUsersResponse, error) {
	out := new(RoleUsersResponse)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserService/RetrieveRoleUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	RetrieveUsers(context.Context, *RetrieveUsersRequest) (*UsersResponse, error)
	RetrieveUser(context.Context, *RetrieveUserRequest) (*User, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*User, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*google_protobuf1.Empty, error)
	RetrieveUserRooms(context.Context, *RetrieveUserRoomsRequest) (*UserRoomsResponse, error)
	RetrieveContacts(context.Context, *RetrieveContactsRequest) (*UsersResponse, error)
	RetrieveProfile(context.Context, *RetrieveProfileRequest) (*User, error)
	// rpc RetrieveDeviceUsers (RetrieveDeviceUsersRequest) returns (DeviceUsersResponse) {}
	RetrieveRoleUsers(context.Context, *RetrieveRoleUsersRequest) (*RoleUsersResponse, error)
}

func RegisterUserServiceServer(s *grpc.Server, srv UserServiceServer) {
	s.RegisterService(&_UserService_serviceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RetrieveUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RetrieveUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/RetrieveUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RetrieveUsers(ctx, req.(*RetrieveUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RetrieveUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RetrieveUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/RetrieveUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RetrieveUser(ctx, req.(*RetrieveUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RetrieveUserRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveUserRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RetrieveUserRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/RetrieveUserRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RetrieveUserRooms(ctx, req.(*RetrieveUserRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RetrieveContacts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveContactsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RetrieveContacts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/RetrieveContacts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RetrieveContacts(ctx, req.(*RetrieveContactsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RetrieveProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RetrieveProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/RetrieveProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RetrieveProfile(ctx, req.(*RetrieveProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RetrieveRoleUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveRoleUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RetrieveRoleUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserService/RetrieveRoleUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RetrieveRoleUsers(ctx, req.(*RetrieveRoleUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "swagchat.protobuf.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "RetrieveUsers",
			Handler:    _UserService_RetrieveUsers_Handler,
		},
		{
			MethodName: "RetrieveUser",
			Handler:    _UserService_RetrieveUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserService_DeleteUser_Handler,
		},
		{
			MethodName: "RetrieveUserRooms",
			Handler:    _UserService_RetrieveUserRooms_Handler,
		},
		{
			MethodName: "RetrieveContacts",
			Handler:    _UserService_RetrieveContacts_Handler,
		},
		{
			MethodName: "RetrieveProfile",
			Handler:    _UserService_RetrieveProfile_Handler,
		},
		{
			MethodName: "RetrieveRoleUsers",
			Handler:    _UserService_RetrieveRoleUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "userService.proto",
}

func init() { proto.RegisterFile("userService.proto", fileDescriptorUserService) }

var fileDescriptorUserService = []byte{
	// 346 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xc1, 0x4e, 0xfa, 0x40,
	0x10, 0xc6, 0x39, 0x71, 0x98, 0xff, 0xdf, 0x28, 0x7b, 0xd0, 0x04, 0x3d, 0x18, 0x42, 0x34, 0xa2,
	0xb4, 0x89, 0xbe, 0x81, 0xe8, 0x4d, 0x12, 0xad, 0x21, 0x26, 0xc6, 0xcb, 0xb6, 0x0e, 0x4b, 0x93,
	0xd2, 0xa9, 0x9d, 0x2d, 0xc6, 0x67, 0xf3, 0xe5, 0x4c, 0x4b, 0x97, 0x85, 0x08, 0x0b, 0xde, 0xba,
	0xfd, 0x7e, 0xf3, 0x7d, 0x3b, 0x33, 0x2d, 0xb4, 0x0a, 0xc6, 0xfc, 0x19, 0xf3, 0x59, 0x1c, 0xa1,
	0x97, 0xe5, 0xa4, 0x49, 0xb4, 0xf8, 0x53, 0xaa, 0x68, 0x22, 0xf5, 0xfc, 0x1c, 0x16, 0xe3, 0xf6,
	0x89, 0x22, 0x52, 0x09, 0xfa, 0x32, 0x8b, 0x7d, 0x99, 0xa6, 0xa4, 0xa5, 0x8e, 0x29, 0xe5, 0x39,
	0xd0, 0x3e, 0xae, 0x55, 0x83, 0xfb, 0x38, 0xcd, 0xf4, 0x57, 0x2d, 0x56, 0x01, 0x43, 0x64, 0x96,
	0xaa, 0x0e, 0xb8, 0xfe, 0x6e, 0xc2, 0xbf, 0x91, 0x8d, 0x15, 0x43, 0x80, 0x41, 0x8e, 0x52, 0x63,
	0xf9, 0x52, 0x74, 0xbd, 0x5f, 0xf9, 0x9e, 0x95, 0x03, 0xfc, 0x28, 0x90, 0x75, 0xfb, 0x68, 0x0d,
	0x55, 0xea, 0x9d, 0x86, 0x78, 0x83, 0xbd, 0x00, 0x75, 0x1e, 0xe3, 0xac, 0xaa, 0x60, 0x71, 0xbe,
	0x86, 0x5d, 0x21, 0x8c, 0xe9, 0xe9, 0x06, 0x53, 0x0e, 0x90, 0x33, 0x4a, 0x19, 0x3b, 0x0d, 0xf1,
	0x04, 0xff, 0x97, 0x6b, 0xc5, 0xd9, 0x16, 0xf3, 0x1d, 0x2e, 0x3c, 0x04, 0x18, 0x65, 0xef, 0xae,
	0xfe, 0xad, 0xbc, 0x83, 0xdd, 0x03, 0xc0, 0x1d, 0x26, 0xe8, 0xb0, 0xb3, 0xb2, 0xb1, 0x3b, 0xf4,
	0xe6, 0x3b, 0xb4, 0xcc, 0x7d, 0xb9, 0xc3, 0x4e, 0x43, 0x4c, 0xa0, 0xb5, 0xd2, 0x0e, 0xd1, 0x94,
	0xc5, 0xe5, 0xb6, 0xa6, 0x4b, 0xca, 0x78, 0x77, 0x37, 0x5c, 0xb5, 0x86, 0x16, 0x93, 0x0d, 0xe1,
	0xc0, 0x78, 0x0c, 0x28, 0xd5, 0x32, 0xd2, 0x2c, 0x7a, 0x8e, 0x20, 0x03, 0xfd, 0x65, 0x7b, 0x2f,
	0xb0, 0x6f, 0xca, 0x1f, 0x73, 0x1a, 0xc7, 0x09, 0x8a, 0x0b, 0x47, 0x44, 0xcd, 0xec, 0x30, 0xf4,
	0xa5, 0x31, 0x05, 0x94, 0xd4, 0x1f, 0x9e, 0x6b, 0x4c, 0x0b, 0xca, 0x35, 0xa6, 0x25, 0xc8, 0xb4,
	0x70, 0x7b, 0xf5, 0xda, 0x53, 0xb1, 0x9e, 0x14, 0xa1, 0x17, 0xd1, 0xd4, 0x37, 0x35, 0xf6, 0xe7,
	0xab, 0x1e, 0xa2, 0xbe, 0xc2, 0xb4, 0xaf, 0x28, 0x6c, 0x56, 0xc7, 0x9b, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x76, 0x28, 0xf5, 0x81, 0xe8, 0x03, 0x00, 0x00,
}