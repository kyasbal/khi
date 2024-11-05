package parameters

// Parameters types has pointer variables but &false,&"",...etc is not valid in Go. This function is just taking the address of literals for testing purpose.
func wrapPointer[T any](t T) *T {
	return &t
}
