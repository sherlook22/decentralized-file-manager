package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "somekey"
	pathKey := CASPathTransormFunc(key)
	expectedOriginalkey := "f16b09e156438c4620530809fc44c5502a5378da"
	expectedPathName := "f16b0/9e156/438c4/62053/0809f/c44c5/502a5/378da"
	if pathKey.Pathname != expectedPathName {
		t.Errorf("have %s, want %s", pathKey.Pathname, expectedPathName)
	}
	if pathKey.filename != expectedOriginalkey {
		t.Errorf("have %s, want %s", pathKey.filename, expectedOriginalkey)
	}
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransFromFunc: CASPathTransormFunc,
	}
	s := NewStore(opts)
	key := "somekey"

	data := []byte("some data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransFromFunc: CASPathTransormFunc,
	}
	s := NewStore(opts)
	key := "somekey"

	data := []byte("some data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("have %t, want %t", ok, true)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("have %s, want %s", b, data)
	}

	s.Delete(key)
}
