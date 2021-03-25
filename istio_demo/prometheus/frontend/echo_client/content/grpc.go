package content

import (
	"client/conf"
	"client/proto"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type Grpc struct {
	config conf.Conf
	name   string
}

func (g *Grpc) Start() {
	conn, err := grpc.Dial(g.config.Host, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Connet Fault:%v", err)
	}
	defer conn.Close()

	echoClient := proto.NewEchoClient(conn)

	for {
		receive, err := echoClient.Send(context.Background(), &proto.EchoReq{
			Datetime: time.Now().Format("2006-01-02 15:04:05"),
		})

		if err != nil {
			fmt.Printf("Receive Fault:%v", err)
		}
		fmt.Printf("%v:%v\n", receive.Name, receive.Datetime)
		time.Sleep(time.Duration(g.config.Interval) * time.Second)
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
