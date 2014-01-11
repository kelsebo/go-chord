package netutils

import (
  "net"
  "net/http"
  "net/rpc"
)

//BEGIN RPC FUNCTIONS

/**
 * Same as rpc.Register, but wil panic if anything goes wrong
 */
func RPCRegister (rcvr interface{}) bool {
  err := rpc.Register (rcvr)
  if err != nil {
    return false
  }
  return true
}

/**
 * Same as rpc.RegisterName, but wil panic if anything goes wrong
 */
func RPCRegisterName (name string, rcvr interface{}) bool {
  err := rpc.RegisterName (name, rcvr)
  if err != nil {
    return false
  }
  return true
}

/**
 * Open RPC connection with host to see if its serving RPC and close connection immediately
 */
func RPCAlive (addr string) bool {
  client, err := rpc.DialHTTP ("tcp4", addr)
  if err != nil {
    return false
  }
  client.Close()
  return true
}

/**
 *  Open RPC connection with host and return client
 */
func RPCConnect (addr string) *rpc.Client {
  conn, err := rpc.DialHTTP("tcp4", addr)
  if err != nil {
    return nil
  }
  return conn
}

/**
 * Serve RPC on the given listener. Should be called in a goroutine
 */
func RPCServe (listener *net.TCPListener) {
  rpc.HandleHTTP()
  http.Serve(listener, nil)
}

/**
 * Connect to RPC server and do a call, close connection and return result
 */
func RPCSingleCall (addr string, serviceMethod string, args interface{}, reply interface{}) error {
  conn, err := rpc.DialHTTP("tcp4", addr)
  if err != nil {
    return err
  }
  defer conn.Close()

  call := <-conn.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done
  return call.Error
}

//END RPC FUNCTIONS
