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
  sensors.setResolution(insideThermometer, 9);
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

void writeResponse(byte b[]) {
  Serial.write(b, 4);
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
