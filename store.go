package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "sherlook22#key"

func CASPathTransormFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key)) // [20]byte => []byte => [:]
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	path := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		path[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(path, "/"),
		filename: hashStr,
	}
}

type PathTransFromFunc func(string) PathKey

type PathKey struct {
	Pathname string
	filename string
}

func (p PathKey) FirstPathName() string {
	return strings.Split(p.Pathname, "/")[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.filename)
}

var DefaultTransFromFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		filename: key,
	}
}

type StoreOpts struct {
	// Root is the folder name of the root, containing all the folders/files of system.
	Root              string
	PathTransFromFunc PathTransFromFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransFromFunc == nil {
		opts.PathTransFromFunc = DefaultTransFromFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransFromFunc(key)
	_, err := os.Stat(s.Root + "/" + pathKey.FullPath())
	return err == nil
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransFromFunc(key)
	return os.RemoveAll(s.Root + "/" + pathKey.FirstPathName())
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransFromFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(fullPathWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransFromFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("wrote %d bytes to %s\n", n, fullPathWithRoot)

	return nil
}
