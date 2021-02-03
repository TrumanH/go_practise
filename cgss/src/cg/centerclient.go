package cg

import (
	"errors"
	"encoding/json"
	"ipc"
)

type CenterClient struct {
	*ipc.IpcClient
}

func (client *CenterClient) AddPlayer(player *Player) error {
	b, err := json.Marshal(*player)
	if err != nil {
		return err
	}

	resp, err := client.Call("addplayer", string(b))
	// resp.Code == "200"
	if err == nil&&resp.Code == "200" {
		return nil
	}
	if err != nil {
		return nil
	}
	return errors.New(resp.Code)
}

func (client *CenterClient) RemovePlayer(name string) error {
	ret, _ := client.Call("removeplayer", name)
	if ret.Code == "200" {
		return nil
	}
	return errors.New(ret.Code)
}

func (client *CenterClient)ListPlayer(params string) (ps []*Player, err error) {
	resp, _ := client.Call("listplayer", params)
	if resp.Code != "200" {
		err := errors.New(resp.Code)
		return
	}

	err = json.Unmarshal([]byte(resp.Body), &ps)
	return
}

func (client *CenterClient)Broadcast(message string) error {
	m : = &Message{content:message}
	b, err := json.Marshal(m)
	fmt.Println("Serialized string b: ", b)
	resp, _ := client.Call("broadcast", string(b))
	if resp.Code == "200" {
		return nil
	}
	return errors.New(resp.Code)
}