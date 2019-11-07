package boot

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/vilbert/go-skeleton/internal/config"

	skeletonData "github.com/vilbert/go-skeleton/internal/data/skeleton"
	skeletonServer "github.com/vilbert/go-skeleton/internal/delivery/http"
	skeletonHandler "github.com/vilbert/go-skeleton/internal/delivery/http/skeleton"
	skeletonService "github.com/vilbert/go-skeleton/internal/service/skeleton"
)

// HTTP will load configuration, do dependency injection and then start the HTTP server
func HTTP() error {
	var (
		s   skeletonServer.Server    // HTTP Server Object
		sd  skeletonData.Data        // BridgingProduct domain data layer
		ss  skeletonService.Service  // BridgingProduct domain service layer
		sh  *skeletonHandler.Handler // BridgingProduct domain handler
		cfg *config.Config           // Configuration object
	)

	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	cfg = config.Get()
	// Open MySQL DB Connection
	db, err := sqlx.Open("mysql", cfg.Database.Master)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	// BridgingProduct domain init
	sd = skeletonData.New(db)
	ss = skeletonService.New(sd)
	sh = skeletonHandler.New(ss)

	// Inject service used on handler
	s = skeletonServer.Server{
		Skeleton: sh,
	}

	if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		return err
	}

	return nil
}
