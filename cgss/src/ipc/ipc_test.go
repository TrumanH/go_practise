package ipc

import (
    "fmt"
    "testing"
)

type EchoServer struct {
}

func (server *EchoServer) Handle(method, params string) *Response {
    return &Response{Code:"200", Body: fmt.Sprintf("method: %v, params: %v", method, params)}
}

func (server *EchoServer) Name() string {
    return "EchoServer"
}

func TestIpc(t *testing.T) {
    server := NewIpcServer(&EchoServer{})
    
    client1 := NewIpcClient(server)
    client2 := NewIpcClient(server)
    
    resp1, error1 := client1.Call("login", "From Client1")
    resp2, error2 := client2.Call("", "From Client2")
    if resp1.Body != "200" || error1 != nil {
        t.Error("IpcClient.Call failed. resp1:", resp1, "resp2:", resp2)
    }
    if resp2.Body != "200" || error2 != nil {
        t.Error("IpcClient.Call failed. resp1:", resp1, "resp2:", resp2)
    }
    client1.Close()
    client2.Close()
}
