package utils

import (
	"fmt"
	"go-pentor-bank/internal/clog"
	"runtime"
	"sync"
)

var waitGroup *sync.WaitGroup

type FunctionWithArgs struct {
	Fn   func(args ...interface{})
	Args []interface{}
}

func SafeGoRoutinesFunc(fwa FunctionWithArgs) {
	defer func() {
		if rec := recover(); rec != nil {
			log := clog.GetLog()
			err, ok := rec.(error)
			if !ok {
				err = fmt.Errorf("%v", rec)
			}
			stack := make([]byte, 4<<10) // 4KB
			length := runtime.Stack(stack, false)
			log.Error().Stack().Err(err).Msg(string(stack[:length]))
		}
	}()
	wg := getWaitGroup()
	wg.Add(1)
	defer wg.Done()
	fwa.Fn(fwa.Args...)
}

func SafeGoRoutines(fn func()) {
	defer func() {
		if rec := recover(); rec != nil {
			log := clog.GetLog()
			err, ok := rec.(error)
			if !ok {
				err = fmt.Errorf("%v", rec)
			}
			stack := make([]byte, 4<<10) // 4KB
			length := runtime.Stack(stack, false)
			log.Error().Stack().Err(err).Msg(string(stack[:length]))
		}
	}()
	wg := getWaitGroup()
	wg.Add(1)
	defer wg.Done()
	fn()
}

func WaitGoRoutines() {
	wg := getWaitGroup()
	wg.Wait()
}

func getWaitGroup() *sync.WaitGroup {
	if waitGroup == nil {
		waitGroup = &sync.WaitGroup{}
	}
	return waitGroup
}
