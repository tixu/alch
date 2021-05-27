package business

import (
	"github.com/andygrunwald/cachet"
	"github.com/sirupsen/logrus"
	"github.com/tixu/alch/alerts"
	"github.com/tixu/alch/cachethq"
	"github.com/tixu/alch/config"
	"github.com/tixu/alch/metrics"
)

func (ctx *business) HandleNotification(msg alerts.Notification) error {
	for i, alert := range msg.Alerts {
		ctx.logger.Infof("alert: %d", i)
		group, service, status, err := extractInfo(alert)
		if err != nil {
			continue
		}
		ctx.logger.Infof("modififying component %s in group %s", service, group)
		ctx.metrics.IncreaseComponentCount()
		error := ctx.client.ChangeComponentStatus(service, group, status)
		if error != nil {
			ctx.metrics.IncreaseFailedComponentCount()
		}
	}
	return nil
}

func extractInfo(alert alerts.Alert) (group string, service string, status int, err error) {
	group = alert.Labels["tenant"]
	service = alert.Labels["service"]
	status = cachet.ComponentStatusOperational
	if alert.Status == "firing" {
		status = cachet.ComponentStatusPartialOutage
	}
	return group, service, status, nil
}

type Business interface {
	HandleNotification(msg alerts.Notification) error
}

// instance CachetHQ instance
type business struct {
	client  cachethq.Instance
	logger  *logrus.Logger
	metrics metrics.Metrics
}

func NewInstanceConf(conf *config.Config, m metrics.Metrics, logger *logrus.Logger) (*business, error) {
	cachetHQ, _ := cachethq.NewInstanceConf(conf.Cachet, logger)
	return &business{client: cachetHQ, metrics: m, logger: logger}, nil
}
