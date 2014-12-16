package automata

import (
  "github.com/tarm/goserial"
  "io"
  "time"
)

const (
  Ping             = 0
  SetDigitalOutput = 1
  DigitalWriteHigh = 2
  DigitalWriteLow  = 3
  GetTemp          = 4

  RespOK = 0
)

type Message struct {
  Command   byte
  Parameter byte
}

type Arduino struct {
  conn      io.ReadWriteCloser
  pins      map[byte]bool
  messages  chan Message
  responses chan byte
}

func NewArduino(conn io.ReadWriteCloser) *Arduino {
  ar := new(Arduino)
  ar.conn = conn
  ar.pins = make(map[byte]bool)
  ar.messages = make(chan Message, 16)
  ar.responses = make(chan byte)
  time.Sleep(2 * time.Second)
  go ar.messageHandler()
  ar.Ping()
  return ar
}

func NewSerial(port string) (*Arduino, error) {
  c := &serial.Config{Name: port, Baud: 57600}
  conn, err := serial.OpenPort(c)
  if err != nil {
    return nil, err
  }
  ar := NewArduino(conn)
  return ar, nil
}

func (ar *Arduino) sendCommand(command byte, parameter byte) byte {
  ar.messages <- Message{Command: command, Parameter: parameter}
  response := <-ar.responses
  return response
}

func (ar *Arduino) Ping() {
  resp := ar.sendCommand(Ping, 0)
  a := make([]byte, 1)
  for resp != RespOK {
    ar.conn.Read(a)
    resp = a[0]
  }
}

func (ar *Arduino) SetDigitalOutput(pin byte) {
  ar.sendCommand(SetDigitalOutput, pin)
}

func (ar *Arduino) On(pin byte) {
  ar.sendCommand(DigitalWriteHigh, pin)
  ar.pins[pin] = true
}

func (ar *Arduino) Off(pin byte) {
  ar.sendCommand(DigitalWriteLow, pin)
  ar.pins[pin] = false
}

func (ar *Arduino) Toggle(pin byte) {
  if ar.pins[pin] == true {
    ar.Off(pin)
  } else {
    ar.On(pin)
  }
}

func (ar *Arduino) Temp() byte {
  return ar.sendCommand(GetTemp, 0)
}

func (ar *Arduino) messageHandler() {
  for {
    message := <-ar.messages
    buf := make([]byte, 1)
    ar.conn.Write([]byte{message.Command, message.Parameter})
    ar.conn.Read(buf)
    ar.responses <- buf[0]
  }
}
