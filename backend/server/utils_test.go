package server

import (
	"reflect"
	"testing"
)

func TestEnumerateArticles(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "Check regular timestamp",
			args: args{
				path: "../../rsrc/test_articles/1",
			},
			want:    1713382077,
			wantErr: false,
		},
		{
			name: "Check minus timestamp",
			args: args{
				path: "../../rsrc/test_articles/2",
			},
			wantErr: true,
		},
		{
			name: "After 32bits are not enough...",
			args: args{
				path: "../../rsrc/test_articles/3",
			},
			want:    3426764154,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnumerateArticles(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnumerateArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil == tt.wantErr {
				return
			}

			if got[0].CreationTimestamp != tt.want {
				t.Errorf("EnumerateArticles() = %v, want %v", got[0].CreationTimestamp, tt.want)
				return
			}
		})
	}
}

func TestEnumerateArticles2(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Check optional field",
			args: args{
				path: "../../rsrc/test_articles/1",
			},
			want:    "In-depth analysis of Authenticode policy on Windows operating system. Hands-on example for creating self-signed certificate",
			wantErr: false,
		},
		{
			name: "Check missing optional field",
			args: args{
				path: "../../rsrc/test_articles/3",
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnumerateArticles(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnumerateArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == true {
				return
			}
			new := got[0].MetaDescription
			if reflect.DeepEqual(new, tt.want) {
				t.Errorf("EnumerateArticles() = %v, want %v", got[0].MetaDescription, tt.want)
				return
			}
		})
	}
}
