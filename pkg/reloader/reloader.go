package reloader

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	promconfig "github.com/prometheus/prometheus/config"
	"k8s.io/klog"
)

// Config for creating a new Reloader
type Config struct {
	ConfigFile    string
	ReloadURL     string
	WatchInterval time.Duration
}

// Reloader orchestrates a process of watching Prometheus config
// and its rule files for changes, calling a Prometheus reload URL
// every time a change is detected.
type Reloader struct {
	configFile    string
	reloadURL     string
	watchInterval time.Duration

	lastConfigHash []byte
	lastRulesHash  []byte
}

// New creates new Reloader using ReloaderConfig
func New(rc *Config) *Reloader {
	return &Reloader{
		configFile:    rc.ConfigFile,
		reloadURL:     rc.ReloadURL,
		watchInterval: rc.WatchInterval,
	}
}

// Watch starts a periodic watch for config and rule file changes
func (r *Reloader) Watch(ctx context.Context) error {
	if err := r.apply(ctx); err != nil {
		return err
	}

	tick := time.NewTicker(r.watchInterval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
		}

		if err := r.apply(ctx); err != nil {
			klog.Error(err)
		}
	}
}

func (r *Reloader) apply(ctx context.Context) error {
	var configHash []byte

	configDir := filepath.Dir(r.configFile)

	configFileContents, err := ioutil.ReadFile(r.configFile)
	if err != nil {
		return fmt.Errorf("failed to read the contents of config file %q: %w", r.configFile, err)
	}

	h := sha256.New()
	h.Write(configFileContents)
	configHash = h.Sum(nil)
	klog.V(4).Infof("configHash: %x", configHash)

	// Load and parse Prometheus config
	config, err := promconfig.Load(string(configFileContents))
	if err != nil {
		return fmt.Errorf("failed to read Prometheus config: %w", err)
	}
	klog.V(5).Infof("config: %#v", config)

	ruleFiles, err := getRuleFiles(config, configDir)
	if err != nil {
		return fmt.Errorf("failed to get rule files from Prometheus config: %w", err)
	}
	klog.V(4).Infof("ruleFiles: %v", ruleFiles)

	h = sha256.New()
	for _, ruleFile := range ruleFiles {
		err := hashFile(h, ruleFile)
		if err != nil {
			return fmt.Errorf("failed to hash rule file %q: %w", ruleFile, err)
		}
	}
	rulesHash := h.Sum(nil)
	klog.V(4).Infof("rulesHash: %x", rulesHash)

	if bytes.Equal(r.lastConfigHash, configHash) && bytes.Equal(r.lastRulesHash, rulesHash) {
		klog.V(2).Info("No config changes detected")
		return nil
	}

	// Do not reload Prometheus on first iteration
	if r.lastConfigHash != nil {
		klog.V(0).Info("Config change detected, reloading Prometheus")
		if err := r.triggerReload(ctx); err != nil {
			return fmt.Errorf("failed to reload Prometheus: %w", err)
		}
	}

	// we only update hashes if triggerReload above succeeds
	r.lastConfigHash = configHash
	r.lastRulesHash = rulesHash

	return nil
}

func (r *Reloader) triggerReload(ctx context.Context) error {
	req, err := http.NewRequest("POST", r.reloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("reload request failed: %w", err)
	}

	defer func() {
		// Drain up to 1 Kb and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 1024)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("received non-200 response: %s; have you set `--web.enable-lifecycle` Prometheus flag?", resp.Status)
	}

	return nil
}

func getRuleFiles(config *promconfig.Config, configDir string) ([]string, error) {
	var ruleFiles []string

	for _, rf := range config.RuleFiles {
		arf := filepath.Join(configDir, rf)
		rfs, err := filepath.Glob(arf)
		if err != nil {
			return nil, err
		}

		ruleFiles = append(ruleFiles, rfs...)
	}

	return ruleFiles, nil
}

func hashFile(h hash.Hash, fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	return nil
}
