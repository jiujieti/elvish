package edit

import (
	"sync"

	"src.elv.sh/pkg/cli/histutil"
	"src.elv.sh/pkg/store"
)

// A wrapper of histutil.Store that is concurrency-safe and supports an
// additional FastForward method.
type histStore struct {
	m  sync.Mutex
	db store.Store
	hs histutil.Store
}

func newHistStore(db store.Store) (*histStore, error) {
	hs, err := histutil.NewHybridStore(db)
	return &histStore{db: db, hs: hs}, err
}

func (s *histStore) AddCmd(cmd store.Cmd) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.hs.AddCmd(cmd)
}

func (s *histStore) AllCmds() ([]store.Cmd, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.hs.AllCmds()
}

func (s *histStore) Cursor(prefix string) histutil.Cursor {
	s.m.Lock()
	defer s.m.Unlock()
	return cursor{&s.m, histutil.NewDedupCursor(s.hs.Cursor(prefix))}
}

func (s *histStore) FastForward() error {
	s.m.Lock()
	defer s.m.Unlock()
	hs, err := histutil.NewHybridStore(s.db)
	s.hs = hs
	return err
}

type cursor struct {
	m *sync.Mutex
	c histutil.Cursor
}

func (c cursor) Prev() {
	c.m.Lock()
	defer c.m.Unlock()
	c.c.Prev()
}

func (c cursor) Next() {
	c.m.Lock()
	defer c.m.Unlock()
	c.c.Next()
}

func (c cursor) Get() (store.Cmd, error) {
	c.m.Lock()
	defer c.m.Unlock()
	return c.c.Get()
}
