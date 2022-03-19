package models

import "server/pkg/grpcService"

type Player struct {
	GrpcInfo            *grpcService.ClientInfoRequest
	Stream              grpcService.Greeter_ConnectToTheHobbyServer
	NotificationChannel chan Notification
	SessionID           int
	StatusName          string
}



//Status: KILL, LIVE

type Notification struct {
	Action     string
	NamePlayer string
}

type Session struct {
	ID     int
	Mafia  []Player
	Police []Player
	People []Player
}
