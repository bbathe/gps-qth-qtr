package main

import (
	"testing"
)

func Test_parseGGA(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name:    "Valid 1",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,1,12,0.96,250.6,M,-33.4,M,,*7A"},
			want:    "GPS fix (SPS)",
			want1:   12,
			wantErr: false,
		},
		{
			name:    "Valid 2",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,1,3,0.96,250.6,M,-33.4,M,,*4A"},
			want:    "GPS fix (SPS)",
			want1:   3,
			wantErr: false,
		},
		{
			name:    "Valid 3",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,2,7,0.96,250.6,M,-33.4,M,,*4D"},
			want:    "DGPS fix",
			want1:   7,
			wantErr: false,
		},
		{
			name:    "Valid 4",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,0,99,0.96,250.6,M,-33.4,M,,*78"},
			want:    "invalid",
			want1:   99,
			wantErr: false,
		},
		{
			name:    "Not enough fields",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,0,"},
			want:    "",
			want1:   0,
			wantErr: true,
		},
		{
			name:    "Bad checksum",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,1,12,0.96,250.6,M,-33.4,M,,*70"},
			want:    "",
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ttt := tt

		t.Run(ttt.name, func(t *testing.T) {
			got, got1, err := parseGGA(ttt.args.s)
			if (err != nil) != ttt.wantErr {
				t.Errorf("parseGGA() error = %v, wantErr %v", err, ttt.wantErr)
				return
			}
			if got != ttt.want {
				t.Errorf("parseGGA() got = %v, want %v", got, ttt.want)
			}
			if got1 != ttt.want1 {
				t.Errorf("parseGGA() got1 = %v, want %v", got1, ttt.want1)
			}
		})
	}
}
