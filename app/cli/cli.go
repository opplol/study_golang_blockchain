package cli

import (
	"flag"
	"fmt"
	"go_crypo_coin/explorer"
	"go_crypo_coin/rest"
	"os"
	"runtime"
)

func usage() {
	fmt.Printf("Welecome to go crypto\n")
	fmt.Printf("Please use the floowing flags:\n\n")
	fmt.Printf("-port=4000:  Set the PORT of the server\n")
	fmt.Printf("-mode=rest:  Choose between 'html' and 'rest'\n")
	runtime.Goexit()
	os.Exit(0)
}
func Start() {
	if len(os.Args) == 0 {
		usage()
	}

	port := flag.Int("port", 4000, "Set Port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")
	flag.Parse()

	fmt.Println(*port, *mode)

	switch *mode {
	case "html":
		fmt.Println("Start Explorer")
		explorer.Start(*port)
	case "rest":
		fmt.Println("Start REST API")
		rest.Start(*port)
	default:
		usage()
	}
}
