package main

import "fmt"
import "net"
import "time"
import "encoding/binary"
const NEURAL_PORT = 52192

type neuralController struct {
  conn net.Conn
}
type neuralMessage struct {
  messageType byte
  length int
}

const ( 
  MsgHelloServer = 1
  MsgHelloClient = 2
  MsgWorm = 3
  MsgOk = 4
  MsgNok = 5
  MsgSense = 6
  MsgMove = 7
  MsgKill = 8
  MsgBye = 9
)
  
func (c *neuralController) getMessage() (*neuralMessage){
  message := new(neuralMessage)
  header := make([]byte, 5)
  c.conn.Read(header)
  message.messageType = header[0]
  message.length = int(binary.LittleEndian.Uint16(header[1:]))
  return message
}

func (m *neuralMessage) getData(conn net.Conn) ([]byte) {
  result := make([]byte, m.length)
  conn.Read(result)
  return result
}

func (c *neuralController) tryConnect(hostname string) (bool){
  fmt.Printf("Connecting to %s\n", hostname)
  address := fmt.Sprintf("%s:%d", hostname, NEURAL_PORT)
  conn, err := net.Dial("tcp", address)
  if err != nil {
  	fmt.Printf("Connection to %s failed.\n", address)
    return false
  }
  c.conn = conn
  fmt.Printf("Connected to %s\n", address)
  return true
}

func (c *neuralController) run() {
  for {
    message := c.getMessage()
    message.getData(c.conn)
    fmt.Printf("Message received of type %d, length %d\n", message.messageType, message.length)
  }
}

func (c *neuralController) start(hostname string) {
  fmt.Printf("NeuralController 0.0.1 starting...\n")
  for {
    connected := false
    for connected == false {
      connected = c.tryConnect(hostname)
      if (connected == false) {time.Sleep(1 * time.Second)}
    }
    c.run()
  }
}