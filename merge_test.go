package gmerge_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/codeginga/gmerge"
)

func t1() error {
	time.Sleep(time.Second * 2)
	return nil
}

func t2() error {
	time.Sleep(time.Second * 3)
	return nil
}

func te1() error {
	time.Sleep(time.Second * 2)
	return errors.New("error1")
}

func te2() error {
	time.Sleep(time.Second * 2)
	return errors.New("error2")
}

func TestMerge(t *testing.T) {
	tt := []struct {
		name string
		f    gmerge.GFunc
		tag  string
		err  error
	}{
		{"success test 1", t1, "test 1", nil},
		{"success test 2", t2, "test 2", nil},
		{"fail test 3", te1, "test 3", errors.New("error1")},
		{"fail test 4", te2, "test 4", errors.New("error2")},
	}

	gm := gmerge.New()
	for _, t := range tt {
		gm.Add(t.tag, t.f)
	}

	merr := gm.Run()

	for _, s := range tt {
		t.Run(s.name, func(t *testing.T) {
			exp := fmt.Sprintf("%v", s.err)
			recv := fmt.Sprintf("%v", merr[s.tag])
			if exp != recv {
				t.Errorf("expected err is %s but got %s", exp, recv)
			}
		})
	}

}

func ExampleMerger() {
	gm := gmerge.New()

	gm.Add("test1", func() error {
		time.Sleep(time.Second * 1)
		return nil
	})

	gm.Add("test2", func() error {
		time.Sleep(time.Second * 2)
		return errors.New("error1")
	})

	merr := gm.Run()
	fmt.Println(merr["test2"])
	// Output: error1
}
