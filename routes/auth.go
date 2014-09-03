package routes

import (
	"encoding/json"
	"fmt"
	"github.com/callumj/auto-convert/shared"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type authResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Uid         int64  `json:"uid,string"`
}

func BeginAuthHandler(c http.ResponseWriter, req *http.Request) {
	redirectTo, err := url.Parse("../complete_auth")
	if err != nil {
		log.Printf("Unable to build relative to complete_auth: %v\n", err)
		return
	}

	finalUrl := shared.GetFullUrl(req).ResolveReference(redirectTo)
	dropboxRedirect := fmt.Sprintf("https://www.dropbox.com/1/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s", shared.Config.DropboxKey, finalUrl.String(), "NOT_IMPLEMENTED")
	http.Redirect(c, req, dropboxRedirect, http.StatusFound)
}

func CompleteAuthHandler(c http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if state != "NOT_IMPLEMENTED" {
		return
	}
	code := req.FormValue("code")

	body := "<h2>Done :)</h2>"
	redirectTo, err := url.Parse("../complete_auth")
	if err != nil {
		log.Printf("Unable to build relative to complete_auth: %v\n", err)
		return
	}

	finalUrl := shared.GetFullUrl(req).ResolveReference(redirectTo)

	go processAuthCode(code, finalUrl.String())

	c.Header().Add("Content-Type", "text/html")
	c.Header().Add("Content-Length", strconv.Itoa(len(body)))
	io.WriteString(c, body)
}

func processAuthCode(code, redirectFrom string) {
	resp, err := http.PostForm("https://api.dropbox.com/1/oauth2/token",
		url.Values{"code": {code},
			"grant_type":    {"authorization_code"},
			"client_id":     {shared.Config.DropboxKey},
			"client_secret": {shared.Config.DropboxSecret},
			"redirect_uri":  {redirectFrom},
		})

	defer resp.Body.Close()
	if err == nil {
		dec := json.NewDecoder(resp.Body)
		var decoded authResp
		err = dec.Decode(&decoded)
		if err != nil {
			log.Print(err)
			return
		}

		if decoded.TokenType != "bearer" {
			log.Printf("Expected bearer, got %v\n", decoded.TokenType)
		} else {
			acc := shared.Account{
				Uid:   decoded.Uid,
				Token: decoded.AccessToken,
			}
			shared.AddUpdateAccount(acc)
		}
	} else {
		log.Printf("Error: %v\n", err)
	}
}
