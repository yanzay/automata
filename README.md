Automata
========

Makes interaction with Arduino fun and easy.

Lightweight and really simple alternative to firmata. Now supports digital writes and reads from Dallas Temperature sensors.

Usage:
```go
package main

import (
  "github.com/tarm/goserial"
  "github.com/yanzay/automata"
  "fmt"
)

func main() {
  c := &serial.Config{Name: "/dev/tty.usbmodem1411", Baud: 57600}
  conn, err := serial.OpenPort(c)
  if err != nil {
    fmt.Println("Error connecting to arduino")
    return
  }
  ar := automata.NewArduino(conn)
  ar.SetDigitalOutput(13)
  for {
    ar.Toggle(13)
    temp := ar.Temp()
    fmt.Printf("Temp: %d C\n", temp)
  }
}
```
