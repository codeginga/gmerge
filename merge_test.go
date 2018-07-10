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
	gm := gmerge.New()
	gm.Add(t1)
	gm.Add(t2)

	gm.AddFs(te1, te2)

	errs := gm.Run()

	fnd := 0
	for _, err := range errs {
		if err.Error() == "error1" || err.Error() == "error2" {
			fnd++
			continue
		}

		t.Errorf("expected errors error1 or error2 but got %s", err.Error())
	}
}

func ExampleMerger() {
	gm := gmerge.New()

	gm.Add(func() error {
		time.Sleep(time.Second * 1)
		return nil
	})

	gm.Add(func() error {
		time.Sleep(time.Second * 2)
		return errors.New("error1")
	})

	errs := gm.Run()
	fmt.Println(errs)
	// Output: [error1]
}
