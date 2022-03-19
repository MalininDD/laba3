package serverFlow

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"server/internal/kernel"
	"server/internal/models"
	"server/pkg/grpcService"
)

type GreeterServer struct {
	kernelUS kernel.UseCase
	grpcService.UnimplementedGreeterServer
}

func NewGreeterServer(kernelUS kernel.UseCase) GreeterServer {
	return GreeterServer{
		kernelUS: kernelUS,
	}
}

//func (s *GreeterServer) ConnectToTheHobby(player *grpcService.ClientInfoRequest, stream grpcService.Greeter_ConnectToTheHobbyServer) error {
//	err := s.kernelUS.AddPlayerToTheHobby(player)
//	if err != nil {
//		return err
//	}
//	for {
//		fmt.Println(stream.Context())
//		err := stream.Send(s.kernelUS.GetConnectedPlayers())
//		if err != nil {
//			err := s.kernelUS.RemovePlayerFromTheHobby(player)
//			if err != nil {
//				fmt.Println(err)
//			}
//			return err
//		}
//		time.Sleep(2 * time.Second)
//	}
//}

func (s *GreeterServer) ShowAllPlayersInTheHobby(ctx context.Context, in *emptypb.Empty) (*grpcService.ClientsResponse, error) {
	return s.kernelUS.GetConnectedPlayers(), nil
}

func (s *GreeterServer) ConnectToTheHobby(player *grpcService.ClientInfoRequest, stream grpcService.Greeter_ConnectToTheHobbyServer) error {
	ch, err := s.kernelUS.AddPlayerToTheHobby(player, stream)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(stream.Context())
	go func(cancel context.CancelFunc, ch chan models.Notification) {
		for {
			if ctx.Err()!= nil {
				err := s.kernelUS.RemovePlayerFromTheHobby(player)
				if err != nil {
					fmt.Println(err)
				}
				cancel()
				<-ch
				return
			}
		}

	}(cancel, ch)
	//notif := s.kernelUS.GetClientNotif(player)
	for {
		value := <- ch
		//stream.Context()
		err := stream.Send(&grpcService.Notification{NamePlayer: value.NamePlayer, Action: value.Action})
		if err != nil {
			err := s.kernelUS.RemovePlayerFromTheHobby(player)
			if err != nil {
				return err
			}
		}


	}
	return nil
	//for {
	//	resp := grpcService.Notification{Action: "Acc", NamePlayer: player.Name}
	//	err := stream.Send(&resp)
	//	if err != nil {
	//		err := s.kernelUS.RemovePlayerFromTheHobby(player)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		return err
	//	}
	//
	//	time.Sleep(2 * time.Second)
	//}
}
