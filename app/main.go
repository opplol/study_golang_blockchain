package main

import (
	"go_crypo_coin/cli"
	"go_crypo_coin/db"
)

func main() {
	defer db.Close()
	cli.Start()

}
