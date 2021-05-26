package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"time"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/tixu/alch/alerts"
	"github.com/tixu/alch/business"
	"github.com/tixu/alch/config"
	"github.com/tixu/alch/errors"
	"github.com/tixu/alch/metrics"
)

func (ctx *server) prom(w http.ResponseWriter, r *http.Request) {

	// [END monitoring_sli_metrics_prometheus_latency]
	// [START monitoring_sli_metrics_prometheus_counts]
	ctx.metrics.IncreaseRequestCount()

	ctx.logger.Infof("receive Request")
	if r.Body == nil {
		ctx.metrics.IncreaseFailedRequestCount()
		err := errors.NewEmptyBodyError(nil)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.metrics.IncreaseFailedRequestCount()
		errs := errors.NewInternalServerError(err)
		http.Error(w, errs.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	msg := alerts.Notification{}
	err = json.Unmarshal(b, &msg)

	if err != nil {
		ctx.metrics.IncreaseFailedRequestCount()
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
	worker, _ := business.NewInstanceConf(ctx.conf, ctx.metrics, ctx.logger)
	worker.HandleNotification(msg)
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
	m := metrics.NewInstance()
	s := &server{server: httpserver, logger: logger, conf: cfg, metrics: m}
	router.HandleFunc("/webhook/prometheus", s.prom)
	router.Handle("/metrics", promhttp.Handler())
	return s
}

type Server interface {
	Run()
}

type server struct {
	logger  *logrus.Logger
	server  *http.Server
	conf    *config.Config
	metrics metrics.Metrics
}

func (ctx *server) Run() {
	log.Fatal(ctx.server.ListenAndServe())
}
