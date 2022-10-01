package main

import (
	"github.com/karim-w/ksec/cmd"
	"github.com/karim-w/ksec/service"
)

func main() {
	go service.GetKubeClient()
	cmd.Execute()
}
