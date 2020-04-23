package main

import (
	"os"

	"DBproject1/boot"
)

func main() {
	if err := boot.Boot(make(chan os.Signal, 1), os.Getenv("FC_MODE") == "test"); err != nil {
		panic(err)
	}
}
