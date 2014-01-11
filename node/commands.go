package main

import "fmt"

func nodeinfo (args []string) error {
  if storage != nil {
    fmt.Printf ("ID : %x\n", storage.Chord.Node.ID )
    fmt.Println ("RPC: ", storage.Chord.Node.Address ())
    fmt.Println ("DHT: ", storage.Chord.Node.HttpAddress ())
    fmt.Println ("   |_ Size: ", storage.Size (nil, nil))
    fmt.Println ()
    if storage.Chord.Predecessor != nil && storage.Chord.Successor != nil {
      fmt.Printf ("%x <- %x -> %x\n", storage.Chord.Predecessor.ID, storage.Chord.Node.ID, storage.Chord.Successor.ID)
    }
  }
  return nil
}

func leave (args []string) error {
  fmt.Println ("Leave!")
  storage.Leave (nil, nil)
  alive = false
  return nil
}
