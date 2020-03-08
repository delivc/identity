package namespace

var namespace string

// SetNamespace sets a new namespace
func SetNamespace(ns string) {
	namespace = ns
}

// GetNamespace returns the namespace
func GetNamespace() string {
	return namespace
}
