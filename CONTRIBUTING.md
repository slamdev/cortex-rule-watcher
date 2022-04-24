# Contributing

To regenerate mocks:

```shell
mockgen -source=internal/syncer.go -package=internal -destination=internal/syncermock_test.go
```

To make a release:

```shell
TAG=x.x.x && git tag -a ${TAG} -m "make ${TAG} release" && git push --tags
```
