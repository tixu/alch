package business

import (
	"github.com/andygrunwald/cachet"
	"github.com/sirupsen/logrus"
	"github.com/tixu/alch/alerts"
	"github.com/tixu/alch/cachethq"
	"github.com/tixu/alch/config"
)

func (ctx *business) ChangeComponentStatus(msg alerts.Notification) error {
	for i, alert := range msg.Alerts {
		ctx.logger.Infof("alert: %d", i)
		group := alert.Labels["tenant"]
		service := alert.Labels["service"]
		status := alert.Status
		ctx.logger.Infof("modififying component %s in group %s", service, group)
		hqstatus := cachet.ComponentStatusOperational
		if status == "firing" {
			hqstatus = cachet.ComponentStatusPartialOutage
		}
		ctx.logger.Infof("modififying component %s in group %s", service, group, hqstatus)
		ctx.client.ChangeComponentStatus(service, group, hqstatus)
	}
	return nil
}

type Business interface {
	ChangeComponentStatus(msg alerts.Notification) error
}

// instance CachetHQ instance
type business struct {
	client cachethq.Instance
	logger *logrus.Logger
}

func NewInstanceConf(conf *config.Config, logger *logrus.Logger) (*business, error) {
	cachetHQ, _ := cachethq.NewInstanceConf(conf.Cachet, logger)
	return &business{client: cachetHQ, logger: logger}, nil
}
