## Modbus reader for AEG PROTECT RCS SPRE-TPRE 

#### aegdcps_modbus_rtu

Read modbus regestry from 1 to 32, and 99 to 111 address, and returns a json object.

Programm flags:

-serial - connected serial port (defaut value "/dev/ttyRS485-1")

-speed - connection speed (defaut value 9600)

-id - modbus slave ID (defaut value 1)

**Note:** serial port static config "8E1" 




Example:

`aegdcps_modbus_rtu -serial=/dev/ttyRS485-1 -speed0=19200 -id=1 `


