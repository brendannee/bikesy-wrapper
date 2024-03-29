package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "go.uber.org/fx"

    "blinktag.com/bikesy-wrapper/config"
    "blinktag.com/bikesy-wrapper/handlers"
    "blinktag.com/bikesy-wrapper/services"
)

// NewLogger ...
func NewLogger() *log.Logger {
    logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
    logger.Print("Executing NewLogger.")
    return logger
}

// NewMux ...
func NewMux(lc fx.Lifecycle, logger *log.Logger, config *config.Configuration) *http.ServeMux {
    logger.Print("Executing NewMux.")
    mux := http.NewServeMux()
    //logger.Print(":%v", config.Application.Port)
    server := &http.Server{
        Addr:    fmt.Sprintf(":%v", config.Application.Port),
        Handler: mux,
    }
    lc.Append(fx.Hook{
        OnStart: func(context.Context) error {
            logger.Print("Starting HTTP server on port ", config.Application.Port)
            go server.ListenAndServe()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Print("Stopping HTTP server.")
            return server.Shutdown(ctx)
        },
    })

    return mux
}

// Register ...
func Register(mux *http.ServeMux, logger *log.Logger, routeService services.RouteService, elevationService services.ElevationService) {
    h := handlers.NewHealthHandler(logger)
    r := handlers.NewBikesyHandler(logger, routeService, elevationService)
    hHandler := h.Handler()
    rHandler := r.Handler()
    mux.Handle("/health", hHandler)
    mux.Handle("/route", rHandler)
}

func main() {
    app := fx.New(
        fx.Provide(
            NewLogger,
            NewMux,

            config.LoadConfig,
            services.NewRouteService,
            services.NewElevationService,
        ),
        fx.Invoke(Register),
    )
    app.Run()
}
