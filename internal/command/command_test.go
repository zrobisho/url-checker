package command

import (
	"reflect"
	"testing"
	"time"
)

func TestParseInput(t *testing.T) {
	type args struct {
		input CLIInput
	}
	tests := []struct {
		name    string
		args    args
		want    ApplicationArgs
		wantErr bool
	}{
		{
			name: "no upstream input",
			args: args{
				input: CLIInput{
					UpstreamUrls:    "",
					PrometheusPort:  9090,
					CheckInterval:   1,
					UpstreamTimeout: 1,
				},
			},
			want:    ApplicationArgs{},
			wantErr: true,
		},
		{
			name: "multiple urls",
			args: args{
				input: CLIInput{
					UpstreamUrls:    "https://test1.com,https://test2.com,https://test3.com",
					PrometheusPort:  9090,
					CheckInterval:   1,
					UpstreamTimeout: 1,
				},
			},
			want: ApplicationArgs{
				PrometheusPort:  9090,
				UpstreamUrls:    []string{"https://test1.com", "https://test2.com", "https://test3.com"},
				CheckInterval:   time.Duration(1) * time.Second,
				UpstreamTimeout: time.Duration(1) * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid urls",
			args: args{
				input: CLIInput{
					UpstreamUrls:    "\n",
					PrometheusPort:  9090,
					CheckInterval:   1,
					UpstreamTimeout: 1,
				},
			},
			want:    ApplicationArgs{},
			wantErr: true,
		},
		{
			name: "invalid interval",
			args: args{
				input: CLIInput{
					UpstreamUrls:    "https://test1.com,https://test2.com,https://test3.com",
					PrometheusPort:  9090,
					CheckInterval:   0,
					UpstreamTimeout: 1,
				},
			},
			want:    ApplicationArgs{},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			args: args{
				input: CLIInput{
					UpstreamUrls:    "https://test1.com,https://test2.com,https://test3.com",
					PrometheusPort:  9090,
					CheckInterval:   1,
					UpstreamTimeout: 0,
				},
			},
			want:    ApplicationArgs{},
			wantErr: true,
		},
		{
			name: "invalid prometheus port: zero",
			args: args{
				input: CLIInput{
					UpstreamUrls:    "https://test1.com,https://test2.com,https://test3.com",
					PrometheusPort:  0,
					CheckInterval:   1,
					UpstreamTimeout: 1,
				},
			},
			want:    ApplicationArgs{},
			wantErr: true,
		},
		{
			name: "invalid prometheus port: too large",
			args: args{
				input: CLIInput{
					UpstreamUrls:    "https://test1.com,https://test2.com,https://test3.com",
					PrometheusPort:  100000,
					CheckInterval:   1,
					UpstreamTimeout: 1,
				},
			},
			want:    ApplicationArgs{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInput(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
