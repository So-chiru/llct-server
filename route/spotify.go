package route

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"crypto/rand"

	"github.com/gorilla/sessions"
)

type SpotifyTokenResponse struct {
	AccessToken      string  `json:"access_token"`
	TokenType        string  `json:"token_type"`
	Scope            string  `json:"scope"`
	ExpiresIn        int     `json:"expires_in"`
	RefreshToken     string  `json:"refresh_token"`
	Error            *string `json:"error"`
	ErrorDescription *string `json:"error_description"`
}

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

var store *sessions.FilesystemStore

func GenerateSessionStore() {
	bytes, err := GenRandomBytes(32)

	if err != nil {
		panic(err)
	}

	store = sessions.NewFilesystemStore("", bytes)
}

func fetchSpotifyToken(code string, redirect_uri *string) (*SpotifyTokenResponse, error) {
	form := url.Values{}
	if redirect_uri == nil {
		form.Add("grant_type", "refresh_token")
		form.Add("refresh_token", code)
	} else {
		form.Add("grant_type", "authorization_code")
		form.Add("code", code)
		form.Add("redirect_uri", *redirect_uri)
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	var SPOTIFY_CLIENT_ID = os.Getenv("SPOTIFY_CLIENT_ID")
	var SPOTIFY_CLIENT_SECRET = os.Getenv("SPOTIFY_CLIENT_SECRET")
	auth_code := b64.StdEncoding.EncodeToString([]byte(SPOTIFY_CLIENT_ID + ":" + SPOTIFY_CLIENT_SECRET))
	req.Header.Add("Authorization", "Basic "+auth_code)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *SpotifyTokenResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("failed to fetch spotify token data")
	}

	if result.ErrorDescription != nil {
		return nil, fmt.Errorf(*result.ErrorDescription)
	}

	if result.Error != nil {
		return nil, fmt.Errorf(*result.Error)
	}

	return result, nil
}

func spotifyLoginHandler(w http.ResponseWriter, r *http.Request) {
	if store == nil {
		GenerateSessionStore()
	}

	w.Header().Add("Cache-Control", "no-cache, private, max-age=0")

	var SPOTIFY_CLIENT_ID = os.Getenv("SPOTIFY_CLIENT_ID")

	// if SPOTIFY_CLIENT_ID == "" || pass != SPOTIFY_CLIENT_ID {
	// 	var error_string = []byte("올바르지 않은 요청입니다.")
	// 	CreateJsonResponse(&w, false, &error_string)

	// 	return
	// }

	session, _ := store.Get(r, "session-name")

	from := r.URL.Query().Get("from")
	if from != "" {
		escape, err := url.QueryUnescape(r.URL.Query().Get("from"))
		if err == nil {
			session.Values["from"] = escape
		}
	}

	if session.Values["spotifyCode"] != nil {
		expireAt := session.Values["spotifyExpireAt"].(int64)
		if session.Values["spotifyExpireAt"] == nil || time.Now().Unix() < expireAt {
			if from != "" {
				http.Redirect(w, r, from+"?spotifyToken="+fmt.Sprintf("%v", session.Values["spotifyCode"]), http.StatusFound)
			}

			var error_string = []byte("이미 로그인된 상태입니다.")
			CreateJsonResponse(&w, false, &error_string)
			return
		}
	}

	var state string
	if session.Values["spotifyState"] == nil {
		bytes, err := GenRandomBytes(16)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		state = fmt.Sprintf("%x", bytes)

		session.Values["spotifyState"] = state

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		state = fmt.Sprintf("%v", session.Values["spotifyState"])
	}

	url := url.URL{
		Scheme: "https",
		Host:   "accounts.spotify.com",
		Path:   "/authorize",
	}

	scheme := r.Header.Get("Scheme")
	if scheme == "" {
		scheme = "http"
	}

	query := url.Query()
	query.Set("response_type", "code")
	query.Set("scope", "streaming user-read-email user-read-private")
	query.Set("client_id", SPOTIFY_CLIENT_ID)
	query.Set("redirect_uri", scheme+"://"+r.Host+"/spotify/callback")
	query.Set("state", state)

	url.RawQuery = query.Encode()

	http.Redirect(w, r, url.String(), http.StatusFound)
}

func spotifyCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if store == nil {
		GenerateSessionStore()
	}

	w.Header().Add("Cache-Control", "no-cache, private, max-age=0")

	if r.URL.Query().Get("code") == "" || r.URL.Query().Get("state") == "" {
		var error_string = []byte("code 또는 state 값이 지정되지 않았습니다.")
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	session, _ := store.Get(r, "session-name")

	if r.URL.Query().Has("error") {
		if session.Values["from"] != nil {
			http.Redirect(w, r, fmt.Sprintf("%v", session.Values["from"])+"?error="+r.URL.Query().Get("error"), http.StatusFound)

			session.Values["from"] = nil
			session.Save(r, w)
		} else {
			var error_string = []byte("오류가 발생했습니다: " + r.URL.Query().Get("error"))
			CreateJsonResponse(&w, false, &error_string)
		}

		return
	}

	state := fmt.Sprintf("%v", session.Values["spotifyState"])
	if state == "<nil>" || state != r.URL.Query().Get("state") {
		var error_string = []byte("잘못된 요청입니다.")
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	if session.Values["spotifyCode"] != nil {
		expireAt := session.Values["spotifyExpireAt"].(int64)
		from := fmt.Sprintf("%v", session.Values["from"])
		if session.Values["spotifyExpireAt"] == nil || time.Now().Unix() < expireAt {
			if from != "" {
				http.Redirect(w, r, from+"?spotifyToken="+fmt.Sprintf("%v", session.Values["spotifyCode"]), http.StatusFound)
			}

			var error_string = []byte("이미 로그인된 상태입니다.")
			CreateJsonResponse(&w, false, &error_string)
			return
		}
	}

	code := r.URL.Query().Get("code")

	scheme := r.Header.Get("Scheme")
	if scheme == "" {
		scheme = "http"
	}

	cb_url := scheme + "://" + r.Host + "/spotify/callback"

	res, err := fetchSpotifyToken(code, &cb_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["spotifyCode"] = res.AccessToken
	session.Values["spotifyRefreshToken"] = res.RefreshToken
	session.Values["spotifyExpireAt"] = time.Now().Unix() + int64(res.ExpiresIn)

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["from"] != nil {
		http.Redirect(w, r, fmt.Sprintf("%v", session.Values["from"])+"?spotifyToken="+res.AccessToken+"&spotifyExpires="+fmt.Sprint(res.ExpiresIn), http.StatusFound)

		session.Values["from"] = nil
		session.Save(r, w)

		return
	}

	var data = []byte("로그인에 성공했습니다. 애플리케이션으로 돌아가세요.")
	CreateJsonResponse(&w, true, &data)
}

func spotifyRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if store == nil {
		GenerateSessionStore()
	}

	w.Header().Add("Cache-Control", "no-cache, private, max-age=0")

	session, _ := store.Get(r, "session-name")

	if session.Values["spotifyState"] == nil || session.Values["spotifyRefreshToken"] == nil {
		var error_string = []byte("잘못된 요청입니다.")
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	refresh_token := fmt.Sprintf("%v", session.Values["spotifyRefreshToken"])
	res, err := fetchSpotifyToken(refresh_token, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["spotifyCode"] = res.AccessToken
	session.Values["spotifyRefreshToken"] = res.RefreshToken
	session.Values["spotifyExpireAt"] = time.Now().Unix() + int64(res.ExpiresIn)

	var data = []byte(res.AccessToken)
	CreateJsonResponse(&w, true, &data)
}
