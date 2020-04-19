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
