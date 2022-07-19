package routine

// Cloneable interface to support copy itself.
type Cloneable interface {
	// Clone create and returns a copy of this object.
	Clone() any
}
