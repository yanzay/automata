package automata

import (
  "bufio"
  "errors"
  "github.com/op/go-logging"
  "github.com/tarm/goserial"
  "io"
  "net"
  "os"
  "time"
)

const (
  SerialArduino   = 0
  EthernetArduino = 1

  Ping             = 0
  SetDigitalOutput = 1
  DigitalWriteHigh = 2
  DigitalWriteLow  = 3
  GetTemp          = 4
)

type Message struct {
  Command   byte
  Parameter byte
}

type Arduino struct {
  conn      io.ReadWriteCloser
  pins      map[byte]bool
  messages  chan Message
  responses chan []byte
}

var log = logging.MustGetLogger("example")

func initLogger() {
  backend1 := logging.NewLogBackend(os.Stderr, "", 0)
  backend1Leveled := logging.AddModuleLevel(backend1)
  backend1Leveled.SetLevel(logging.ERROR, "")
  logging.SetBackend(backend1Leveled)
}

func New(t byte, addr string) (*Arduino, error) {
  var err error
  var conn io.ReadWriteCloser
  initLogger()
  switch t {
  case SerialArduino:
    c := &serial.Config{Name: addr, Baud: 57600}
    conn, err = serial.OpenPort(c)
  case EthernetArduino:
    conn, err = net.Dial("tcp", addr)
  default:
    return nil, errors.New("Undefined arduino connection type.")
  }
  if err != nil {
    log.Debug("Error connecting to arduino on %s, error: %v", addr, err)
    return nil, err
  }
  ar := newArduino(conn)
  return ar, nil
}

func newArduino(conn io.ReadWriteCloser) *Arduino {
  log.Debug("Instantiate arduino on connection %v", conn)
  ar := new(Arduino)
  log.Debug("Arduino isn't ready yet.")
  ar.conn = conn
  ar.pins = make(map[byte]bool)
  log.Debug("Launching messageHandler in goroutine")
  go ar.messageHandler()
  log.Debug("Pinging...")
  ar.Ping()
  return ar
}

func (ar *Arduino) Close() {
  ar.conn.Close()
}

func (ar *Arduino) sendCommand(command byte, parameter byte) []byte {
  ar.messages <- Message{Command: command, Parameter: parameter}
  response := <-ar.responses
  return response
}

func (ar *Arduino) Ping() {
  log.Debug("Sleeping for 2 seconds...")
  time.Sleep(2 * time.Second)
  ar.sendCommand(Ping, 0)
}

func (ar *Arduino) SetDigitalOutput(pin byte) {
  ar.sendCommand(SetDigitalOutput, pin)
}

func (ar *Arduino) On(pin byte) byte {
  ar.pins[pin] = true
  return ar.sendCommand(DigitalWriteHigh, pin)[3]
}

func (ar *Arduino) Off(pin byte) byte {
  ar.pins[pin] = false
  return ar.sendCommand(DigitalWriteLow, pin)[3]
}

func (ar *Arduino) Toggle(pin byte) byte {
  if ar.pins[pin] == true {
    return ar.Off(pin)
  }
  return ar.On(pin)
}

func (ar *Arduino) Temp() []byte {
  return ar.sendCommand(GetTemp, 0)
}

func (ar *Arduino) messageHandler() {
  log.Debug("Starting messageHandler")
  ar.messages = make(chan Message, 256)
  ar.responses = make(chan []byte)
  log.Debug("Channels initialized")

  for {
    log.Debug("Trying to get message from channel")
    message := <-ar.messages
    log.Debug("Message received from channel: %v", message)
    buf := make([]byte, 4)
    log.Debug("Writing message to connection %v", []byte{message.Command, message.Parameter})

    ar.conn.Write([]byte{message.Command, message.Parameter})
    log.Debug("Reading message from connection")
    reader := bufio.NewReader(ar.conn)

    buf[0], _ = reader.ReadByte()
    buf[1], _ = reader.ReadByte()
    buf[2], _ = reader.ReadByte()
    buf[3], _ = reader.ReadByte()

    log.Debug("Response from connection received, adding to channel %v", buf)
    ar.responses <- buf
  }
}
