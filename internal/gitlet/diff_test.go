package gitlet

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"testing"
	"time"
)

func TestNewGitInfo(t *testing.T) {
	// create .git folder and test config
	_, err := RunCommand(`git init`, ".")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = RunCommand(`git config user.name Michael`, ".")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = RunCommand(`git config user.email a@dbc.de`, ".")
	if err != nil {
		t.Error(err)
		return
	}

	t.Cleanup(func() {
		_, err = RunCommand(`rm -rf .git`, ".")
		if err != nil {
			t.Error(err)
		}
	})

	type args struct {
		path string
		msg  string
	}
	tests := []struct {
		name        string
		args        args
		wantGitInfo *GitInfo
		wantErr     bool
	}{
		{name: "test1", args: args{".", "test"},
			wantGitInfo: &GitInfo{
				Msg: "test",
				Signature: object.Signature{
					Name:  "Michael",
					Email: "a@dbc.de",
					When:  time.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGitInfo, err := NewGitInfo(tt.args.path, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGitInfo() error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Logf("NewGitInfo() error = %v", err)
				}
				return
			}

			if gotGitInfo.Msg != tt.wantGitInfo.Msg {
				t.Errorf("NewGitInfo() gotGitInfo.Msg = %v, want %v", gotGitInfo.Msg, tt.wantGitInfo.Msg)
			}

			if gotGitInfo.Name != tt.wantGitInfo.Name {
				t.Errorf("NewGitInfo() gotGitInfo.Name = %v, want %v", gotGitInfo.Name, tt.wantGitInfo.Name)
			}

			if gotGitInfo.Email != tt.wantGitInfo.Email {
				t.Errorf("NewGitInfo() gotGitInfo.Email = %v, want %v", gotGitInfo.Email, tt.wantGitInfo.Email)
			}
		})
	}
}

func TestGetDiff(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test 1", args{"."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDiff(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}