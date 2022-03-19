package usecase

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"server/config"
	"server/internal/kernel"
	"server/internal/models"
	"server/pkg/grpcService"
	"server/pkg/serverFlow"
	"time"
)

type kernelUS struct {
	cfg     *config.Config
	players []models.Player
	games   map[int]models.Session
}

func NewKernelUS(cfg *config.Config) kernel.UseCase {
	var players []models.Player
	games := make(map[int]models.Session)
	return &kernelUS{
		cfg:     cfg,
		players: players,
		games:   games,
	}
}

func (k *kernelUS) GetConnectedPlayers() *grpcService.ClientsResponse {
	var clients grpcService.ClientsResponse
	for _, r := range k.players {
		clients.Client = append(clients.Client, r.GrpcInfo)
	}
	return &clients
}

func (k *kernelUS) CheckStartGame() (ok bool, players []models.Player) {
	lenPlayers := 0
	for _, r := range k.players {
		if r.SessionID == 0 {
			lenPlayers += 1
			players = append(players, r)
		}
	}
	if lenPlayers > k.cfg.Server.LenPlayersForStart {
		return true, players
	} else {
		return false, players
	}
}

func (k *kernelUS) GetLastIDGame() int {
	lastID := 0
	for _, r := range k.players {
		if r.SessionID > lastID {
		}
		lastID = r.SessionID
	}
	return lastID
}

func (k *kernelUS) StartingGameSession(players []models.Player) {
	for _, r := range players {
		go func(channel chan models.Notification) {
			channel <- models.Notification{Action: "Starting the game session"}
		}(r.NotificationChannel)
	}
	rand.Seed(time.Now().UnixNano())
	newIDGame := k.GetLastIDGame() + 1
	var newSession models.Session
	k.games[newIDGame] = newSession
	for i, r := range k.players {
		for _, l := range players {
			if r.GrpcInfo.Name == l.GrpcInfo.Name {
				k.players[i].SessionID = newIDGame
				k.players[i].StatusName = LiveStatus
			}
		}
	}

	mafia := []models.Player{players[0]}
	people := []models.Player{players[1], players[2]}
	police := []models.Player{players[3]}

	newSession.ID = newIDGame
	newSession.People = people
	newSession.Mafia = mafia
	newSession.Police = police

	for _, r := range mafia {
		go func(channel chan models.Notification) {
			channel <- models.Notification{Action: "You are mafia"}
		}(r.NotificationChannel)
	}
	for _, r := range police {
		go func(channel chan models.Notification) {
			channel <- models.Notification{Action: "You are police"}
		}(r.NotificationChannel)
	}
	for _, r := range people {
		go func(channel chan models.Notification) {
			channel <- models.Notification{Action: "You are people"}
		}(r.NotificationChannel)
	}
}

func (k *kernelUS) AddPlayerToTheHobby(player *grpcService.ClientInfoRequest, stream grpcService.Greeter_ConnectToTheHobbyServer) (chan models.Notification, error) {
	for _, r := range k.players {
		if r.GrpcInfo.Name == player.Name {
			return nil, errors.New("Such player name already exists")
		}
	}

	fmt.Println("New client", player.Name, "has been connected to our game")
	channel := make(chan models.Notification)
	playerr := models.Player{GrpcInfo: player, Stream: stream, NotificationChannel: channel, SessionID: 0}
	k.players = append(k.players, playerr)
	for _, r := range k.players {
		go func(channel chan models.Notification) {
			channel <- models.Notification{NamePlayer: player.Name, Action: "Enter the lobby"}
		}(r.NotificationChannel)
	}
	ok, players := k.CheckStartGame()
	if ok {
		k.StartingGameSession(players)
	}

	return channel, nil
}

func (k *kernelUS) RemovePlayerFromTheHobby(player *grpcService.ClientInfoRequest) error {
	for i, r := range k.players {
		if r.GrpcInfo.Name == player.Name {
			k.players = append(k.players[:i], k.players[i+1:]...)
			fmt.Println("Client", player.Name, "has been disconnected from our game")
			for _, r := range k.players {
				go func(channel chan models.Notification) {
					channel <- models.Notification{NamePlayer: player.Name, Action: "Exited the lobby"}
				}(r.NotificationChannel)
			}
			return nil
		}
	}
	return errors.New("Cannot find such player")
}

func (k *kernelUS) StartListening() {
	fmt.Println(k.cfg.Server.IP)
	lis, err := net.Listen("tcp", k.cfg.Server.IP)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server has been successfully started")
	s := grpc.NewServer()
	greeterServer := serverFlow.NewGreeterServer(k)
	grpcService.RegisterGreeterServer(s, &greeterServer)
	if err := s.Serve(lis); err != nil {
		fmt.Println(err)
	}
}

func (k *kernelUS) PushAllClients(notification *grpcService.Notification) {
	for _, r := range k.players {
		err := r.Stream.Send(notification)
		if err != nil {
			fmt.Println(err)
		}

	}
}
