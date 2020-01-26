package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseRMCTime(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "2020-01-18 20:34:34 +0000 UTC",
			args:    args{fields: []string{"", "203434.00", "", "", "", "", "", "", "", "180120"}},
			want:    time.Date(2020, time.Month(1), 18, 20, 34, 34, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ttt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRMCTime(ttt.args.fields)
			if (err != nil) != ttt.wantErr {
				t.Errorf("parseRMCTime() error = %v, wantErr %v", err, ttt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, ttt.want) {
				t.Errorf("parseRMCTime() = %v, want %v", got, ttt.want)
			}
		})
	}
}

func Test_parseDegMinToFloat(t *testing.T) {
	type args struct {
		dm string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "Budapest Latitude",
			args:    args{dm: "4726.5824"},
			want:    47.44304,
			wantErr: false,
		},
		{
			name:    "Budapest Longitude",
			args:    args{dm: "01900.0581"},
			want:    19.00096833333333333,
			wantErr: false,
		},
		{
			name:    "McMurdo Station Longitude",
			args:    args{dm: "16642.4"},
			want:    166.70666666666667,
			wantErr: false,
		},
		{
			name:    "Nowhere 1",
			args:    args{dm: "18100.0"},
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "Nowhere 2",
			args:    args{dm: "9100.0"},
			want:    0.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ttt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDegMinToFloat(ttt.args.dm)
			if (err != nil) != ttt.wantErr {
				t.Errorf("parseDegMinToFloat() error = %v, wantErr %v", err, ttt.wantErr)
				return
			}
			if got != ttt.want {
				t.Errorf("parseDegMinToFloat() = %v, want %v", got, ttt.want)
			}
		})
	}
}

func Test_latLonToGridsquare(t *testing.T) {
	type args struct {
		lat float64
		lon float64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Budapest",
			args:    args{lat: 47.44304, lon: 19.000968},
			want:    "JN97mk",
			wantErr: false,
		},
		{
			name:    "Rio De Janeiro",
			args:    args{lat: -22.912328, lon: -43.182617},
			want:    "GG87jc",
			wantErr: false,
		},
		{
			name:    "Washington DC",
			args:    args{lat: 38.92, lon: -77.065},
			want:    "FM18lw",
			wantErr: false,
		},
		{
			name:    "McMurdo Station",
			args:    args{lat: -77.855000, lon: 166.706667},
			want:    "RB32id",
			wantErr: false,
		},
		{
			name:    "South Pole",
			args:    args{lat: -90.000000, lon: 0.0},
			want:    "JA00aa",
			wantErr: false,
		},
		{
			name:    "North Pole",
			args:    args{lat: 90.000000, lon: 0.0},
			want:    "JS00aa",
			wantErr: false,
		},
		{
			name:    "Equator West",
			args:    args{lat: 0.0, lon: 0.0},
			want:    "JJ00aa",
			wantErr: false,
		},
		{
			name:    "Equator East 1",
			args:    args{lat: 0.0, lon: 180.0},
			want:    "SJ00aa",
			wantErr: false,
		},
		{
			name:    "Equator East 2",
			args:    args{lat: 0, lon: -180.0},
			want:    "AJ00aa",
			wantErr: false,
		},
		{
			name:    "Lost 1",
			args:    args{lat: 90.0, lon: 180.0},
			want:    "SS00aa",
			wantErr: false,
		},
		{
			name:    "Lost 2",
			args:    args{lat: 90.0, lon: -180.0},
			want:    "AS00aa",
			wantErr: false,
		},
		{
			name:    "Lost 3",
			args:    args{lat: -90.0, lon: 180.0},
			want:    "SA00aa",
			wantErr: false,
		},
		{
			name:    "Lost 4",
			args:    args{lat: -90.0, lon: -180.0},
			want:    "AA00aa",
			wantErr: false,
		},
		{
			name:    "Nowhere",
			args:    args{lat: 91.0, lon: 181.0},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ttt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := latLonToGridsquare(ttt.args.lat, ttt.args.lon)
			if (err != nil) != ttt.wantErr {
				t.Errorf("latLonToGridsquare() error = %v, wantErr %v", err, ttt.wantErr)
				return
			}
			if got != ttt.want {
				t.Errorf("latLonToGridsquare() = %v, want %v", got, ttt.want)
			}
		})
	}
}

func Test_parseRMC(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		want1   string
		want2   float64
		want3   float64
		wantErr bool
	}{
		{
			name:    "Budapest",
			args:    args{s: "GNRMC,203434.00,A,4726.5824,N,01900.0581,E,0.149,,180120,,,A*62"},
			want:    time.Date(2020, time.Month(1), 18, 20, 34, 34, 0, time.UTC),
			want1:   "JN97mk",
			want2:   47.44304,
			want3:   19.000968333333333,
			wantErr: false,
		},
		{
			name:    "Rio De Janeiro",
			args:    args{s: "GNRMC,203434.00,A,2254.7397,S,04310.957,W,0.149,,300333,,,A*59"},
			want:    time.Date(2033, time.Month(3), 30, 20, 34, 34, 0, time.UTC),
			want1:   "GG87jc",
			want2:   -22.912328333333335,
			want3:   -43.18261666666667,
			wantErr: false,
		},
		{
			name:    "Washington DC",
			args:    args{s: "GNRMC,020202.00,A,3855.2,N,07703.9,W,0.149,,180120,,,A*7B"},
			want:    time.Date(2020, time.Month(1), 18, 2, 2, 2, 0, time.UTC),
			want1:   "FM18lw",
			want2:   38.92,
			want3:   -77.065,
			wantErr: false,
		},
		{
			name:    "McMurdo Station",
			args:    args{s: "GNRMC,203434.00,A,7751.3,S,16642.4,E,0.149,,180120,,,A*73"},
			want:    time.Date(2020, time.Month(1), 18, 20, 34, 34, 0, time.UTC),
			want1:   "RB32id",
			want2:   -77.855000,
			want3:   166.70666666666668,
			wantErr: false,
		},
		{
			name:    "Bad number of fields",
			args:    args{s: "GNRMC,203434.00,A,7751.3,S,16642.4,E,0.149,,"},
			want:    time.Time{},
			want1:   "",
			want2:   0.0,
			want3:   0.0,
			wantErr: true,
		},
		{
			name:    "Bad checksum",
			args:    args{s: "GNRMC,203434.00,A,7751.3,S,16642.4,E,0.149,,180120,,,A*74"},
			want:    time.Time{},
			want1:   "",
			want2:   0.0,
			want3:   0.0,
			wantErr: true,
		},
		{
			name:    "Bad state",
			args:    args{s: "GNRMC,203434.00,V,7751.3,S,16642.4,E,0.149,,180120,,,A*73"},
			want:    time.Time{},
			want1:   "",
			want2:   0.0,
			want3:   0.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ttt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, err := parseRMC(ttt.args.s)
			if (err != nil) != ttt.wantErr {
				t.Errorf("parseRMC() error = %v, wantErr %v", err, ttt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, ttt.want) {
				t.Errorf("parseRMC() got = %v, want %v", got, ttt.want)
			}
			if got1 != ttt.want1 {
				t.Errorf("parseRMC() got1 = %v, want %v", got1, ttt.want1)
			}
			if got2 != ttt.want2 {
				t.Errorf("parseRMC() got2 = %v, want %v", got2, ttt.want2)
			}
			if got3 != ttt.want3 {
				t.Errorf("parseRMC() got3 = %v, want %v", got3, ttt.want3)
			}
		})
	}
}
