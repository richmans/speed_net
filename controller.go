package main

import "fmt"
import "net"
import "time"
import "encoding/binary"
import "runtime"
const NEURAL_PORT = 52192

type neuralController struct {
  conn net.Conn
  runner networkRunner
}
type neuralMessage struct {
  messageType byte
  length int
  dataDone bool
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

func CreateShortMessage(messageType byte) (*neuralMessage){
  return CreateMessage(messageType, 0)
}

func CreateMessage(messageType byte, length int) (*neuralMessage){
  return &neuralMessage{messageType, length, false}
}

func (m *neuralMessage) getHeader() ([]byte) {
  header := make([]byte, 5)
  header[0] = m.messageType
  binary.LittleEndian.PutUint32(header[1:], uint32(m.length))
  return header
}
func (c *neuralController) getMessage() (*neuralMessage){
  message := new(neuralMessage)
  header := make([]byte, 5)
  c.conn.Read(header)
  message.messageType = header[0]
  message.length = int(binary.LittleEndian.Uint16(header[1:]))
  return message
}

func (m *neuralMessage) getData(conn net.Conn) ([]byte) {
  if(m.dataDone) { return nil }
  result := make([]byte, m.length)
  conn.Read(result)
  m.dataDone = true
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

func (c *neuralController) sendShortMessage(messageType byte) {
  c.sendMessage(messageType, nil)
}

func (c *neuralController) sendHello() {
  data := []byte(" NeuralController")
  data[0] = byte(len(data)) - 1
   c.sendMessage(MsgHelloClient, data)
}
func (c *neuralController) sendMessage(messageType byte, data []byte) {
  messageLength := 0
  if (data != nil) {
    messageLength = len(data)
  }
  message := CreateMessage(messageType, messageLength)
  c.conn.Write(message.getHeader())
  if (data != nil) {
    c.conn.Write(data)
  }
}

func (c *neuralController) createNetwork(dna []byte) {
  
}

func (c *neuralController) handleMessage(m *neuralMessage){
  if m.messageType == MsgHelloServer {
    c.sendHello()
  }else if m.messageType == MsgWorm {
    dna := m.getData(c.conn)
    c.createNetwork(dna)
  }
}

func (c *neuralController) run() {
  for {
    message := c.getMessage()
    c.handleMessage(message)
    message.getData(c.conn)
    fmt.Printf("Message received of type %d, length %d\n", message.messageType, message.length)
  }
}

func (c *neuralController) start(hostname string) {
  num_workers := 8
  runtime.GOMAXPROCS(num_workers)
  
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