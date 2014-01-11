package types

import (
  "fmt"
  "io"
  "strconv"
  "crypto/sha1"
//  "hash/fnv"
)

func hashfunc(in string) ID_t {
	hasher := sha1.New()
	io.WriteString(hasher, in)
  sum := hasher.Sum(nil)
  tmp := sum[(len(sum)-(THE_M/8)):]
  res, err := strconv.ParseUint(fmt.Sprintf("%x",tmp), 16, THE_M)
  if err != nil {
    fmt.Println(err)
  }
  return ID_t(res)
}
/*
func hashfunc (in string) ID_t {
  hash := fnv.New32 ()
  hash.Write ([]byte (in))
  return ID_t (hash.Sum32 ())
}
*/
//END Functions

