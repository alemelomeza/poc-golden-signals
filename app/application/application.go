package application

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

var (
	counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	histogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "myapp_greeting_seconds",
		Help:    "Time take to greet someone",
		Buckets: []float64{1, 2, 5, 6, 10},
	}, []string{"code"})
)

func StartApp() {
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	log.WithFields(logrus.Fields{
		"payload": "myPayload",
	}).Info("myMessage")

	e := echo.New()
	e.HideBanner = true
	e.Debug = false
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "/")
	})
	e.GET("/counter", func(c echo.Context) error {
		counter.Inc()
		return c.String(http.StatusOK, "counter")
	})
	e.GET("/histogram", func(c echo.Context) error {
		start := time.Now()
		defer func() {
			httpDuration := time.Since(start)
			histogram.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(httpDuration.Seconds())
		}()
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(2)
		time.Sleep(time.Duration(n) * time.Second)
		return c.String(http.StatusOK, "histogram")
	})
	e.GET("/metrics", func(c echo.Context) error {
		h := promhttp.Handler()
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	s := &http.Server{
		Addr:         ":1313",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		if err := e.StartServer(s); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
