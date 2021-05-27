package business

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/tixu/alch/alerts"
	"github.com/tixu/alch/cachethq"
	"github.com/tixu/alch/metrics"
)

type testCachetInstance struct {
	err error
}

func (tci *testCachetInstance) ChangeComponentStatus(name string, groupName string, status int) error {
	return tci.err
}

type testMetricsInstance struct{}

func (t *testMetricsInstance) IncreaseRequestCount()         {}
func (t *testMetricsInstance) IncreaseFailedRequestCount()   {}
func (t *testMetricsInstance) IncreaseComponentCount()       {}
func (t *testMetricsInstance) IncreaseFailedComponentCount() {}

func TestContext_ManageHook(t *testing.T) {
	type args struct {
		alert alerts.Alert
	}
	type fields struct {
		client  cachethq.Instance
		metrics metrics.Metrics
		logger  *logrus.Logger
	}
	logger := logrus.New()
	// Set JSON as default for the moment
	logger.SetFormatter(&logrus.JSONFormatter{})

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test",
			fields: fields{client: &testCachetInstance{}, metrics: &testMetricsInstance{}, logger: logger},
			args: args{
				alert: alerts.Alert{
					Labels: map[string]string{
						"tenant":  "ddd",
						"service": "ddd",
					},
					Status: "rrr",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if _, _, _, err := extractInfo(tt.args.alert); (err != nil) != tt.wantErr {
				t.Errorf("Context.ManageHook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
