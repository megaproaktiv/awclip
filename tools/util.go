package tools

var empty = ""

func EmptyWhenNil(s *string) *string {
	if s == nil {
		return &empty
	}
	return s
}
