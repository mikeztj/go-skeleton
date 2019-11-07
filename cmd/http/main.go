package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vilbert/ProdukOnline/internal/boot"
)

func main() {
	if err := boot.HTTP(); err != nil {
		fmt.Println("[HTTP] failed to boot http server due to " + err.Error())
	}
}
