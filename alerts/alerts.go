package alerts

import (
	"errors"
	"time"

	gerr "github.com/tixu/alch/errors"
)

var ErrHookVersionNotSupported = errors.New("prometheus alert hook not supported (not in version 4)")

type Alert struct {
	Annotations  map[string]string `json:"annotations"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     time.Time         `json:"startsAt"`
	Status       string            `json:"status"`
}

type Notification struct {
	Alerts            []Alert           `json:"alerts"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	ExternalURL       string            `json:"externalURL"`
	GroupLabels       map[string]string `json:"groupLabels"`
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`

	// Timestamp records when the alert notification was received
	Timestamp string `json:"@timestamp"`
}

func (n *Notification) Validate() error {
	// Check hook version
	if n.Version != "4" {
		return gerr.NewBadInputError(ErrHookVersionNotSupported)
	}
	return nil
}
