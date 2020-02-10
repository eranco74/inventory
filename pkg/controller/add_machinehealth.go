package controller

import (
	"github.com/eranco74/inventory/pkg/controller/machinehealth"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, machinehealth.Add)
}
