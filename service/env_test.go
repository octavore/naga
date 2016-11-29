package service

import (
	"os"
	"testing"
)

func TestGetEnvironment(t *testing.T) {
	for e, envs := range EnvMap {
		for _, env := range envs {
			os.Setenv(EnvVarName, env)
			if e != GetEnvironment() {
				t.Errorf("expected %q but got %q", e, GetEnvironment())
			}
		}
	}
}
