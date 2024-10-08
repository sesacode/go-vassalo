package memorydb

import (
	"math/rand"
	"sync"

	"github.com/sesanetwork/go-vassalo/hash"
	"github.com/sesanetwork/go-vassalo/sesadb"
)

type fakeFS struct {
	Namespace string
	Files     map[string]sesadb.Store

	sync.RWMutex
}

var (
	fakeFSs = make(map[string]*fakeFS)
	fakeFSl = new(sync.Mutex)
)

func newFakeFS(namespace string) *fakeFS {
	if namespace == "" {
		namespace = uniqNamespace()
	}

	fakeFSl.Lock()
	defer fakeFSl.Unlock()

	if fs, ok := fakeFSs[namespace]; ok {
		return fs
	}

	fs := &fakeFS{
		Namespace: namespace,
		Files:     make(map[string]sesadb.Store),
	}
	fakeFSs[namespace] = fs
	return fs
}

func uniqNamespace() string {
	return hash.FakeHash(rand.Int63()).Hex() // nolint:gosec
}

func (fs *fakeFS) ListFakeDBs() []string {
	fs.RLock()
	defer fs.RUnlock()

	ls := make([]string, 0, len(fs.Files))
	for f := range fs.Files {
		ls = append(ls, f)
	}

	return ls
}

func (fs *fakeFS) OpenFakeDB(name string) sesadb.Store {
	fs.Lock()
	defer fs.Unlock()

	drop := func() {
		delete(fs.Files, name)
	}

	db := NewWithDrop(drop)

	if oldDB, ok := fs.Files[name]; ok {
		return oldDB
	}
	fs.Files[name] = db

	return db
}
