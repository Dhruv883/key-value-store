package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	store := InitKVStore[string, string]()
	key, value := "testKey", "testValue"

	err := store.Put(key, value)
	assert.NoError(t, err)

	err = store.Put(key, value)
	assert.Error(t, err)

	val, err := store.Get(key)
	assert.Equal(t, value, val)
}

func TestGet(t *testing.T) {
	store := InitKVStore[string, string]()
	key, invalidKey, value := "testKey", "invalidTestKey", "testValue"

	err := store.Put(key, value)
	assert.NoError(t, err)

	val, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, val)

	val, err = store.Get(invalidKey)
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	store := InitKVStore[string, string]()
	key, initialValue, newValue := "testKey", "testValue", "newValue"

	err := store.Put(key, initialValue)
	assert.NoError(t, err)

	err = store.Update(key, newValue)
	assert.NoError(t, err)

	val, err := store.Get(key)
	assert.Equal(t, newValue, val)
}

func TestDelete(t *testing.T) {
	store := InitKVStore[string, string]()
	key, initialValue := "testKey", "testValue"

	err := store.Put(key, initialValue)
	assert.NoError(t, err)

	err = store.Delete(key)
	assert.NoError(t, err)

	err = store.Delete(key)
	assert.Error(t, err)
}

func TestPutWithTTL(t *testing.T) {
	store := InitKVStore[string, string]()
	key, value := "ttlKey", "ttlValue"

	err := store.PutWithTTL(key, value, 100*time.Millisecond)
	assert.NoError(t, err)

	val, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, val)

	// Wait for the key to expire.
	time.Sleep(150 * time.Millisecond)

	_, err = store.Get(key)
	assert.Error(t, err)
}

func TestSetTTL(t *testing.T) {
	store := InitKVStore[string, string]()
	key, value := "testKey", "testValue"

	err := store.Put(key, value)
	assert.NoError(t, err)

	err = store.SetTTL(key, 100*time.Millisecond)
	assert.NoError(t, err)

	val, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, val)

	time.Sleep(150 * time.Millisecond)

	_, err = store.Get(key)
	assert.Error(t, err)
}

func TestSetTTL_PreservesValue(t *testing.T) {
	store := InitKVStore[string, string]()
	key, value := "testKey", "testValue"

	err := store.Put(key, value)
	assert.NoError(t, err)

	err = store.SetTTL(key, time.Minute)
	assert.NoError(t, err)

	val, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, val)
}

func TestUpdate_PreservesTTL(t *testing.T) {
	store := InitKVStore[string, string]()
	key := "ttlKey"

	err := store.PutWithTTL(key, "original", 150*time.Millisecond)
	assert.NoError(t, err)

	err = store.Update(key, "updated")
	assert.NoError(t, err)

	// Value should reflect the update.
	val, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, "updated", val)

	// TTL should still be in effect — key must expire after the original duration.
	time.Sleep(200 * time.Millisecond)

	_, err = store.Get(key)
	assert.Error(t, err)
}
