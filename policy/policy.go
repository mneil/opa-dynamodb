package policy

import (
	"github.com/mneil/opa-dynamodb/store"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	log "github.com/sirupsen/logrus"
)

// Policy is a type that contains a store for policy data
type Policy struct {
	name  string
	store store.PolicyStore
}

// NewPolicy creates a new named policy
func NewPolicy(name string, store store.PolicyStore) *Policy {
	return &Policy{
		name:  name,
		store: store,
	}
}

// Get returns a policy given 2 terms. A namespace and a principal
func (p *Policy) Get(bctx rego.BuiltinContext, a, b *ast.Term) (*ast.Term, error) {
	var namespace, principal string
	if err := ast.As(a.Value, &namespace); err != nil {
		log.Error(err)
		return nil, err
	} else if ast.As(b.Value, &principal); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("Getting Policy from DynamoDB in %s for %s", namespace, principal)
	res, err := p.store.Get(namespace, principal)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	v, err := ast.InterfaceToValue(res)
	return ast.NewTerm(v), err
}
