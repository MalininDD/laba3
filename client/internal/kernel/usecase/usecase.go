package usecase

import (
	"client/config"
	"client/internal/kernel"
	"client/pkg/grpcService"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
)

type kernelUS struct {
	cfg *config.Config
}

func NewKernelUS(cfg *config.Config) kernel.UseCase{
	return &kernelUS{
		cfg: cfg,
	}
}

func (k *kernelUS) ConnectClient() {
	var name string
	fmt.Print("Enter your name = ")
	fmt.Scanf("%s\n", &name)
	conn, err := grpc.Dial(k.cfg.GrpcServer.Url, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	c := grpcService.NewGreeterClient(conn)

	in := emptypb.Empty{}
	palyers, err := c.ShowAllPlayersInTheHobby(context.Background(),&in)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Players before: ")
	fmt.Println(palyers)

	stream, err := c.ConnectToTheHobby(context.Background(), &grpcService.ClientInfoRequest{Name: name})
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream read failed: %v", err)
		}

		fmt.Println(msg)
	}

	//client := grpcService.NewGreeterClient(conn)
	//resp, err := client.ConnectToTheHobby(context.TODO(), &grpcService.ClientInfoRequest{Name: "ABOBA"})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("Hobby players: ")
	//for i, r := range resp.Client {
	//	fmt.Println(i+1, ")", r.Name)
	//}
	//fmt.Println(resp)

}
