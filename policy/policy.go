// Copyright 2020 Michael Neil

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
