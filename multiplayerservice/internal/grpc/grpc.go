package grpc

import (
	"google.golang.org/grpc"
	proto "github.com/horahoradev/nexus2/multiplayerservice/protocol"
	"context"
	log "github.com/sirupsen/logrus"
	guuid "github.com/google/uuid"
)


var _ *proto.MultiplayerServiceServer = (*GrpcServer)(nil)

type GrpcServer struct {
}

func New() (GrpcServer, error){
	return GrpcServer{}, nil
}

const (
	defaultMapID = "default"
)

func (g *GrpcServer) Login(stream proto.MultiplayerService_LoginServer) error {
	for {


		// Client gets their own UUID
		uuid, err := guuid.NewUUID()
		if err != nil {
			log.Errorf("Could not generate uuid")
			return err
		}

		clientUUID := uuid.String()

		clientMsg, err := stream.Recv()
		if err != nil {
			log.Errorf("recv err: %s", err))
		}

		switch clientMsg.(type) {
		case proto.ClientMessage_Movemsg:


		case proto.ClientMessage_Navigatemsg:

		case proto.ClientMessage_Chatmsg:

		default:
			log.Errorf("Unknown client message!")
		}
	}
}


type MapPubsubManager struct {
	// Map ID to player UUID
	playerChanMap map[string]map[string]playerChans
	playerLocMap map[string]string
}

type playerChans struct {
	moveChan chan moveMsg
	chatChan chan chatMsg
}

type moveMsg struct {
	x, y       int64
	playerUUID string
}

type chatMsg struct {
	playerUUID, message string
}

func (m *MapPubsubManager) Subscribe(playerUUID, mapID string) error {
	_, ok := m.playerChanMap[defaultMapID]
	if !ok {
		m.playerChanMap[defaultMapID] = make(map[string]chan string)
	}

	err := m.movePlayerChanInfo(playerUUID, mapID)
	if err != nil {
		return err
	}

	return nil
}

func (m *MapPubsubManager) PublishMove(msg moveMsg, mapID string) error {
	for _, playerChans := range m.playerChanMap[mapID] {
		playerChans.moveChan <- msg
	}
}

func (m *MapPubsubManager) PublishChat(msg chatMsg, mapID string) error {
	for _, playerChans := range m.playerChanMap[mapID] {
		playerChans.chatChan <- msg
	}
}

func (m *MapPubsubManager) movePlayerChanInfo(playerUUID, mapID string) error {
	// Where is the player's current location?
	oldLoc, ok := m.playerLocMap[playerUUID]
	if !ok {
		// Player doesn't have a location, new login?
		m.playerLocMap[playerUUID] = mapID

		newChans :=  playerChans{
			moveChan: make(chan moveMsg),
			chatChan: make(chan chatMsg),
		}
		m.playerChanMap[mapID][playerUUID] = newChans
	}

	val := m.playerChanMap[oldLoc][playerUUID]
	// Move to new location
	delete(m.playerChanMap[oldLoc], playerUUID)
	m.playerChanMap[mapID][playerUUID] = val

	return nil
}

