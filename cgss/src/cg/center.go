package cg

import (
	"encoding/json"
	"errors"
	"sync"

	"ipc"
)

// valid ipc implemented the 'CenterServer' interface
var _, ipc = &CenterServer{}

type Message struct {
	From 	string 	"from"
	To 		string  "to"
	Content string  "content"
}

type CenterServer struct {
	servers map[string] ipc.Server
	players []*Player
	rooms []*Room
	mutex sync.RWMutex
}

func NewCenterSever() *CenterServer {
	servers := make(map[string] ipc.Server)
	players := make([]*Player, 0)
	// rooms := make([]*Room, 0)

	return &CenterServer{servers:servers, players:players}
}

func (server *CenterServer)addPlayer(params string) error {
	player := NewPlayer()
	err := json.Unmarshal([]byte(params), &player)
	if err != nil {
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()

	append(server.players, player) // todo: optimise it with a unique elements structure
	return nil
}

func (server *CenterServer)removePlayer(params string) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	// todo: if used a unique elements structure, here don't need for loop
	for i, v := range server.players {
		if v.Name == params {
			if len(server.players) == 1 {
				server.players = make([]*Player, 0)
			}
			elseif i == len(server.players) - 1 {
				server.players = server.players[:i-1]
			} // in the end 
			elseif i == 0 {
				server.players = server.players[1:]
			} // at first
			else {
				server.players = append(server.players[:i-1], server.players[i+1:]...)
			} // at middle
		}
		return nil // found 
	}
	return errors.New("Player not found.")
}
// recieve a params but not use?
func (server *CenterServer) listPlayer(params string) (players string, err error) {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if len(server.players) > 0 {
		b, _ := json.Marshal(server.players)
		players = string(b)
	} else {
		err = errors.New("No player online")
	}
	return
}

func (server *CenterServer)broadcast(params string) error {
	var msg Message
	err := json.Unmarshal([]byte(params), &msg)
	if err != nil {
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()
	if len(server.players) > 0 {
		for _, p := range server.players {
			p.mq <- &msg
		}
	} else {
		err = errors.New("No player online.")
	}
	return err
}

func (server *CenterServer) Handle(method, params string) *ipc.Response {
	switch method {
	case "addplayer":
		err := server.addPlayer(params)
		if err != nil {
			return &ipc.Response{Code:err.Error()}
		}
	case "removeplayer":
		err := server.removePlayer(params)
		if err != nil {
			return &ipc.Response{Code:err.Error()}
		}
	case "listplayer":
		players, err := server.listPlayer(params)
		if err != nil {
			return &ipc.Response{Code:err.Error()}
		}
		return &ipc.Response("200", players)
	case "broadcast":
		players, err := server.broadcast(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
		return &ipc.Response{Code: "200"}
	default:
		return &ipc.Response{Code: "404", Body: method+":"+params}
	}
	return &ipc.Response{Code: "200"}

}

func (server *CenterServer)Name() string {
	return "CenterServer"
}