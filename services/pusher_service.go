package services

import (
	"log"

	"github.com/pusher/pusher-http-go/v5"
	"github.com/rogerjeasy/go-letusconnect/config"
)

// PusherClient is the instance of the Pusher client
var PusherClient *pusher.Client

// InitializePusher sets up the Pusher client with environment variables
func InitializePusher() {
	PusherClient = &pusher.Client{
		AppID:   config.PusherAppID,
		Key:     config.PusherKey,
		Secret:  config.PusherSecret,
		Cluster: config.PusherCluster,
		Secure:  true,
	}

	log.Println("Pusher client initialized successfully")
}
