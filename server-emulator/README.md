# Prerequisites

## RIOT

Clone the RIOT OS repository into your local drive.

```sh
git clone https://github.com/RIOT-OS/RIOT
```

This release has been tested with RIOT release 2017.04 on an Atmel SAM R21 Xplained Pro board. Other branches and boards may be used as well but bugs may occur. In that case, you can switch to the branch corresponding to the release that has been tested:

```sh
git fetch
git checkout 2017.04-branch
```

# Build, flash and connect to terminal

By default, server-emulator will build with the following parameters:

- BOARD ?= samr21-xpro
- RIOTBASE ?= $(CURDIR)/../RIOT

This means that the application will be built for an Atmel SAM R21 Xplained Pro board and the path of the RIOT root folder will be assumed to be two directories above the server-emulator directory (e.g. if you copy or clone the entire team2 folder into the RIOT root folder). In this case, the project can be built simply by running

```sh
make all
```

Flash the application by running

```sh
make flash
```

Connect to the board's terminal via serial connection by running

```sh
make term
```

These commands can be combined as

```sh
make all flash term
```

If the default parameters for BOARD and RIOTBASE do not apply to your configuration, run the commands from above and append the parameters that correspond to your configuration, e.g.

```sh
make all flash term RIOTBASE=YOURPATH BOARD=YOURBOARD
```

