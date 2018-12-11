package test

import (
	"os"
)

func init() {
	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}
}
