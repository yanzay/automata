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
  ar := automata.New(automata.SerialArduino, "/dev/tty.usbmodem1411")
  // or
  // ar := automata.New(automata.EthernetArduino, "192.168.0.13:13666")
  ar.SetDigitalOutput(13)
  for {
    ar.Toggle(13)
    temp := ar.Temp()
    fmt.Printf("Temp: %d C\n", temp)
  }
}
```
