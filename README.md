# Coffer

## Summary

Service responsible for asset/recording management via the Pulse API


### External dependencies

- Go 1.6 or greater should be installed in your $PATH
- GOPATH should be set as described in http://golang.org/doc/code.html
- MongoDB GridFS.  Assets are primarily stored in GridFS
- Govendor
- Ginkgo for testing

### Development Setup


```bash
go get github.com/wallywest/coffer
cd $GOPATH/src/github.com/wallywest/coffer
```

to install vendored dependencies run

```bash
make tools 
```


### Running Tests

We are using [Ginkgo](https://github.com/onsi/ginkgo), for our test framework.

```bash
make test
```


### Building

```bash
make build
```


### Command Line Flags

- `--port <port>`: Port from which to serve. Default is `6000`.
- `--advertise-address <address>`: External address of the server for it to pass to users for various confirmation codes. Default is `127.0.0.1:6000`.
- `--mongo-servers <sddress>`: Address of MongoDB instance to use.
- `--mongo-db`: Mongo database.
- `--mongo-prefix`: Mongo prefix for the GridFS collection.
- `--log-level`: Log level. Debug logging is off by default.
- `--skip-registration`: Flag to disable registration of service.
- `--registry-type`: Registry type for service discovery options are inmem,consul.
- `--registry-nodes`: Nodes for the registry to connect to *only applicable for consul*.

## Next Steps

[Reference](docs/reference.md)

