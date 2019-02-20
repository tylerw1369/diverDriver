package common

import (
	"sync"

	"github.com/tylerw1369/iotago"
)

const (
	DiverDriverVersion = "0.0.1"
)

type PowFuncDefinition func(trytes giota.Trytes, minWeightMagnitude int) (result giota.Trytes, Error error)
type GetPowInfoDefinition func(ServerVersion string, PowType string, PowVersion string, Error error)

type ClientAPI struct {
	PowFuncDefinition    PowFuncDefinition
	GetPowInfoDefinition GetPowInfoDefinition
}
