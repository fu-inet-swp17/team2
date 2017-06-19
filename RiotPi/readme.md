# Prerequisites

## Go runtime

```sh
sudo apt-get install golang
```

## Git

```sh
sudo apt-get install git
```

## Go Libraries

### Go-logging

```sh
go get github.com/op/go-logging
```

### go mysql driver

```sh
go get github.com/go-sql-driver/mysql
```

### go CoAP

```sh
go get github.com/dustin/go-coap
```

### go SenML

```sh
go get github.com/nkristek/go-senml
```

## Go Path

In case, there is an error when running/building the `main.go` file and in the output is something about the GOPATH missing, execute this line:

```sh
export GOPATH=~/go
```

# First run

Run this command in `.../RiotPi/src/` to create a sample config file:

```sh
go run main.go --sampleconfig
```

Now edit the newly created `conf.json` in a text editor of your choice.

Now you have to initialize the SQL database with the provided credentials in the configuration file by using the following command:

```sh
go run main.go --initdb
```

The server is now up and running.

# Run

Run this command in `.../RiotPi/src/`:

```sh
go run main.go
```


# Build

Run this command in `.../RiotPi/src/`:

```sh
go build main.go
```