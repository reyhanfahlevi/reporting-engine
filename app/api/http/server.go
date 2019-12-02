package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reyhanfahlevi/pkg/go/log"
	"github.com/tokopedia/reporting-engine/app/api"
	"github.com/tokopedia/reporting-engine/app/api/http/reporting"
)

// Server struct
type Server struct {
	http.Server
	reportingHandler *reporting.Handler
}

// New will instantiate the http server
func New(ucReporting api.ReportingUseCase) *Server {
	reportingHandler := reporting.New(ucReporting)
	return &Server{
		reportingHandler: reportingHandler,
	}
}

func (s *Server) v1Endpoint(r *gin.Engine) {
	g := r.Group("/v1")

	/* reporting endpoint */

	rg := g.Group("/reporting")
	s.reportingHandler.HandlerStoreReport(rg)
	s.reportingHandler.HandlerGetReport(rg)
}

// Run run http server
func (s *Server) Run(port string) {
	r := gin.Default()

	/* Register router here */

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "You know, for reporting...")
	})

	/* add middleware here for authentication */

	/* v1 endpoint */
	s.v1Endpoint(r)

	/* End of registering router */

	s.Addr = port
	s.Handler = r

	go func() {
		// service connections
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Graceful Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
