# Linux Job Worker

A Golang client and server to manage Linux jobs (processes).

## Building

Clone the repo:

```bash
$ git clone https://github.com/christarazi/gravitational-challenge /path/to/repo
```

Building the client:

```bash
pushd /path/to/repo/client
  go build -o workerctl
popd
```

Building the server:

```bash
pushd /path/to/repo/server
  go build -o workerd
popd
```

## Usage

Start the server either in one terminal or you can background it with `&`:

```bash
$ /path/to/workerd
```

The client will now be able to send requests:

```bash
$ /path/to/workerctl start -- ls -al
1

$ /path/to/workerctl status 1
Running

$ /path/to/workerctl stop 1

$ /path/to/workerctl status 1
Stopped (ec: 42)
```
