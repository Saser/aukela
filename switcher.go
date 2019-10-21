package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/icccm"
)

type Switcher struct {
	logger *log.Logger
	xu     *xgbutil.XUtil

	focused <-chan xproto.Window

	done chan struct{}
}

func NewSwitcher(
	logger *log.Logger,
	xu *xgbutil.XUtil,
	focused <-chan xproto.Window,
) *Switcher {
	return &Switcher{
		logger: logger,
		xu:     xu,

		focused: focused,

		done: make(chan struct{}, 1),
	}
}

func (s *Switcher) Start() error {
	s.logger.Println("starting switcher")
loop:
	for {
		select {
		case window := <-s.focused:
			s.logger.Printf("getting class of window %d", window)
			wmClass, err := icccm.WmClassGet(s.xu, window)
			if err != nil {
				s.logger.Printf("getting class of window %d failed: %+v", window, err)
				continue
			}
			s.logger.Printf("got class of window %d: %+v", window, wmClass)
		case <-s.done:
			s.logger.Println("received done signal")
			break loop
		}
	}
	s.logger.Println("stopped switcher")
	return nil
}

func (s *Switcher) Stop() {
	s.logger.Println("sending done signal")
	s.done <- struct{}{}
	s.logger.Println("sent done signal")
}
