package utils

// PtrOrNil returns a pointer to the string if it's not empty, otherwise nil
func PtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntPtrOrNil returns a pointer to the int if it's not zero, otherwise nil
func IntPtrOrNil(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// UintPtrOrNil returns a pointer to the uint if it's not zero, otherwise nil
func UintPtrOrNil(i uint) *uint {
	if i == 0 {
		return nil
	}
	return &i
}