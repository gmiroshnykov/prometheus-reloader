package main

import (
	"context"
	"flag"
	"time"

	"k8s.io/klog"

	"github.com/laggyluke/prometheus-reloader/pkg/reloader"
)

var (
	// Version is the build version
	Version = "<development build>"
	// GitHash is the git hash of the build
	GitHash = "N/A"
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	rc := reloader.Config{}

	flag.StringVar(&rc.ConfigFile, "config-file", "/etc/prometheus/prometheus.yml", "Prometheus configuration file path")
	flag.StringVar(&rc.ReloadURL, "reload-url", "http://127.0.0.1:9090/-/reload", "Prometheus reload endpoint")
	flag.DurationVar(&rc.WatchInterval, "watch-interval", 10*time.Second, "Interval for watching config and rules files for changes")

	flag.Parse()

	klog.Infof("Starting prometheus-reloader %s (%s)", Version, GitHash)

	r := reloader.New(&rc)
	err := r.Watch(context.Background())
	if err != nil {
		klog.Fatal(err)
	}
}
