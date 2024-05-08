package synchronizer

func (s *storage) Get(key string) *property {
	return nil
}
func (s *storage) Keys() []string {
	return nil
}
func (s *storage) Count() int {
	return 0
}
func (s *storage) Set(key string, value []byte, reliable bool) error {
	return nil
}
func (s *storage) Upsert(key string, reliable bool, cb func(exists bool, current []byte) []byte) error {
	return nil
}
func (s *storage) Remove(key string) bool {
	return false
}
