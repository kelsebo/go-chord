package types

func (s *Slave_t) Kill (ign string, _ign *string) error {
  s.KillFunc ()
  return nil
}
func (s *Slave_t) Leave (ign string, _ign *string) error {
  s.LeaveFunc ()
  return nil
}
