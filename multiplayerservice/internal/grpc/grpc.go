package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	proto "github.com/horahoradev/nexus2/multiplayerservice/protocol"
	"context"
	log "github.com/sirupsen/logrus"
	guuid "github.com/google/uuid"
)


var _ proto.MultiplayerServiceServer = (*GrpcServer)(nil)

type GrpcServer struct {
	playerLocMap map[string]string
	pubsubHelper MapPubsubManager
	proto.UnimplementedMultiplayerServiceServer
}

func New() (GrpcServer, error){
	return GrpcServer{
		playerLocMap: make(map[string]string),
		pubsubHelper: NewPubsubManager(),
	}, nil
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
		g.playerLocMap[clientUUID] =  defaultMapID


		clientMsg, err := stream.Recv()
		if err != nil {
			log.Errorf("recv err: %s", err))
		}

		switch msg := clientMsg.Payload.(type) {
		case *proto.ClientMessage_Movemsg:
			err = g.pubsubHelper.PublishMove(moveMsg{
				x:          msg.Movemsg.X,
				y:          msg.Movemsg.Y,
				playerUUID: clientUUID,
			}, g.playerLocMap[clientUUID])
			if err != nil {
				log.Errorf("could not publish move. Err: %s", err)
			}

		case *proto.ClientMessage_Navigatemsg:
			oldMapID, ok := g.playerLocMap[clientUUID]
			g.playerLocMap[clientUUID] =  msg.Navigatemsg.MapID
			if ok {
				err = g.pubsubHelper.Subscribe(clientUUID, g.playerLocMap[clientUUID], &oldMapID )
				if err != nil {
					log.Errorf("Could not subscribe. Err: %s", err)
				}
			} else {
				g.pubsubHelper.Subscribe(clientUUID, g.playerLocMap[clientUUID], nil )
				if err != nil {
					log.Errorf("Could not subscribe. Err: %s", err)
				}
			}

			resp := proto.ServerMessage{Payload:
				&proto.ServerMessage_Navigateresp{
					Navigateresp: &proto.ServerNavigate{
						Maploc:   fmt.Sprintf("./%s.tmx", msg.Navigatemsg.MapID),
						Audioloc: "", // TODO
						Players:  nil,
					},
				},
			}
			
			err = stream.Send(&resp)
			if err != nil {
				log.Errorf("Error sending navigate resp: %s", err)
			}

		case *proto.ClientMessage_Chatmsg:
			err = g.pubsubHelper.PublishChat(chatMsg{
				message: msg.Chatmsg.Message,
				playerUUID: clientUUID,
			}, g.playerLocMap[clientUUID])
			if err != nil {
				log.Errorf("could not publish chat message. Err: %s", err)
			}
		default:
			log.Errorf("Unknown client message!")
		}
	}
}


type MapPubsubManager struct {
	// Map ID to player UUID
	playerChanMap map[string]map[string]playerChans
}

func NewPubsubManager() MapPubsubManager{
	return MapPubsubManager{playerChanMap:make(map[string]map[string]playerChans),
	}
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

func (m *MapPubsubManager) Subscribe(playerUUID, mapID string, oldMapID *string) error {
	_, ok := m.playerChanMap[defaultMapID]
	if !ok {
		m.playerChanMap[defaultMapID] = make(map[string]playerChans)
	}

	err := m.movePlayerChanInfo(playerUUID, mapID, oldMapID)
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

func (m *MapPubsubManager) movePlayerChanInfo(playerUUID, mapID string, oldmapID *string) error {
	if oldmapID == nil {
		// Player doesn't have a location, new login?

		newChans :=  playerChans{
			moveChan: make(chan moveMsg),
			chatChan: make(chan chatMsg),
		}
		m.playerChanMap[mapID][playerUUID] = newChans

		// Return early, no move to perform
		return nil
	}

	val := m.playerChanMap[*oldmapID][playerUUID]
	// Move to new location
	delete(m.playerChanMap[*oldmapID], playerUUID)
	m.playerChanMap[mapID][playerUUID] = val

	return nil
}

