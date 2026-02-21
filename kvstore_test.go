package main

import (
	"testing"

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
