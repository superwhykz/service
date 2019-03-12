package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/service/cmd/gcp-sales-api/handlers"
	"github.com/ardanlabs/service/internal/platform/flag"
	"github.com/ardanlabs/service/internal/platform/gcp/ds"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"github.com/kelseyhightower/envconfig"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// It is set using build flags in the makefile.
var (
	version string
	commit  string
	repo    string
)

func main() {
	// =========================================================================
	// Logging
	// =========================================================================
	log := log.New(os.Stdout, "GCP-SALES-API : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	// =========================================================================
	var cfg struct {
		Project struct {
			ID string `default:"PROJECT_ID_NEEDED" envconfig:"PROJECT_ID"`
		}
		Trace struct {
			Prefix string `default:"gcp-sales-api"`
		}
		Profiler struct {
			Name string `default:"gcp-sales-api"`
		}
		Web struct {
			APIHost         string        `default:"0.0.0.0:8080" envconfig:"API_HOST"`
			ReadTimeout     time.Duration `default:"5s" envconfig:"READ_TIMEOUT"`
			WriteTimeout    time.Duration `default:"5s" envconfig:"WRITE_TIMEOUT"`
			ShutdownTimeout time.Duration `default:"5s" envconfig:"SHUTDOWN_TIMEOUT"`
		}
		Service struct {
			APIEnv string `default:"development" envconfig:"API_ENV"`
		}
	}

	if err := envconfig.Process("GCP-SALES-API", &cfg); err != nil {
		log.Fatalf("main : Parsing Config : %v", err)
	}

	if err := flag.Process(&cfg); err != nil {
		if err != flag.ErrHelp {
			log.Fatalf("main : Parsing Command Line : %v", err)
		}
		return // We displayed help.
	}

	// =========================================================================
	// Initialize Google Profiler
	// =========================================================================
	log.Println("main : Initializing Google Profiler...")
	if err := profiler.Start(profiler.Config{
		ProjectID:      cfg.Project.ID,
		Service:        cfg.Profiler.Name,
		ServiceVersion: version,
	}); err != nil {
		log.Fatalf("main : Failed Google Profile... : %v", err)
	}

	// =========================================================================
	// Start Tracing/Metrics Support. Create and register a OpenCensus Stackdriver Trace exporter.
	// =========================================================================
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:         cfg.Project.ID,
		MonitoredResource: monitoredresource.Autodetect(),
	})
	if err != nil {
		log.Fatalf("main : Failed to create the Stackdriver exporter : %v", err)
	}

	// Flush before main function exits
	defer exporter.Flush()

	// Export to Stackdriver Monitoring
	view.RegisterExporter(exporter)
	view.SetReportingPeriod(60 * time.Second)

	// Export to Stackdriver Trace.
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// =========================================================================
	// Initialize Google Pub/Sub. Example...
	// =========================================================================
	// log.Println("main : Initializing Google Pub/Sub...")
	// cPubSub, err := ps.New(cfg.Project.ID)
	// if err != nil {
	// 	log.Fatalf("main : Failed Initialize Pub/Sub... : %v", err)
	// }
	// defer cPubSub.Close()

	// =========================================================================
	// Initialize Google Firestore. Example...
	// =========================================================================
	// log.Println("TRACE MAIN : Initializing Google Firestore...")
	// cFirestore, err := fs.New(cfg.Project.ID)
	// if err != nil {
	// 	log.Fatalf("TRACE STARTUP : Failed Initialize Firestore... : %v", err)
	// }
	// defer cFirestore.Close()

	// =========================================================================
	// Initialize Google Datastore
	// =========================================================================
	log.Println("main : Initializing Google Datastore...")
	cDatastore, err := ds.New(cfg.Project.ID)
	if err != nil {
		log.Fatalf("main : Failed Initialize Datastore... : %v", err)
	}
	defer cDatastore.Close()

	// =========================================================================
	// App Starting
	// =========================================================================
	log.Printf("main : Started : GCP-SALES-API Initializing version %q : %q : %q", repo, version, commit)
	defer log.Println("main : Completed")

	cfgJSON, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatalf("main : Marshalling Config to JSON : %v", err)
	}

	// Validate what is being written to the logs.
	log.Printf("main : Config : %v\n", string(cfgJSON))

	// =========================================================================
	// Start Chameleon API Service
	// =========================================================================
	api := http.Server{
		Addr:           cfg.Web.APIHost,
		Handler:        handlers.API(log, cDatastore, cfg.Service.APIEnv),
		ReadTimeout:    cfg.Web.ReadTimeout,
		WriteTimeout:   cfg.Web.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API Listening %s", cfg.Web.APIHost)
		serverErrors <- api.ListenAndServe()
	}()

	// Shutdown
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Stop API Service
	// =========================================================================
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("main : Error starting server: %v", err)

	case <-osSignals:
		log.Println("main : Start shutdown...")

		// Create context for Shutdown call.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		if err := api.Shutdown(ctx); err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			if err := api.Close(); err != nil {
				log.Fatalf("main : Could not stop http server: %v", err)
			}
		}
	}

}
