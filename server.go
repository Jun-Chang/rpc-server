package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"

	libhttp "github.com/CyberAgent/car-golib/http"
	"github.com/CyberAgent/car-golib/log"
	"github.com/Jun-Chang/rpc-server/codec"
	"github.com/Jun-Chang/rpc-server/proto"
	"github.com/Jun-Chang/rpc-server/service"
	"github.com/julienschmidt/httprouter"
	msgpackrpc "github.com/msgpack-rpc/msgpack-rpc-go/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// grpc
type grpcService struct {
}

func (ts *grpcService) Call(context.Context, *proto.RequestType) (*proto.Response, error) {
	seq := service.Run()

	return &proto.Response{Seq: int32(seq)}, nil
}

// messagepack rpc
type messagepackResolver map[string]reflect.Value

func (rslvr messagepackResolver) Resolve(name string, arguments []reflect.Value) (reflect.Value, error) {
	return rslvr[name], nil
}

func MessagepackCall() (map[string]interface{}, fmt.Stringer) {
	res := map[string]interface{}{}
	seq := service.Run()
	res["seq"] = seq

	return res, nil
}

// http
func HttpCall(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	seq := service.Run()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("%d", seq)))
}

func main() {
	mode := os.Getenv("RPC_MODE")
	log.Info("MODE", mode)

	switch mode {
	case "grpc":
		grpcL, err := net.Listen("tcp", ":11111")
		if err != nil {
			panic(err)
		}
		server := grpc.NewServer(grpc.CustomCodec(codec.CodecJson{}))

		proto.RegisterTestServiceServer(server, new(grpcService))
		server.Serve(grpcL)
	case "messagepack":
		rslvr := messagepackResolver{"call": reflect.ValueOf(MessagepackCall)}
		server := msgpackrpc.NewServer(rslvr, true, nil)
		lis, err := net.Listen("tcp", ":11112")
		if err != nil {
			panic(err)
		}
		server.Listen(lis)
		server.Run()
	case "http":
		router := httprouter.New()
		router.GET("/", HttpCall)

		sock := os.Getenv("HTTP_SOCK")
		libhttp.ListenUnixSocket(sock, router)
	default:
		panic("invalid mode. " + mode)
	}
}
