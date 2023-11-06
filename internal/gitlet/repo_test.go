package gitlet

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isDirInHomeDir(t *testing.T) {
	tests := []struct {
		pathItems  []string
		gitGSItems []string
		want       bool
	}{
		{[]string{"Users", "username", "repos"}, []string{"Users", "username", "repos"}, false},
		{[]string{"Users", "username", "repos"}, []string{"Users", "username", "repos", "test"}, false},
		{[]string{"Users", "username", "repos", "test"}, []string{"Users", "username", "repos"}, true},
		{[]string{"Users", "username", "confetos"}, []string{"Users", "username", "repos"}, false},
	}
	for _, tt := range tests {
		t.Run(filepath.Join(tt.pathItems...)+" in "+filepath.Join(tt.gitGSItems...), func(t *testing.T) {
			assert.Equalf(t, tt.want, isDirInHomeDir(tt.pathItems, tt.gitGSItems), "isValidGSDir(%v, %v)", tt.pathItems, tt.gitGSItems)
		})
	}
}
