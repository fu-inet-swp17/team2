# Smart Environment

This is a suite of applications for visualizing and persisting sensor data in different ways. Users can deploy an IoT application for the RIOT operating system on any number of low-power IoT boards that will measure data of connected sensors whenever it is polled from a central server. <br>
In this example the server is running on a [Raspberry Pi](https://www.raspberrypi.org).
The server stores a database that a web server can connect to in order to visualize a history of sensor data on a web application. <br>
An alternate way of visualizing sensor data is on an LED stripe connected to any of the IoT boards. Drivers for WS2812B are provided for using on a _Phytec phyNODE KW22_ board.

The sensor data sent from a board to the server is serialized following the draft for [SenML](https://tools.ietf.org/html/draft-ietf-core-senml-08). <br>
The protocol of communication between the boards and the Raspberry Pi is [CoAP](https://tools.ietf.org/html/rfc7252). <br>
CoAP is transmitted using [6LoWPAN](https://tools.ietf.org/html/rfc4944). <br>
Further details of the application suite's architecture and architectural choices made during development can be looked up in the [doc](https://github.com/fu-inet-swp17/team2/tree/master/doc) folder or on the wiki.

Refer to each module's `README.md` for details regarding respective prerequisites and installation instructions.

## board-app

This is the C-based RIOT application that is to be deployed on the IoT boards. A _Phytec phyNODE KW22_ board was used but other boards can be used as long as drivers implementing the SAUL interface are available for the sensors to be polled. All SAUL compatible sensors will be automatically initialized but CoAP resources for the corresponding sensors can and must be created manually, which is a simple process that allows for increased flexibility.

### Phytec phyNODE KW22

| MCU                    | MKW22D512 |
|------------------------|------------------------------------|
| Family                 | ARM Cortex-M4 |
| Vendor                 | Kinetis |
| RAM/ROM                | 65Kb / 512Kb |
| Frequency              | up to 50MHz |
| Timers                 | up to 12 (16-bit, 24-bit, 32-bit) |
| ADCs                   | 1x 16-bit (8 channels) |
| UARTs / SPIs / I2Cs    | 3 / 1 / 2 |
| Vcc                    | 1.8V - 3.6V |
| Radio                  | IEEE802.15.4 @ 2,4GHZ |
| Sensors                | diverse |

## SenPi

This is a Golang-based application running on the Raspberry Pi that regularly polls registered devices for sensor data.
It then saves those datapoints in a MySQL database specified in a configuration file preferrably also running on the Raspberry Pi. <br>
Please consider checking the provided `README.md` for more details.

## RiotDashboard

This is the PHP-based website to be deployed on a web server, e.g. running on XAMPP, on the Raspberry Pi or another device. It presents the history of sensor data from RiotPi's database graphically on a dynamic web frontend.

## server-emulator

If a Raspberry Pi is not available, this C-based RIOT application can be deployed on another IoT board for testing purposes. A desired CoAP resource will be requested from the board-app running on the other IoT board. The response will be printed on the terminal.
