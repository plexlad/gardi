package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

type JsonDB struct {
	basePath string
	mu       sync.RWMutex
}

func NewJsonDB(basePath string) *JsonDB {
	return &JsonDB{basePath: basePath}
}

func (db *JsonDB) Set(collection, user, entry string, data any) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dir := filepath.Join(db.basePath, collection, user)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	filePath := filepath.Join(dir, entry+".json")

	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

func (db *JsonDB) Get(collection, user, entry string, dest any) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// checks if dest is a pointer
	if reflect.ValueOf(dest).Kind() != reflect.Pointer {
		return fmt.Errorf("dest must be a pointer")
	}

	filepath := filepath.Join(db.basePath, collection, user, entry+".json")

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("entry not found")
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

func (db *JsonDB) Delete(collection, user, entry string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	filePath := filepath.Join(db.basePath, collection, user, entry+".json")

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("entry not found")
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (db *JsonDB) List(collection, user string) ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dir := filepath.Join(db.basePath, collection, user)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	names := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			// remove .json extension
			name := entry.Name()[:len(entry.Name())-5]
			names = append(names, name)
		}
	}

	return names, nil
}

func (db *JsonDB) ListAll(collection string) (map[string][]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	collectionPath := filepath.Join(db.basePath, collection)

	users, err := os.ReadDir(collectionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]string{}, nil
		}
		return nil, fmt.Errorf("failed to read collection: %w", err)
	}

	result := make(map[string][]string)
	for _, userDir := range users {
		if !userDir.IsDir() {
			continue
		}

		entries, err := db.List(collection, userDir.Name())
		if err != nil {
			continue
		}

		result[userDir.Name()] = entries
	}

	return result, nil
}
