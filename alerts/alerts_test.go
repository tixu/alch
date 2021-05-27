package alerts

import (
	"testing"
)

func TestPrometheusAlertHook_Validate(t *testing.T) {
	type fields struct {
		Version string
		Alerts  []Alert
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid",
			fields: fields{
				Version: "4",
				Alerts:  []Alert{{Status: "fake"}},
			},
			wantErr: false,
		},
		{
			name: "Not valid",
			fields: fields{
				Version: "3",
				Alerts:  []Alert{{Status: "fake"}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			pah := &Notification{
				Version: tt.fields.Version,
				Alerts:  tt.fields.Alerts,
			}
			if err := pah.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PrometheusAlertHook.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
