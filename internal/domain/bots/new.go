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

type RSAParams struct {
	PublicKey *rsa.PublicKey
	Timestamp string
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

func (sc *SteamBot) apiCall(method, endpoint string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, SteamAPIURL+endpoint, nil)
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

func (sc *SteamBot) SignTradeOffer(tradeURL string) string {
	secret := strings.TrimSpace(sc.IdentitySecret)
	if len(secret)%4 != 0 {
		secret += strings.Repeat("=", 4-len(secret)%4)
	}

	key, _ := base64.StdEncoding.DecodeString(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(tradeURL))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// func (sc *SteamBot) GetParamsForTrade(tradeURL string, itemID string, message string, apiKey string) (string, url.Values) {
// 	signature := sc.SignTradeOffer(tradeURL)

// 	apiURL := "https://api.steampowered.com/IEconService/CreateTradeOffer/v1/"
// 	params := url.Values{}
// 	params.Add("key", apiKey)
// 	params.Add("tradeoffer_message", message)
// 	params.Add("trade_offer_access_token", sc.AccessToken)
// 	params.Add("trade_url", tradeURL)
// 	params.Add("items_to_give[]", itemID)
// 	params.Add("signature", signature)

// 	return apiURL, params
// }

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
	resp, err := sc.apiCall("GET", "/ITwoFactorService/QueryTime/v1/", nil)
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
			ServerTime interface{} `json:"server_time"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return time.Time{}, fmt.Errorf("failed to unmarshal server time: %w", err)
	}

	var serverTime int64

	switch v := response.Response.ServerTime.(type) {
	case float64:
		serverTime = int64(v)
	case string:
		serverTime, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid server time format: %w", err)
		}
	default:
		return time.Time{}, fmt.Errorf("unexpected server time type: %T", v)
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

	return nil
}

// >>>>>>>>??
// if sc.RefreshToken != "" {
// 	finalizeResp, err := sc.finalizeLogin()
// 	if err == nil && finalizeResp != nil {
// sc.performRedirects(finalizeResp)
// 	}

// sc.syncCookies()
// }

// }
// >>>>>>>>??

// func (sc *SteamBot) RefreshTokens() error {
// 	data := url.Values{}
// 	data.Set("steamid", sc.SteamID)
// 	data.Set("refresh_token", sc.RefreshToken)

// 	req, err := http.NewRequest(
// 		"POST",
// 		SteamAPIURL+"/IAuthenticationService/GenerateAccessTokenForApp/v1/",
// 		strings.NewReader(data.Encode()),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	resp, err := sc.Client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	var result struct {
// 		Response struct {
// 			AccessToken string `json:"access_token"`
// 		} `json:"response"`
// 	}

// 	json.Unmarshal(body, &result)

// 	if result.Response.AccessToken == "" {
// 		return fmt.Errorf("failed to refresh token: %s", string(body))
// 	}

// 	slog.Info("Token", "is", result.Response.AccessToken)

// 	sc.AccessToken = result.Response.AccessToken
// 	return nil
// }

// func (sc *SteamBot) syncCookies() {
// 	u, _ := url.Parse(SteamCommunityURL)
// 	cookies := sc.Client.Jar.Cookies(u)

// 	domains := []string{
// 		"steamcommunity.com",
// 		"store.steampowered.com",
// 		"help.steampowered.com",
// 		"login.steampowered.com",
// 	}

// 	for _, domain := range domains {
// 		domainURL := &url.URL{Scheme: "https", Host: domain}
// 		currentCookies := sc.Client.Jar.Cookies(domainURL)
// 		cookieMap := make(map[string]*http.Cookie)
// 		for _, cookie := range currentCookies {
// 			cookieMap[cookie.Name] = cookie
// 		}

// 		for _, cookie := range cookies {
// 			if cookie.Name == "sessionid" || cookie.Name == "steamLoginSecure" ||
// 				cookie.Name == "steamRefresh_steam" || cookie.Name == "steamCountry" {
// 				newCookie := &http.Cookie{
// 					Name:   cookie.Name,
// 					Value:  cookie.Value,
// 					Path:   "/",
// 					Domain: domain,
// 				}
// 				cookieMap[newCookie.Name] = newCookie
// 			}
// 		}

// 		newCookies := make([]*http.Cookie, 0, len(cookieMap))
// 		for _, cookie := range cookieMap {
// 			newCookies = append(newCookies, cookie)
// 		}

// 		sc.Client.Jar.SetCookies(domainURL, newCookies)
// 	}
// }
