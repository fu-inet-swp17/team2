# Prerequisites

Prepare your Raspberry Pi like described [here](https://github.com/RIOT-Makers/wpan-raspbian/wiki/Spice-up-Raspbian-for-the-IoT).
Now follow these steps:

## Go runtime

```sh
sudo apt-get install golang
```

## Git

```sh
sudo apt-get install git
```

## Go Path

In case, there is an error when installing and in the output is something about the GOPATH missing, execute this line:

```sh
export GOPATH=$HOME/go
```

If you also wish to run this program like any other program on the computer, run:

```sh
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```



# Install and first run

```sh
go get github.com/fu-inet-swp17/team2/SenPi
go install github.com/fu-inet-swp17/team2/SenPi
```

Now run this command to create a sample config file:

```sh
SenPi --sampleconfig
```

When this fails (e.g. because the GOPATH is not in your PATH), you can run
```sh
$GOPATH/bin/SenPi --sampleconfig
```

Now edit the newly created config file`$GOPATH/bin/conf.json` with a text editor of your choice.

You can also move the config file to a folder of you choice, you only have to specify the config file by adding this attribute to the command everytime you start the server:
```sh
--config <PATH TO CONFIGFILE>
```

Now you have to initialize the SQL database with the provided credentials in the configuration file by using the following command:

```sh
SenPi --initdb
```

__The server is now up and running.__

If you need to clear the database, run
```sh
SenPi --cleardb
```

If you need to clear and initialisethe database, run
```sh
SenPi --reinitdb
```
This basically combines `--cleardb` and `--initdb`.

# Run

```sh
SenPi
```

When this fails (e.g. because the GOPATH is not in your PATH), you can run
```sh
$GOPATH/bin/SenPi
```

# Tests

To execute the existing unit tests, execute the following line (including the trailing ...):

```sh
go test github.com/fu-inet-swp17/team2/SenPi/...
```

## Used Go Libraries

### Go-logging

`github.com/op/go-logging`

### Go mysql driver

`github.com/go-sql-driver/mysql`

### Go CoAP

`github.com/dustin/go-coap`

### Go SenML

`github.com/nkristek/go-senml`

### Autostart
Es gibt im großen und ganzen zwei Möglichkeiten für einen Autostart. Die
erste ist das einfügen eines bashscripts in den Ordner '/etc/init.d'. Siehe
als Beispiel das SenPi Skript im SenPi Ordner oder als Vorlage das skeleton im
init.d Ordner. Das Skript muss den '### BEGIN INIT INFO' Kommentar am Anfang
enthalten. Nach dem Erstellen des Skripts (das Skript muss für alle User
executable sein 'chmod 755') muss der Befehl:
	sudo update-rc.d meinprog defaults
ausgeführt werden. Beim nächsten Systemstart wird das Skript mit ausgeführt.

Die andere Möglichkeit ist in das Skript /etc/rc.local einen entsprechenden
Befehl einzufügen. Das kann als etwas uneleganter bezeichnet werden, da dieses
Skript dafür eigentlich nicht gedacht ist. Wenn aber ein Befehl zum starten
eines Programms wie z.B. '/usr/bin/tool' zwischen dem Kopfkommentar und den
'exit 0' Befehl eingfügt wird, wird dieses Programm auch bei jedem Systemstart
gestartet.

