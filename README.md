# Coffer

## Summary

Service responsible for asset/recording management via the Pulse API


### External dependencies

- Go 1.5 or greater should be installed in your $PATH
- GOPATH should be set as described in http://golang.org/doc/code.html
- MongoDB GridFS.  Assets are primarily stored in GridFS
- Glide


### Development Setup


```bash
go get gitlab.vailsys.com/jerny/coffer
cd $GOPATH/src/gitlab.vailsys.com/jerny/coffer
```

to install vendored dependencies run

```bash
make setup
make install
```


### Running Tests

We are using [Ginkgo](https://github.com/onsi/ginkgo), for our test framework.

```bash
make test
```

## Next Steps

Design [Design](docs/design.md)

