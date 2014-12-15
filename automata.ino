#include <OneWire.h>
#include <DallasTemperature.h>

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
  Serial.begin(57600);
  sensors.begin();
  sensors.getAddress(insideThermometer, 0);
}

void loop() {
  int command;
  int parameter;
  byte buf[2];
  if (Serial.available()) {
    Serial.readBytes(buf, 2);
    command = buf[0];
    parameter = buf[1];
    switch (command) {
      case SET_DIGITAL_OUT:
        pinMode(parameter, OUTPUT);
        Serial.write(0);
        break;
      case DIGITAL_WRITE_HIGH:
        digitalWrite(parameter, HIGH);
        Serial.write(0);
        break;
      case DIGITAL_WRITE_LOW:
        digitalWrite(parameter, LOW);
        Serial.write(0);
        break;
      case READ_TEMP:
        byte temp;
        sensors.requestTemperatures();
        temp = byte(sensors.getTempC(insideThermometer));
        Serial.write(temp);
        break;
      case PING:
        Serial.write(0);
        break;
      default:
        Serial.write(1);
    }
  }
}
