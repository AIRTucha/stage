package main

import (
	"time"
)

type Debouncer struct {
	timer    *time.Timer
	debounce time.Duration
}

func NewDebounce(debounce int) *Debouncer {
	return &Debouncer{
		timer:    nil,
		debounce: time.Millisecond * time.Duration(debounce),
	}
}

func (debouncer *Debouncer) Run(action func()) {
	if debouncer.timer != nil {
		debouncer.timer.Stop()
	}
	debouncer.timer = time.AfterFunc(debouncer.debounce, action)
}
