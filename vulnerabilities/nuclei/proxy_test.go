package main

import (
	"strings"
	"testing"
)

func TestProxyFromArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr string
	}{
		{name: "no args", args: nil, want: ""},
		{name: "-proxy value", args: []string{"-proxy", "http://192.0.2.1:8080"}, want: "http://192.0.2.1:8080"},
		{name: "-p value", args: []string{"-p", "socks5://192.0.2.1:1080"}, want: "socks5://192.0.2.1:1080"},
		{name: "scheme normalised", args: []string{"-p", "SOCKS5://192.0.2.1:1080"}, want: "socks5://192.0.2.1:1080"},
		{name: "--proxy=value", args: []string{"--proxy=https://proxy.example.com:3128"}, want: "https://proxy.example.com:3128"},
		{name: "last wins", args: []string{"-p", "http://192.0.2.1:1", "-proxy=http://192.0.2.2:2"}, want: "http://192.0.2.2:2"},
		{name: "missing value", args: []string{"-proxy"}, wantErr: "needs a value"},
		{name: "empty value", args: []string{"-proxy="}, wantErr: "invalid -proxy"},
		{name: "no scheme", args: []string{"-proxy", "192.0.2.1:8080"}, wantErr: "invalid -proxy"},
		{name: "bad scheme", args: []string{"-proxy", "ftp://192.0.2.1:21"}, wantErr: "scheme"},
		{name: "comma list", args: []string{"-proxy", "http://192.0.2.1:1,http://192.0.2.2:2"}, wantErr: "invalid -proxy"},
		{name: "other flag", args: []string{"-severity", "high"}, wantErr: "unsupported custom arg"},
		{name: "credentials not echoed", args: []string{"-proxy", "ftp://user:s3cret@192.0.2.1:21"}, wantErr: "scheme"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proxyFromArgs(tt.args)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
				}
				if strings.Contains(err.Error(), "s3cret") {
					t.Fatalf("error leaks credentials: %v", err)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}
