package types

import (
  "os"
  "os/signal"
  "syscall"
)

func (s *Sighandler_t) RegisterExitHandler (exitcallback func ()) {
  s.ExitCallback = exitcallback
  exitch := make (chan os.Signal, 1)
  signal.Notify (exitch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
  go func () {
    <-exitch
    s.ExitCallback ()
  }()
}

