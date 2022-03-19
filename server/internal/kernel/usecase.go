package kernel

import (
	"server/internal/models"
	"server/pkg/grpcService"
)

type UseCase interface {
	StartListening()
	AddPlayerToTheHobby(player *grpcService.ClientInfoRequest, stream grpcService.Greeter_ConnectToTheHobbyServer) (chan models.Notification, error)
	GetConnectedPlayers() *grpcService.ClientsResponse
	RemovePlayerFromTheHobby(player *grpcService.ClientInfoRequest) error
	PushAllClients(notification *grpcService.Notification)
	CheckStartGame() (ok bool, players []models.Player)
//	GetClientNotif(player *grpcService.ClientInfoRequest) *chan *grpcService.Notification
}
