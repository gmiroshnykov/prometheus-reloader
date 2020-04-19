# prometheus-reloader

A sidecar that reloads Prometheus when config file or rule files change.


## Usage

```
$ prometheus-reloader -config-file testdata/prometheus.yml
I0419 15:30:46.124751   72508 main.go:33] Starting prometheus-reloader 7f0d912 (7f0d912bba736278b2752943bf6287b4658b3490)
```

After changing `prometheus.yml` or any of the rule files:

```
I0419 15:31:06.127397   72508 reloader.go:115] Config change detected, reloading Prometheus at http://127.0.0.1:9090/-/reload
```

## Flags

```
  -config-file string
    	Prometheus configuration file path (default "/etc/prometheus/prometheus.yml")
  -reload-url string
    	Prometheus reload endpoint (default "http://127.0.0.1:9090/-/reload")
  -watch-interval duration
    	Interval for watching config and rules files for changes (default 10s)
```

Run `prometheus-reloader -h` to see the rest of the flags.


## Testing with skaffold

1. Run:
```
skaffold dev --port-forward
```

2. Open [http://127.0.0.1:9090/config](http://127.0.0.1:9090/config) in your browser.

3. Make some changes to `k8s/configmap.yaml`.
  Skaffold will apply the changes automatically once you save the file.

4. Wait for [up to two minutes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#mounted-configmaps-are-updated-automatically) for Kubernetes to propagate ConfigMap changes.

5. Observe something like:
```
[prometheus reloader] I0419 14:27:22.139774       1 reloader.go:115] Config change detected, reloading Prometheus at http://127.0.0.1:9090/-/reload
[prometheus] level=info ts=2020-04-19T14:27:22.143Z caller=main.go:788 msg="Loading configuration file" filename=/etc/prometheus/prometheus.yml
[prometheus] level=info ts=2020-04-19T14:27:22.147Z caller=main.go:816 msg="Completed loading of configuration file" filename=/etc/prometheus/prometheus.yml
```

6. Refresh Prometheus config page in your browser to make sure Prometheus has been reloaded.
