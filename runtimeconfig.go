package runtimeconfig

import (
	"fmt"
	"os"
	"sync"
)

type RuntimeConfig struct {
	data       map[string]string // where our data is stored
	ignoreKeys map[string]bool   // mainly used for validation step
	mu         sync.RWMutex      // mutex for thread safe
}

const mKeyDefaultValue string = ""

func NewRuntimeConfig(defaultKeys, ignoreKeys []string) *RuntimeConfig {
	cm := &RuntimeConfig{
		data:       make(map[string]string),
		ignoreKeys: make(map[string]bool),
	}
	for _, key := range defaultKeys {
		cm.data[key] = mKeyDefaultValue
	}
	for _, key := range ignoreKeys {
		cm.ignoreKeys[key] = true
	}
	return cm
}

func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
func getKeysBool(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func (rconfig *RuntimeConfig) Copy() *RuntimeConfig {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()

	newData := make(map[string]string, len(rconfig.data))
	for key, value := range rconfig.data {
		newData[key] = value
	}

	newIgnoreKeys := make(map[string]bool, len(rconfig.ignoreKeys))
	for key, value := range rconfig.ignoreKeys {
		newIgnoreKeys[key] = value
	}

	return &RuntimeConfig{
		data:       newData,
		ignoreKeys: newIgnoreKeys,
	}
}

func (rconfig *RuntimeConfig) ClearData() {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	rconfig.data = make(map[string]string)
}

func (rconfig *RuntimeConfig) ClearIgnoreKeys() {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	rconfig.ignoreKeys = make(map[string]bool)
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

func (rconfig *RuntimeConfig) AddIgnoreKeys(keys ...string) {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	for _, key := range keys {
		if rconfig.ignoreKeys[key] {
			fmt.Printf("Key '%s' is already in ignoreKeys.\n", key)
			continue
		}

		rconfig.ignoreKeys[key] = true
		fmt.Printf("Key '%s' added to ignoreKeys.\n", key)
	}
}

func (rconfig *RuntimeConfig) AddIgnoreKey(key string) {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()

	if rconfig.ignoreKeys[key] {
		fmt.Printf("Key '%s' is already in ignoreKeys.\n", key)
		return
	}

	rconfig.ignoreKeys[key] = true
	fmt.Printf("Key '%s' added to ignoreKeys.\n", key)
}

func (rconfig *RuntimeConfig) RemoveIgnoreKey(key string) {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()

	if !rconfig.ignoreKeys[key] {
		fmt.Printf("Key '%s' is not in ignoreKeys.\n", key)
		return
	}

	delete(rconfig.ignoreKeys, key)
	fmt.Printf("Key '%s' removed from ignoreKeys.\n", key)
}

func (rconfig *RuntimeConfig) IgnoreKeys() []string {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	keys := make([]string, 0, len(rconfig.data))
	for key := range rconfig.data {
		keys = append(keys, key)
	}
	return keys
}

func (rconfig *RuntimeConfig) LoadValueFromEnv() {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	for key := range rconfig.data {
		rconfig.data[key] = os.Getenv(key)
	}
}

func (rconfig *RuntimeConfig) ValuesLoaded() bool {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	for key, value := range rconfig.data {
		if rconfig.ignoreKeys[key] {
			continue // skip current item if ignore
		}
		if value == "" {
			return false // if any item empty return false
		}
	}
	return true
}

func (rconfig *RuntimeConfig) PrintMissingValues() {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	for key, value := range rconfig.data {
		if rconfig.ignoreKeys[key] {
			continue
		}
		if value == "" {
			fmt.Printf("%s: (not set)\n", key)
		}
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
