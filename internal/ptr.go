package internal

// StringP returns a pointer string
func StringP(str string) *string { return &str }

// Float64P returns a pointer float
func Float64P(v float64) *float64 { return &v }
