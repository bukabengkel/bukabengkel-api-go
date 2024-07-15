package config

import (
	"path/filepath"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewCasbinEnfocer will create new casbin enfocer instance
func NewCasbinEnfocer(config *Config) (enfocer *casbin.Enforcer, err error) {
	var m model.Model
	var a persist.Adapter

	modelPath, err := filepath.Abs(config.CasbinModelFilePath)
	if err != nil {
		return nil, err
	}

	policyPath, err := filepath.Abs(config.CasbinPolicyFilePath)
	if err != nil {
		return nil, err
	}

	// Load the model from file
	if m, err = model.NewModelFromFile(modelPath); err != nil {
		return nil, err
	}
	// Load the policy from file
	a = fileadapter.NewAdapter(policyPath)

	// Initialize the enforcer
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	return e, nil
}
