package content

import (
	"context"
	"fmt"
	"log"
	"net"
	"server/conf"
	database "server/db"
	"server/proto"
	"time"

	"google.golang.org/grpc"
)

type EchoSrv struct {
	proto.UnimplementedEchoServer
}

func (e *EchoSrv) Send(ctx context.Context, req *proto.EchoReq) (res *proto.EchoRes, err error) {
	fmt.Printf("receive client request, client send:%s\n", req.Datetime)
	name := database.GetTestingSQLService()
	return &proto.EchoRes{
		Name:     name,
		Datetime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

type Grpc struct {
	config conf.Conf
	name   string
}

func (g *Grpc) Start() {
	apiListener, err := net.Listen("tcp", g.config.Proto.Port)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("grpc server start.")
	fmt.Println(g.config.Proto.Port)

	es := &EchoSrv{}
	grpc := grpc.NewServer()
	proto.RegisterEchoServer(grpc, es)
	if err := grpc.Serve(apiListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (g *Grpc) GetName() string {
	return g.name
}

func (g *Grpc) SetConf(in conf.Conf) {
	g.config = in
}

func init() {
	GetFactory().Add(&Grpc{
		name: "grpc",
	})
}
