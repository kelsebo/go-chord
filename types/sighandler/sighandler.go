package sighandler

import "types"

func New () *types.Sighandler_t {
  return new (types.Sighandler_t)
}
