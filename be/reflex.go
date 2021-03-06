// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package be

import (
	. "github.com/gocircuit/escher/circuit"
	// . "github.com/gocircuit/escher/faculty"
)

// Reflex is a bundle of not yet attached sense endpoints (synapses).
type Reflex map[Name]*Synapse

//
type Materializer func(Reflex, Circuit) interface{}

// Material represents a materializable object implemented as a Go type.
type Material interface {
	Spark(eye *Eye, matter Circuit, aux ...interface{}) Value // Initializer
}

type Sparkless struct{}

func (Sparkless) Spark(eye *Eye, matter Circuit, aux ...interface{}) Value {
	return nil
}
