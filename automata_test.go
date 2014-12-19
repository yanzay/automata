package automata

import (
  "testing"
  "time"
)

func TestNewSerial(t *testing.T) {
  ar, err := NewSerial("/dev/invalid")
  if ar != nil || err == nil {
    t.Error("Should not create arduino for invalid port")
  }
}

func TestNewEthernet(t *testing.T) {
  _, err := NewEthernet("169.254.68.110:13666")
  if err != nil {
    t.Error("Should create arduino via ethernet")
  }
}

func TestSendCommand(t *testing.T) {
  // ar, _ := NewSerial("/dev/tty.usbmodem1421")
  ar, _ := NewEthernet("169.254.68.110:13666")
  ar.SetDigitalOutput(12)
  for i := 0; i < 10; i++ {
    go func() {
      resp := ar.On(12)
      if resp != DigitalWriteHigh {
        t.Error("On should respond with DigitalWriteHigh")
      }
    }()

    go func() {
      temp := ar.Temp()
      t.Log(temp)
      if temp[2] == 0 && temp[3] == 0 {
        t.Errorf("Temp should not be 0, got %v", temp)
      }
    }()

    go func() {
      resp := ar.Off(12)
      if resp != DigitalWriteLow {
        t.Error("Off should respond with DigitalWriteHigh")
      }
    }()

    go func() {
      ar.Toggle(12)
    }()
  }
  time.Sleep(4 * time.Second)
}

func TestToggle(t *testing.T) {
  ar, _ := NewSerial("/dev/tty.usbmodem1411")
  ar.SetDigitalOutput(12)
  resp := ar.Toggle(12)
  if resp != DigitalWriteHigh {
    t.Error("Toggle should write digital high first")
  }
  resp = ar.Toggle(12)
  if resp != DigitalWriteLow {
    t.Error("Toggle should write digital low second")
  }
}
