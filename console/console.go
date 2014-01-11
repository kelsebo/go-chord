package console

import (
  "bufio"
  "errors"
  "fmt"
  "os"
  "strings"
)

//BEGIN Variables
var registeredCommands = map[string]*consoleCommand{}
//END Variables

//BEGIN Types
type commandCallback func(args []string) error
//END Types

//BEGIN Structs
type consoleCommand struct {
  description string
  callback    commandCallback
}
//END Structs

//BEGIN Functions
func byteArrayToStringArray (b []byte) []string {
  s := string(b)
  return strings.Fields(s)
}

func help() {
  fmt.Println("Available commands:")
  fmt.Println("\t\"help\",  Print this help")
  fmt.Println("\t\"exit\",  Exit console")
  for key, val := range registeredCommands {
    fmt.Println("\t\"" + key + "\",  " + val.description)
  }
}

func RegisterCommand (name string, description string, callback commandCallback) {
  if len(name) < 1 {
    panic (errors.New("Missing name"))
  }
  if name == "help" {
    panic (errors.New("Name 'help' is reserved"))
  }
  if name == "exit" {
    panic (errors.New("Name 'exit' is reserved"))
  }
  if callback == nil {
    panic (errors.New("Missing callback"))
  }
  cmd := new(consoleCommand)
  cmd.description = description
  cmd.callback = callback
  registeredCommands[name] = cmd
}

func RunConsole(exitcallback func()) {
  bio := bufio.NewReader(os.Stdin)
  fmt.Println ("Welcome, type 'help' for help")
  for {
    line, _, err := bio.ReadLine()
    if err != nil {
      continue
    }
    splat := byteArrayToStringArray(line)
    name := splat[0]
    if name == "help" {
      help()
      continue
    }
    if name == "exit" {
      break
    }

    f := registeredCommands[name]
    if f == nil {
      fmt.Println (name + ": no such command!")
      continue
    }

    err = f.callback (splat[1:])
    if err != nil {
      fmt.Println (err)
    }
  }
  exitcallback()
}
//END Functions
