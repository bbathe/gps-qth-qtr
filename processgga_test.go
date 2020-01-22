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
		want2   float64
		wantErr bool
	}{
		{
			name:    "Valid 1",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,1,12,0.96,250.6,M,-33.4,M,,*7A"},
			want:    "GPS fix (SPS)",
			want1:   12,
			want2:   0.96,
			wantErr: false,
		},
		{
			name:    "Valid 2",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,1,3,99.3,250.6,M,-33.4,M,,*46"},
			want:    "GPS fix (SPS)",
			want1:   3,
			want2:   99.3,
			wantErr: false,
		},
		{
			name:    "Valid 3",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,2,7,1.33,250.6,M,-33.4,M,,*43"},
			want:    "DGPS fix",
			want1:   7,
			want2:   1.33,
			wantErr: false,
		},
		{
			name:    "Valid 4",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,0,99,5.2,250.6,M,-33.4,M,,*40"},
			want:    "invalid",
			want1:   99,
			want2:   5.2,
			wantErr: false,
		},
		{
			name:    "Not enough fields",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,0,"},
			want:    "",
			want1:   0,
			want2:   0.0,
			wantErr: true,
		},
		{
			name:    "Bad checksum",
			args:    args{s: "GNGGA,013016.00,7751.3,S,16642.4,E,1,12,0.96,250.6,M,-33.4,M,,*70"},
			want:    "",
			want1:   0,
			want2:   0.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ttt := tt

		t.Run(ttt.name, func(t *testing.T) {
			got, got1, got2, err := parseGGA(ttt.args.s)
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
			if got2 != ttt.want2 {
				t.Errorf("parseGGA() got2 = %v, want %v", got2, ttt.want2)
			}
		})
	}
}
