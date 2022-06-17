package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var expectedOutput = `// Code generated by protoc-gen-go-kratos. DO NOT EDIT
// versions:
// protoc-gen-go-kratos v2.3.1

package tests

import (
	context "context"
	wrr "github.com/go-kratos/kratos/v2/selector/wrr"
	grpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)

type LibraryServiceGRPCClient struct {
	cli LibraryServiceClient
}

//NewLibraryServiceGRPCClient create grpc client for kratos
func NewLibraryServiceGRPCClient(ctx context.Context, opts ...grpc.ClientOption) (cli *LibraryServiceGRPCClient, err error) {
	opts = append(opts, grpc.WithBalancerName(wrr.Name))
	conn, err := grpc.DialInsecure(ctx, opts...)
	if err != nil {
		return nil, err
	}
	client := NewLibraryServiceClient(conn)
	return &LibraryServiceGRPCClient{cli: client}, nil
}

type BookServiceGRPCClient struct {
	cli BookServiceClient
}

//NewBookServiceGRPCClient create grpc client for kratos
func NewBookServiceGRPCClient(ctx context.Context, opts ...grpc.ClientOption) (cli *BookServiceGRPCClient, err error) {
	opts = append(opts, grpc.WithBalancerName(wrr.Name))
	conn, err := grpc.DialInsecure(ctx, opts...)
	if err != nil {
		return nil, err
	}
	client := NewBookServiceClient(conn)
	return &BookServiceGRPCClient{cli: client}, nil
}

type LibraryServiceV2GRPCClient struct {
	cli LibraryServiceV2Client
}

//NewLibraryServiceV2GRPCClient create grpc client for kratos
func NewLibraryServiceV2GRPCClient(ctx context.Context, opts ...grpc.ClientOption) (cli *LibraryServiceV2GRPCClient, err error) {
	opts = append(opts, grpc.WithBalancerName(wrr.Name))
	conn, err := grpc.DialInsecure(ctx, opts...)
	if err != nil {
		return nil, err
	}
	client := NewLibraryServiceV2Client(conn)
	return &LibraryServiceV2GRPCClient{cli: client}, nil
}

type LibraryServiceV3GRPCClient struct {
	cli LibraryServiceV3Client
}

//NewLibraryServiceV3GRPCClient create grpc client for kratos
func NewLibraryServiceV3GRPCClient(ctx context.Context, opts ...grpc.ClientOption) (cli *LibraryServiceV3GRPCClient, err error) {
	opts = append(opts, grpc.WithBalancerName(wrr.Name))
	conn, err := grpc.DialInsecure(ctx, opts...)
	if err != nil {
		return nil, err
	}
	client := NewLibraryServiceV3Client(conn)
	return &LibraryServiceV3GRPCClient{cli: client}, nil
}
`

func createPlugin(t *testing.T) *protogen.Plugin {
	file, err := ioutil.ReadFile("request_data.json")
	if err != nil {
		t.Fatal(err)
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := json.Unmarshal(file, req); err != nil {
		assert.Nil(t, err)
	}
	target, err := protogen.Options{}.New(req)
	assert.Nil(t, err)
	return target
}

func TestGenerate(t *testing.T) {
	gen := createPlugin(t)
	file := gen.FilesByPath["tests/library.proto"]
	generateFile(gen, file)
	response := gen.Response()
	content := response.File[0].Content
	assert.Equal(t, expectedOutput, *content)
}
