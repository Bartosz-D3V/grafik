// Package ds (Data Structure) contains all golang data structures used by generator.
package ds

// Const represents simplified constant in Golang AST.
// Name is the name of the constant.
// Val is an interface{} and represents any possible value of the struct.
type Const struct {
	Name string
	Val  interface{}
}

// IsString returns true if passed interface{} is a string.
func (c Const) IsString() bool {
	switch c.Val.(type) {
	case string:
		return true
	default:
		return false
	}
}
