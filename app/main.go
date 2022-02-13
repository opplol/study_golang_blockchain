package main

import (
	"go_crypo_coin/explorer"
	"go_crypo_coin/rest"
)

func main() {
	go explorer.Start(4000)
	rest.Start(5000)
}
