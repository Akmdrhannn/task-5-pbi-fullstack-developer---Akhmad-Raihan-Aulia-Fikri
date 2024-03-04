package main

import (
	route "btpn-golang/routers"
	"fmt"
)

func main() {
	route := route.SetupRouter()
	alamat := fmt.Sprintf("localhost:8080")
	fmt.Println("Server berjalan di alamat:", alamat)
	route.Run(alamat)

}
