// Package runtimeconfig provides a map derived struct that can load its
// keys from environment variables
package runtimeconfig

import (
	"fmt"
	"os"
	"sync"
)

// RuntimeConfig a struct for managing environment variables
type RuntimeConfig struct {
	data       map[string]string // where our data is stored
	ignoreKeys map[string]bool   // mainly used for validation step
	mu         sync.RWMutex      // mutex for thread safe
}

// mKeyDefaultValue package const for empty string
const mKeyDefaultValue string = ""

// NewRuntimeConfig returns a RuntimeConfig initialized with defaultKeys
// and ignoreKeys
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

// CreateCopy returns a Copy of RuntimeConfig
func (rconfig *RuntimeConfig) CreateCopy() *RuntimeConfig {
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

// ClearData empties the data from a RuntimeConfig
func (rconfig *RuntimeConfig) ClearData() {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	rconfig.data = make(map[string]string)
}

// ClearIgnoreKeys empties the ignoreKeys map from a RuntimeConfig
func (rconfig *RuntimeConfig) ClearIgnoreKeys() {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	rconfig.ignoreKeys = make(map[string]bool)
}

// Set assigns a key value pair in the RuntimeConfig data prop
func (rconfig *RuntimeConfig) Set(key, value string) {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	rconfig.data[key] = value
}

// Get returns the value provided a key from RuntimeConfig data prop
func (rconfig *RuntimeConfig) Get(key string) string {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	return rconfig.data[key]
}

// Delete removes key value pair from RuntimeConfig data prop
func (rconfig *RuntimeConfig) Delete(key string) {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	delete(rconfig.data, key)
}

// Keys returns the keys from the RuntimeConfig data prop
func (rconfig *RuntimeConfig) Keys() []string {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	keys := make([]string, 0, len(rconfig.data))
	for key := range rconfig.data {
		keys = append(keys, key)
	}
	return keys
}

// Size the size of RuntimeConfig data prop
func (rconfig *RuntimeConfig) Size() int {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	return len(rconfig.data)
}

// AddIgnoreKeys appends multiple keys to the RuntimeConfig ignoreKeys map
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

// AddIgnoreKey appends a single key to the RuntimeConfig ignoreKeys map
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

// RemoveIgnoreKey removes a key from ignore keys in the RuntimeConfig
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

// IgnoreKeys returns a list of ignoreKeys RuntimeConfig
func (rconfig *RuntimeConfig) IgnoreKeys() []string {
	rconfig.mu.RLock()
	defer rconfig.mu.RUnlock()
	keys := make([]string, 0, len(rconfig.data))
	for key := range rconfig.data {
		keys = append(keys, key)
	}
	return keys
}

// LoadValueFromEnv iterates over each key in the data prop
// and calls an os.Getenv to get the value
func (rconfig *RuntimeConfig) LoadValueFromEnv() {
	rconfig.mu.Lock()
	defer rconfig.mu.Unlock()
	for key := range rconfig.data {
		rconfig.data[key] = os.Getenv(key)
	}
}

// ValuesLoaded returns a bool based on all values being populated
// note: items in the ignoreKeys will not count against the overall
// loaded status
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

// PrintMissingValues prints a lists of what values are missing (unset)
// note: items in the ignoreKeys will not count against missing
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

// PrintStatus prints a lists of what values are missing (unset)
// note: this does not take into account ignore list
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
