package gmerge

import (
	"context"
	"sync"
)

type fholder struct {
	gfunc GFunc
	tag   string
}

type eholder struct {
	err error
	tag string
}

type merge struct {
	fholders []fholder
	errch    chan eholder
	wg       sync.WaitGroup
}

func (m *merge) Add(tag string, gf GFunc) {
	m.fholders = append(m.fholders, fholder{
		gfunc: gf,
		tag:   tag,
	})
}

func (m *merge) runI(i int) {
	defer m.wg.Done()

	if err := m.fholders[i].gfunc(); err != nil {
		m.errch <- eholder{
			tag: m.fholders[i].tag,
			err: err,
		}
	}
}

func (m *merge) Run() map[string]error {

	for i := range m.fholders {
		m.wg.Add(1)
		go m.runI(i)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer close(m.errch)

		m.wg.Wait()
		cancel()
	}()

	merr := make(map[string]error)
	done := false

	for !done {
		select {
		case errh, ok := <-m.errch:
			if !ok {
				m.errch = nil
				continue
			}

			merr[errh.tag] = errh.err

		case <-ctx.Done():
			done = true
		}
	}

	return merr
}

// New returns instance of Merger
func New() Merger {
	return &merge{
		fholders: []fholder{},
		errch:    make(chan eholder),
		wg:       sync.WaitGroup{},
	}
}
