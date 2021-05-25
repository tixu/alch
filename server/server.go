package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/tixu/alch/alerts"
	"github.com/tixu/alch/business"
	"github.com/tixu/alch/config"
	"github.com/tixu/alch/errors"
)

var (
	requestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_request_count",
		Help: "total request count",
	})
	failedRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_failed_request_count",
		Help: "failed request count",
	})
	responseLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "go_response_latency",
		Help: "response latencies",
	})
)

func (ctx *server) prom(w http.ResponseWriter, r *http.Request) {

	requestReceived := time.Now()
	defer func() {
		responseLatency.Observe(time.Since(requestReceived).Seconds())
	}()
	// [END monitoring_sli_metrics_prometheus_latency]
	// [START monitoring_sli_metrics_prometheus_counts]
	requestCount.Inc()

	ctx.logger.Infof("receive Request")
	if r.Body == nil {
		failedRequestCount.Inc()
		err := errors.NewEmptyBodyError(nil)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		failedRequestCount.Inc()
		errs := errors.NewInternalServerError(err)

		http.Error(w, errs.Error(), http.StatusInternalServerError)

		return
	}
	defer r.Body.Close()

	msg := alerts.Notification{}
	err = json.Unmarshal(b, &msg)

	if err != nil {
		failedRequestCount.Inc()
		ctx.logger.Error(err)
		errs := errors.NewBadInputError(err)
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}
	err = msg.Validate()
	if err != nil {
		errs := errors.NewBadInputError(err)
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}

	ctx.logger.Infof(msg.Status)
	worker, _ := business.NewInstanceConf(ctx.conf, ctx.logger)
	worker.ChangeComponentStatus(msg)
}

func NewServer(cfg *config.Config, logger *logrus.Logger) *server {
	router := mux.NewRouter()
	addr := fmt.Sprintf("%s:%d", cfg.Server.ListenAddr, cfg.Server.Port)
	httpserver := &http.Server{
		Handler: router,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s := &server{server: httpserver, logger: logger, conf: cfg}
	router.HandleFunc("/webhook/prometheus", s.prom)
	router.Handle("/metrics", promhttp.Handler())
	return s
}

type Server interface {
	Run()
}

type server struct {
	logger *logrus.Logger
	server *http.Server
	conf   *config.Config
}

func (ctx *server) Run() {
	log.Fatal(ctx.server.ListenAndServe())
}
