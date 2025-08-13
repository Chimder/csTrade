package bots

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	SteamAPIURL       = "https://api.steampowered.com"
	SteamCommunityURL = "https://steamcommunity.com"
	SteamStoreURL     = "https://store.steampowered.com"
	SteamLoginURL     = "https://login.steampowered.com"
)

type SteamBot struct {
	Username       string
	Password       string
	SteamID        string
	SharedSecret   string
	IdentitySecret string
	DeviceID       string
	AccessToken    string
	RefreshToken   string
	Client         *http.Client
}

func NewSteamClient(username, password, steamID, sharedSecret, identitySecret, deviceID string) *SteamBot {
	jar, _ := cookiejar.New(nil)

	u, _ := url.Parse(SteamCommunityURL)
	jar.SetCookies(u, []*http.Cookie{
		{Name: "Steam_Language", Value: "english"},
		{Name: "timezoneOffset", Value: "0,0"},
	})
	if deviceID != "" && !strings.HasPrefix(deviceID, "android:") {
		deviceID = "android:" + deviceID
	}

	return &SteamBot{
		Username:       username,
		Password:       password,
		SteamID:        steamID,
		SharedSecret:   sharedSecret,
		IdentitySecret: identitySecret,
		DeviceID:       deviceID,
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
				req.Header.Set("Origin", SteamCommunityURL)
				return nil
			},
		},
	}
}

type Asset struct {
	AppID      int    `json:"appid"`
	ContextID  string `json:"contextid"`
	AssetID    string `json:"assetid"`
	Amount     int    `json:"amount"`
	ClassID    string `json:"classid,omitempty"`
	InstanceID string `json:"instanceid,omitempty"`
}

func (sc *SteamBot) GetSessionID() string {
	u, _ := url.Parse(SteamCommunityURL)
	if sc.Client == nil || sc.Client.Jar == nil {
		return ""
	}
	for _, c := range sc.Client.Jar.Cookies(u) {
		if c.Name == "sessionid" {
			return c.Value
		}
	}
	return ""
}
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func parseTradeURL(tradeURL string) (partnerID, token string, err error) {
	u, err := url.Parse(tradeURL)
	if err != nil {
		return "", "", err
	}
	q := u.Query()
	partnerID = q.Get("partner")
	token = q.Get("token")
	if partnerID == "" || token == "" {
		return "", "", fmt.Errorf("invalid tradeURL")
	}
	return partnerID, token, nil
}

func (sc *SteamBot) ReceiveFromUser(assetID, tradeURL string) error {
	partner, token, err := parseTradeURL(tradeURL)
	if err != nil {
		return err
	}
	offer := map[string]interface{}{
		"newversion": true,
		"version":    2,
		"me":         map[string]interface{}{"assets": []map[string]string{}},
		"them": map[string]interface{}{"assets": []map[string]string{
			{"appid": "730", "contextid": "2", "assetid": assetID},
		}},
	}
	form := url.Values{
		"sessionid":                 {sc.GetSessionID()},
		"serverid":                  {"1"},
		"partner":                   {partner},
		"tradeoffermessage":         {""},
		"trade_offer_create_params": {fmt.Sprintf(`{"trade_offer_access_token":"%s"}`, token)},
		"json_tradeoffer":           {toJSON(offer)},
	}
	req, _ := http.NewRequest("POST", "https://steamcommunity.com/tradeoffer/new/send", strings.NewReader(form.Encode()))
	req.Header.Set("Referer", tradeURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := sc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Err code offer FromUser: %d", resp.StatusCode)
	}
	return nil
}

func (sc *SteamBot) SendToBuyer(assetID, tradeURL string) error {
	partner, token, err := parseTradeURL(tradeURL)
	if err != nil {
		return err
	}
	offer := map[string]interface{}{
		"newversion": true,
		"version":    2,
		"me": map[string]interface{}{"assets": []map[string]string{
			{"appid": "730", "contextid": "2", "assetid": assetID},
		}},
		"them": map[string]interface{}{"assets": []map[string]string{}},
	}
	form := url.Values{
		"sessionid":                 {sc.GetSessionID()},
		"serverid":                  {"1"},
		"partner":                   {partner},
		"tradeoffermessage":         {""},
		"trade_offer_create_params": {fmt.Sprintf(`{"trade_offer_access_token":"%s"}`, token)},
		"json_tradeoffer":           {toJSON(offer)},
	}
	req, _ := http.NewRequest("POST", "https://steamcommunity.com/tradeoffer/new/send", strings.NewReader(form.Encode()))
	req.Header.Set("Referer", tradeURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var res struct {
		TradeOfferID            string `json:"tradeofferid"`
		StrError                string `json:"strError"`
		NeedsMobileConfirmation bool   `json:"needs_mobile_confirmation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("err parse tradeOffer to buyer: %w", err)
	}
	log.Info().Interface("resp", res).Msg("To buyer")

	if res.StrError != "" {
		return fmt.Errorf("steam error: %s", res.StrError)
	}

	return nil
}

func (sc *SteamBot) DeclineTrade(tradeOfferID, apiKey string) error {
	params := map[string]string{
		"key":          apiKey,
		"tradeofferid": tradeOfferID,
	}
	resp, err := sc.apiCall("POST", "/IEconService/DeclineTradeOffer/v1/", params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Err trade statusCode: %d", resp.StatusCode)
	}
	return nil
}

func (sc *SteamBot) apiCall(method, endpoint string, params map[string]string) (*http.Response, error) {
	urlStr := endpoint
	if !strings.HasPrefix(endpoint, "http") {
		urlStr = SteamAPIURL + endpoint
	}
	req, err := http.NewRequest(method, urlStr, nil)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
		"Origin":     []string{SteamCommunityURL},
	}

	if method == "GET" && params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	} else if method == "POST" && params != nil {
		data := url.Values{}
		for k, v := range params {
			data.Set(k, v)
		}
		req.Body = io.NopCloser(strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return sc.Client.Do(req)
}

func (sc *SteamBot) GenerateTOTPCode() (string, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(sc.SharedSecret))
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %v", err)
	}

	steamTime, err := sc.GetSteamTime()
	if err != nil {
		return "", fmt.Errorf("failed to get Steam time: %v", err)
	}

	timestamp := steamTime.Unix() / 30

	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timestamp))

	h := hmac.New(sha1.New, key)
	h.Write(timeBytes)
	hash := h.Sum(nil)

	offset := hash[19] & 0x0F
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

	chars := "23456789BCDFGHJKMNPQRTVWXY"
	result := make([]byte, 5)
	for i := 0; i < 5; i++ {
		result[i] = chars[code%uint32(len(chars))]
		code /= uint32(len(chars))
	}

	return string(result), nil
}

type RSAParams struct {
	PublicKey *rsa.PublicKey
	Timestamp string
}

func (sc *SteamBot) fetchRSAParams() (*RSAParams, error) {
	params := map[string]string{"account_name": sc.Username}
	resp, err := sc.apiCall("GET", "/IAuthenticationService/GetPasswordRSAPublicKey/v1/", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var rsaResp struct {
		Response struct {
			PublicKeyMod string `json:"publickey_mod"`
			PublicKeyExp string `json:"publickey_exp"`
			Timestamp    string `json:"timestamp"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rsaResp); err != nil {
		return nil, err
	}

	mod := new(big.Int)
	mod.SetString(rsaResp.Response.PublicKeyMod, 16)

	exp := new(big.Int)
	exp.SetString(rsaResp.Response.PublicKeyExp, 16)

	return &RSAParams{
		PublicKey: &rsa.PublicKey{N: mod, E: int(exp.Int64())},
		Timestamp: rsaResp.Response.Timestamp,
	}, nil
}

type LoginResponse struct {
	Response struct {
		ClientID             string `json:"client_id"`
		RequestID            string `json:"request_id"`
		Interval             int    `json:"interval"`
		AllowedConfirmations []struct {
			ConfirmationType int `json:"confirmation_type"`
		} `json:"allowed_confirmations"`
		Steamid              string `json:"steamid"`
		WeakToken            string `json:"weak_token"`
		ExtendedErrorMessage string `json:"extended_error_message"`
	} `json:"response"`
}

func (sc *SteamBot) startAuth() (*LoginResponse, error) {
	rsaParams, err := sc.fetchRSAParams()
	if err != nil {
		return nil, err
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaParams.PublicKey, []byte(sc.Password))
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"account_name":         sc.Username,
		"encrypted_password":   base64.StdEncoding.EncodeToString(encrypted),
		"encryption_timestamp": rsaParams.Timestamp,
		"persistence":          "1",
	}

	resp, err := sc.apiCall("POST", "/IAuthenticationService/BeginAuthSessionViaCredentials/v1/", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, err
	}

	return &loginResp, nil
}

func (sc *SteamBot) submitTOTP(clientID string) error {
	code, err := sc.GenerateTOTPCode()
	if err != nil {
		return err
	}
	slog.Info("Generated TOTP", "code", code)

	data := map[string]string{
		"client_id": clientID,
		"steamid":   sc.SteamID,
		"code":      code,
		"code_type": "3",
	}

	resp, err := sc.apiCall("POST", "/IAuthenticationService/UpdateAuthSessionWithSteamGuardCode/v1/", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	slog.Debug("TOTP submit response", "status", resp.Status, "body", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TOTP failed status: %d", resp.StatusCode)
	}
	return nil
}

func (sc *SteamBot) getTokens(clientID, requestID string) error {
	data := map[string]string{
		"client_id":  clientID,
		"request_id": requestID,
	}

	resp, err := sc.apiCall("POST", "/IAuthenticationService/PollAuthSessionStatus/v1/", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Response.AccessToken == "" {
		return fmt.Errorf("empty access token")
	}

	sc.AccessToken = result.Response.AccessToken
	sc.RefreshToken = result.Response.RefreshToken
	return nil
}

func (sc *SteamBot) GetSteamTime() (time.Time, error) {
	resp, err := sc.apiCall("POST", "/ITwoFactorService/QueryTime/v1/", nil)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}

	var response struct {
		Response struct {
			ServerTime string `json:"server_time"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return time.Time{}, fmt.Errorf("failed to unmarshal server time: %w", err)
	}

	serverTime, err := strconv.ParseInt(response.Response.ServerTime, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid server time format: %w", err)
	}

	return time.Unix(serverTime, 0), nil
}

func (sc *SteamBot) Login() error {
	loginResp, err := sc.startAuth()
	if err != nil {
		return fmt.Errorf("start auth failed: %w", err)
	}

	if err := sc.submitTOTP(loginResp.Response.ClientID); err != nil {
		return fmt.Errorf("TOTP submit failed: %w", err)
	}

	if err := sc.getTokens(loginResp.Response.ClientID, loginResp.Response.RequestID); err != nil {
		return fmt.Errorf("get tokens failed: %w", err)
	}

	sc.syncCookies()

	return sc.testSession()
}

func (sc *SteamBot) syncCookies() {
	req, _ := http.NewRequest("GET", "https://steamcommunity.com/login/home/", nil)
	resp, err := sc.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	mainCookies := sc.Client.Jar.Cookies(&url.URL{Scheme: "https", Host: "steamcommunity.com"})

	domains := []string{
		"store.steampowered.com",
		"help.steampowered.com",
		"login.steampowered.com",
	}

	for _, domain := range domains {
		domainURL := &url.URL{Scheme: "https", Host: domain}
		copied := []*http.Cookie{}

		for _, c := range mainCookies {
			if c.Name == "sessionid" || c.Name == "steamLoginSecure" || c.Name == "steamCountry" {
				newC := *c
				newC.Domain = domain
				newC.Path = "/"
				newC.Secure = true
				copied = append(copied, &newC)
			}
		}

		if len(copied) > 0 {
			sc.Client.Jar.SetCookies(domainURL, copied)
		}
	}
}

func (sc *SteamBot) testSession() error {
	req, _ := http.NewRequest("GET", "https://store.steampowered.com/account", nil)
	resp, err := sc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("TestSession Status: %s\n", resp.Status)

	return nil
}
