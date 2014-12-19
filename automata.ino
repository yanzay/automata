#include <OneWire.h>
#include <DallasTemperature.h>
#include <SPI.h>
#include <Ethernet.h>

#define ETHERNET true

#ifdef ETHERNET
EthernetServer server = EthernetServer(13666);
byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
byte ip[] = { 169, 254, 68, 110 };
#endif

#define PING 0
#define SET_DIGITAL_OUT 1
#define DIGITAL_WRITE_HIGH 2
#define DIGITAL_WRITE_LOW 3
#define READ_TEMP 4

#define ONE_WIRE_BUS 2

OneWire oneWire(ONE_WIRE_BUS);

DallasTemperature sensors(&oneWire);

DeviceAddress insideThermometer;

void setup() {
#ifdef ETHERNET
  Ethernet.begin(mac, ip);
  server.begin();
#else
  Serial.begin(57600);
#endif
  sensors.begin();
  sensors.getAddress(insideThermometer, 0);
  sensors.setResolution(insideThermometer, 12);
}

void loop() {
#ifdef ETHERNET
  EthernetClient client = server.available();
#endif

  int command;
  int parameter;
  byte buf[2];
  bool avail;
  avail = readRequest(buf);
  if (avail) {
    command = buf[0];
    parameter = buf[1];
    switch (command) {
      case SET_DIGITAL_OUT:
        pinMode(parameter, OUTPUT);
        writeByte(1);
        break;
      case DIGITAL_WRITE_HIGH:
        digitalWrite(parameter, HIGH);
        writeByte(2);
        break;
      case DIGITAL_WRITE_LOW:
        digitalWrite(parameter, LOW);
        writeByte(3);
        break;
      case READ_TEMP:
        float temp;
        sensors.requestTemperatures();
        temp = sensors.getTempC(insideThermometer);
        writeFloat(temp);
        break;
      case PING:
        writeByte(0);
        break;
      default:
        writeByte(1);
        break;
    }
  }
}

#ifdef ETHERNET
bool readRequest(byte *buf) {
  EthernetClient client = server.available();
  if (client.available()) {
    buf[0] = client.read();
    buf[1] = client.read();
    return true;
  } else {
    return false;
  }
}
#else
bool readRequest(byte *buf) {
  if (Serial.available()) {
    Serial.readBytes(buf, 2);
    return true;
  } else {
    return false;
  }
}
#endif

void writeResponse(byte b[]) {
#ifdef ETHERNET
  server.write(b, 4);
#else
  Serial.write(b, 4);
#endif
}

void writeByte(byte resp) {
  byte array[4];
  array[0] = 0;
  array[1] = 0;
  array[2] = 0;
  array[3] = resp;
  writeResponse(array);
}

void writeFloat(float resp) {
  union u_tag {
    byte b[4];
    float fval;
  } u;

  u.fval = resp;
  writeResponse(u.b);
}
