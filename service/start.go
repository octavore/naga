package service

import (
	"os"
	"os/signal"
	"syscall"
)

// Start is a convenience method equivalent to `service.Load(m).Run()` and starting the
// app with `./<myapp> start`. Prefer using `Run()` as it is more flexible.
func (s *Service[App]) Start() error {
	return s.RunCommand("start")
}

// StartForTest starts the app with the environment set to test.
// Returns the started module and stop function as a convenience.
func (s *Service[App]) StartForTest() (App, func()) {
	s.Env = EnvTest
	err := s.setup()
	if err != nil {
		panic(err)
	}
	go s.start()
	s.started.Wait()
	return s.root, s.Stop
}

// start calls Start on each module, in goroutines. Assumes that
// setup() has already been called. Start command must not block.
func (s *Service[App]) start() {
	for _, m := range s.modules {
		n := getModuleName(m)
		c := s.configs[n]
		BootPrintln("[service] starting", n)
		if c.Start != nil {
			c.Start()
		}
	}
	// mark as started
	s.started.Done()

	// mark process as running
	s.running.Add(1)

	// wait for a stop signal to be received
	// note: that c might crash without this parent goroutine knowing.
	s.wait()

	// mark process as done
	s.running.Done()
}

// wait blocks until a signal is received, or the stopper channel is closed
func (s *Service[App]) wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)
	select {
	case sig := <-c:
		BootPrintln("[service] got signal:", sig)
	case <-s.stopper:
		BootPrintln("[service] app stop")
	}
	s.stop()
}
