package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/kimitzu/kimitzu-go/schema"
	"github.com/ipfs/go-ipfs/repo/fsrepo"

	"os"

	"github.com/kimitzu/kimitzu-go/core"
)

func ChangeAPICredentials(repoPath, username, password string, authenticated bool) error {
	// Set repo path
	cfgPath := path.Join(repoPath, "config")
	configFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	_, err = fsrepo.Open(repoPath)
	if _, ok := err.(fsrepo.NoRepoError); ok {
		return fmt.Errorf(
			"IPFS repo in the data directory '%s' has not been initialized."+
				"\nRun openbazaar with the 'start' command to initialize.",
			repoPath)
	}
	if err != nil {
		return err
	}

	configJson := make(map[string]interface{})
	err = json.Unmarshal(configFile, &configJson)
	if err != nil {
		return err
	}

	apiCfg, err := schema.GetAPIConfig(configFile)
	if err != nil {
		log.Error(err)
		return err
	}

	pw := password
	if strings.Contains(username, "\r\n") {
		apiCfg.Username = strings.Replace(username, "\r\n", "", -1)
	} else if strings.Contains(username, "\n") {
		apiCfg.Username = strings.Replace(username, "\n", "", -1)
	} else {
		apiCfg.Username = username
	}
	apiCfg.Authenticated = authenticated
	h := sha256.Sum256([]byte(pw))
	apiCfg.Password = hex.EncodeToString(h[:])
	if len(apiCfg.AllowedIPs) == 0 {
		apiCfg.AllowedIPs = []string{}
	}

	configJson["JSON-API"] = apiCfg

	out, err := json.MarshalIndent(configJson, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cfgPath, out, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println("Changed API credentials")

	return nil

}

// KimitzuInfo - returns the essential information such as repository path
func KimitzuInfo(node *core.OpenBazaarNode, cookie http.Cookie, config schema.APIConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			http.Error(w, fmt.Sprintf(`{"error": "Unauthorized", "origin": "%v"}`, r.RemoteAddr), 401)
			return
		}
		fmt.Fprintf(w,
			`{"repoPath": "%v", "cookie": "%v", "username": "%v", "password": "%v", "authenticated": %v}`,
			strings.ReplaceAll(node.RepoPath, "\\", "\\\\"),
			cookie.String(),
			config.Username,
			config.Password,
			config.Authenticated)
		return
	}

}

type Config struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Authenticated bool   `json:"authenticated"`
}

// KimitzuConfig - changes authentication // POTENTIALLY UNSAFE, IT'S A HACK FOR THE TIME BEING.
func KimitzuConfig(node *core.OpenBazaarNode, cookie http.Cookie, config schema.APIConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			http.Error(w, fmt.Sprintf(`{"error": "Unauthorized", "origin": "%v"}`, r.RemoteAddr), 401)
			return
		}

		if r.Method != "POST" {
			http.Error(w, `{"error": "MethodNotSupported"}`, 405)
			return
		}

		newConfig := Config{}
		b, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		json.Unmarshal(b, &newConfig)
		fmt.Println("New Credentials", newConfig)

		err := ChangeAPICredentials(node.RepoPath, newConfig.Username, newConfig.Password, newConfig.Authenticated)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed Changing credentials", "details": "%v"}`, err), 500)
			return
		}

		fmt.Fprintf(w,
			`{"repoPath": "%v", "cookie": "%v", "username": "%v", "authenticated": %v}`,
			strings.ReplaceAll(node.RepoPath, "\\", "\\\\"),
			cookie.String(),
			newConfig.Username,
			config.Authenticated)
		return
	}

}
