package runtimeconfig

import (
	"fmt"
	"os"
	"sync"
)

type RuntimeConfig struct {
	data map[string]string
	mu   sync.RWMutex
}

const mKeyDefaultValue string = ""

func NewRuntimeConfig(defaultKeys []string) *RuntimeConfig {
	cm := &RuntimeConfig{
		data: make(map[string]string),
	}
	for _, key := range defaultKeys {
		cm.data[key] = mKeyDefaultValue
	}
	return cm
}

func (rconfig *RuntimeConfig) Set(key, value string) {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	rconfig.data[key] = value
}

func (rconfig *RuntimeConfig) Get(key string) string {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	return rconfig.data[key]
}

func (rconfig *RuntimeConfig) Delete(key string) {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	delete(rconfig.data, key)
}

func (rconfig *RuntimeConfig) Keys() []string {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	keys := make([]string, 0, len(rconfig.data))
	for key := range rconfig.data {
		keys = append(keys, key)
	}
	return keys
}

func (rconfig *RuntimeConfig) Size() int {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	return len(rconfig.data)
}

func (rconfig *RuntimeConfig) LoadValueFromEnv() {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	for key := range rconfig.data {
		rconfig.data[key] = os.Getenv(key)
	}
}

func (rconfig *RuntimeConfig) PrintStatus() {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	for key, value := range rconfig.data {
		if value == "" {
			fmt.Printf("%s: (not set)\n", key)
		} else {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}
