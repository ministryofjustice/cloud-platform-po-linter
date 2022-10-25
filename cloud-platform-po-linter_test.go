package main

import (
	"testing"
)

func Test_poLint(t *testing.T) {
	type args struct {
		dir []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				dir: []string{"template/00-success-prometheus.yaml"},
			},
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				dir: []string{"template/00-failure-prometheus.yaml"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := poLint(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("poLint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
