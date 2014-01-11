package organizer

import (
  "fmt"
  "types"
)

func New(orgaddr string, orgport int) *types.Organizer_t {
  o := types.Organizer_t(fmt.Sprintf("%s:%d", orgaddr, orgport))
  return &o
}
