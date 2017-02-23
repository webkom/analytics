# Analytics

> Simple metric collector

This service exposes a simple event ingestion endpoint. We use this to track events
and trends. The metrics gets stored in postgres, services like Redash is used to query the
data.


## Configuration

We use environment variables to configure the service.

* SENTRY_DSN - Send errors to Sentry
* LISTEN_ADDRESS="127.0.0.1:8000" - Address and port to listen on
* POSTGRES_URL="postgres://analytics:analytics@localhost/analytics?sslmode=disable" - Postgres connection details

You need to set the LISTEN_ADDRESS and POSTGRES_URL variables to start the service.


## CLI

Start the apiserver
```
./analytics
```

Create the database table
```
./analytics -migrate
```


## Getting started

You need to install go and configure the GOPATH before you can build and run this project.
We use godep to manage dependencies, run the line bellow to install it.

```
go get github.com/tools/godep
```

Pull dependencies and run the program.

```
godep get
go build
./analytics
```