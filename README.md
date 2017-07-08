# Smart Environment

This is a suite of applications for visualizing sensor data in different ways. Users can deploy an IoT application for the RIOT operating system on any number of low-power IoT boards that will measure data of connected sensors whenever it is polled from a central server running on a Raspberry Pi. The Raspberry Pi stores a database that a web server can connect to in order to visualize a history of sensor data on a website. An alternate way of visualizing sensor data is on an LED stripe connected to any of the IoT boards. Drivers are provided for WS2812B running on a Phytec phyNODE KW22 board.

*Phytec phyNODE KW22*

    MCU                    | MKW22D512
    -----------------------|------------------------------------
    Family                 | ARM Cortex-M4
    Vendor                 | Kinetis
    RAM/ROM                | 65Kb / 512Kb
    Frequency              | up to 50MHz
    Timers                 | up to 12 (16-bit, 24-bit, 32-bit)
    ADCs                   | 1x 16-bit (8 channels)
    UARTs / SPIs / I2Cs    | 3 / 1 / 2
    Vcc                    | 1.8V - 3.6V
    Radio                  | IEEE802.15.4 @ 2,4GHZ
    Sensors                | diverse
    
Sensor data sent from a board to a Raspberry Pi is serialized with SenML. The protocol of communication between the boards and the Raspberry Pi is CoAP. Further details of the application suite's architecture and architectural choices made during development can be looked up in the doc folder or on the wiki.

Refer to each module's README.md for details regarding respective prerequisites and installation instructions.

## board-app

This is the C-based RIOT application that is to be deployed on the IoT boards. A Phytec phyNODE KW22 board was used but other boards can be used as long as drivers implementing the SAUL interface are available for the sensors to be polled. All SAUL compatible sensors will be automatically initialized but CoAP resources for the corresponding sensors can and must be created manually, which is a simple process that allows for increased flexibility.

## RiotPi

This is the GO-based application running on the Raspberry Pi that polls registered devices for sensor data regularly and hosts and initializes the database for sensor data.

## RiotDashboard

This is the PHP-based website to be deployed on a web server, e.g. running on XAMPP, on the Raspberry Pi or another device. It presents the history of sensor data from RiotPi's database graphically on a dynamic web frontend.

## server-emulator

If a Raspberry Pi is not available, this C-based RIOT application can be deployed on another IoT board for testing purposes. A desired CoAP resource will be requested from the board-app running on the other IoT board. The response will be printed on the terminal.
