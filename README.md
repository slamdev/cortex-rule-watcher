# cortex-rule-watcher <a href="https://hub.docker.com/r/slamdev/cortex-rule-watcher"><img alt="status" src="https://img.shields.io/docker/v/slamdev/cortex-rule-watcher"></a>

Ruler component in [Cortex](https://github.com/cortexproject/cortex)/[Mimir](https://github.com/grafana/mimir) 
does not support rules extraction from [PrometheusRule](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/alerting.md#prometheusrule-labelling) resources.

**cortex-rule-watcher** can help you with that. It can be added as a sidecar container to the ruler deployment and
will watch for the changes in PrometheusRule resources in the cluster. As soon as it sees a change in PrometheusRule resource,
it will dump the PrometheusRule spec to the file, that will be synced by ruler component.

[Click to see a full example](example/example.yaml)

Or use this configuration for [mimir helm chart](https://github.com/grafana/helm-charts/tree/mimir-distributed-2.0.7/charts/mimir-distributed):

```yaml
ruler:
  extraContainers:
    - name: cortex-rule-watcher
      args:
        - --rule-path=/local-rules
      image: slamdev/cortex-rule-watcher:latest
      imagePullPolicy: Always
      ports:
        - name: http-health
          containerPort: 8081
        - name: http-metrics
          containerPort: 8082
      livenessProbe:
        httpGet:
          path: /healthz
          port: 8081
        initialDelaySeconds: 15
        periodSeconds: 20
      readinessProbe:
        httpGet:
          path: /readyz
          port: 8081
        initialDelaySeconds: 5
        periodSeconds: 10
      resources:
        requests:
          cpu: 10m
          memory: 16Mi
      volumeMounts:
        - name: rules
          mountPath: /local-rules
  extraVolumes:
    - name: rules
      emptyDir: { }
  extraVolumeMounts:
    - name: rules
      mountPath: /ruler/fake
  extraArgs:
    ruler-storage.backend: local
    ruler-storage.local.directory: /ruler
```

Besides that you have to have a cluster role and binding:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: mimir-distributed-cortex-rule-watcher
rules:
  - apiGroups:
      - monitoring.coreos.com
    resources:
      - prometheusrules
    verbs:
      - get
      - list
      - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: mimir-distributed-cortex-rule-watcher
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: mimir-distributed-cortex-rule-watcher
subjects:
  - kind: ServiceAccount
    name: mimir-distributed
    namespace: monitoring
```
