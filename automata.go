package automata

import (
  "github.com/tarm/goserial"
  "io"
  "sync"
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

type Arduino struct {
  conn  io.ReadWriteCloser
  pins  map[byte]bool
  mutex *sync.Mutex
}

func NewArduino(conn io.ReadWriteCloser) *Arduino {
  ar := new(Arduino)
  ar.conn = conn
  ar.pins = make(map[byte]bool)
  time.Sleep(2 * time.Second)
  ar.mutex = &sync.Mutex{}
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
  buf := make([]byte, 1)
  ar.mutex.Lock()
  ar.conn.Write([]byte{command, parameter})
  ar.conn.Read(buf)
  ar.mutex.Unlock()
  return buf[0]
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
