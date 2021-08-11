package awclip
func emptyWhenNil(s *string) *string {
	if s == nil {
		return &empty
	}
	return s
}
