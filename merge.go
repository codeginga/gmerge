package gmerge

import (
	"context"
	"sync"
)

type merge struct {
	gfuncs []GFunc
	ech    chan error
	wg     sync.WaitGroup
}

func (m *merge) Add(gf GFunc) {
	m.gfuncs = append(m.gfuncs, gf)
}

func (m *merge) AddFs(gfs ...GFunc) {
	for _, gf := range gfs {
		m.Add(gf)
	}
}

func (m *merge) runI(i int) {
	defer m.wg.Done()

	if err := m.gfuncs[i](); err != nil {
		m.ech <- err
	}
}

func (m *merge) Run() []error {

	for i := range m.gfuncs {
		m.wg.Add(1)
		go m.runI(i)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer close(m.ech)

		m.wg.Wait()
		cancel()
	}()

	errs := []error{}
	done := false

	for !done {
		select {
		case err, ok := <-m.ech:
			if !ok {
				m.ech = nil
				continue
			}

			errs = append(errs, err)

		case <-ctx.Done():
			done = true
		}
	}

	return errs
}

// New returns instance of Merger
func New() Merger {
	return &merge{
		gfuncs: []GFunc{},
		ech:    make(chan error),
		wg:     sync.WaitGroup{},
	}
}
