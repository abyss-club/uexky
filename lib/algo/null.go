package algo

func NullString(s string) *string {
	return &s
}

func NullToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func NullInt(i int) *int {
	return &i
}

func NullToIntDefault(i *int, def int) int {
	if i == nil {
		return def
	}
	return *i
}
