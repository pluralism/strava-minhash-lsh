package stravaminhashlsh

type set struct {
	m map[string]struct{}
}

func newSet() *set {
	s := &set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *set) Add(value string) {
	s.m[value] = struct{}{}
}

func (s *set) Remove(value string) {
	delete(s.m, value)
}

func (s *set) Contains(value string) bool {
	_, ok := s.m[value]
	return ok
}
