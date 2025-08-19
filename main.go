package main

import (
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.Module).Run()
}
