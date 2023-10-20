package gpt

// import (
// 	"testing"
// )
//
// func Test_getLatestCommits(t *testing.T) {
// 	tests := []struct {
// 		url         string
// 		wantCommits int
// 		wantErr     bool
// 	}{
// 		{"github.com/kaatinga/commit", 10, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.url, func(t *testing.T) {
// 			got, err := getLatestCommits(tt.url, 10)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getLatestCommits() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if len(got) != tt.wantCommits {
// 				t.Errorf("getLatestCommits() got = %v, want %v", len(got), tt.wantCommits)
// 			}
//
// 			for _, commit := range got {
// 				t.Logf("commit: %s, from: %s", commit.Message.Content, commit.Date)
// 			}
// 		})
// 	}
// }
