package main

import (
	"os"
	"runtime"
	"time"

	"math/rand"

	"github.com/Ranper/iam/internal/apiserver"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	apiserver.NewApp("iam-apiserver").Run()
}
