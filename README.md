# FTP Watcher for polling stream data

Polls event data from FTP Server and writes them to firebase service.

Keeps track of last modification timestamp to avoid unnecessary write operations.

Build
> go build

Run
> go run .