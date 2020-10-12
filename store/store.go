package store

// PolicyStore is an interface type for working with different policy backends
type PolicyStore interface {
	Get(namespace string, principal string) (interface{}, error)
}
