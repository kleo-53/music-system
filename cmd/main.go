package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kleo-53/music-system/config"
	"github.com/kleo-53/music-system/internal/app"
	"github.com/kleo-53/music-system/pkg/logger"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

//	@title			Music library
//	@version		0.0.1
//	@description	This server provides information about songs in music library

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Log().Fatal(context.Background(), "Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
