package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/djali-foundation/djali-go/core"
	"github.com/djali-foundation/djali-go/schema"
)

// DjaliInfo - returns the essential information such as repository path
func DjaliInfo(node *core.OpenBazaarNode, cookie http.Cookie, config schema.APIConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			http.Error(w, fmt.Sprintf(`{"error": "Unauthorized", "origin": "%v"}`, r.RemoteAddr), 401)
			return
		}
		fmt.Fprintf(w,
			`{"repoPath": "%v", "cookie": "%v", "username": "%v", "password": "%v"}`,
			strings.ReplaceAll(node.RepoPath, "\\", "\\\\"),
			cookie.String(),
			config.Username,
			config.Password)
		return
	}

}
