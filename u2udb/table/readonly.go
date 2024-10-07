package table

import "github.com/sesanetwork/go-vassalo/sesadb"

type IteratedReader struct {
	prefix     []byte
	underlying sesadb.IteratedReader
}

func (t *IteratedReader) Has(key []byte) (bool, error) {
	return t.underlying.Has(prefixed(key, t.prefix))
}

func (t *IteratedReader) Get(key []byte) ([]byte, error) {
	return t.underlying.Get(prefixed(key, t.prefix))
}

func (t *IteratedReader) NewIterator(itPrefix []byte, start []byte) sesadb.Iterator {
	return &iterator{t.underlying.NewIterator(prefixed(itPrefix, t.prefix), start), t.prefix}
}

func (t *Table) GetSnapshot() (sesadb.Snapshot, error) {
	snap, err := t.underlying.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return &snapshot{
		IteratedReader: IteratedReader{
			prefix:     t.prefix,
			underlying: snap,
		},
		snap: snap,
	}, nil
}

func (t *Table) Stat(property string) (string, error) {
	return t.underlying.Stat(property)
}

/*
 * Iterator
 */

type iterator struct {
	it     sesadb.Iterator
	prefix []byte
}

func (it *iterator) Next() bool {
	return it.it.Next()
}

func (it *iterator) Error() error {
	return it.it.Error()
}

func (it *iterator) Key() []byte {
	return noPrefix(it.it.Key(), it.prefix)
}

func (it *iterator) Value() []byte {
	return it.it.Value()
}

func (it *iterator) Release() {
	it.it.Release()
	*it = iterator{}
}

type snapshot struct {
	IteratedReader
	snap sesadb.Snapshot
}

func (s *snapshot) Release() {
	s.snap.Release()
}
