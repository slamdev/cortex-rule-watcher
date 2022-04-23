# cortex-rule-watcher

Cortex ruler sidecar to sync PrometheusRule resources

To regenerate mocks:

```shell
mockgen -source=internal/syncer.go -package=internal -destination=internal/syncermock_test.go
```
