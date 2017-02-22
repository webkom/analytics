# Analytics

> Simple metric collector

This service exposes a simple event ingestion endpoint. We use this to track events
and trends. The metrics gets stored in postgres, services like Redash is used to query the
data.


## Getting started

You need to install go and configure the GOPATH before you can build and run this project.
We use godep to manage dependencies, run the line bellow to install it.

```
go get github.com/tools/godep
```

Pull dependencies and run the program.

```
godep get
go run main.go
```