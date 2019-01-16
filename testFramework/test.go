package testFramework

import (
  "github.com/mnordsletten/lotto/environment"
)

type Test struct {
  ClientCommandScript string `json:"clientcommandscript"`
	Setup               environment.SSHClients `json:"setup"`
	Cleanup             environment.SSHClients `json:"cleanup"`
}
