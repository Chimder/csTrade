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
	"os"
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

type SteamClient struct {
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

func NewSteamClient(username, password, steamID, sharedSecret, identitySecret, deviceID string) *SteamClient {
	jar, _ := cookiejar.New(nil)

	u, _ := url.Parse(SteamCommunityURL)
	jar.SetCookies(u, []*http.Cookie{
		{Name: "Steam_Language", Value: "english"},
		{Name: "timezoneOffset", Value: "0,0"},
	})
	if deviceID != "" && !strings.HasPrefix(deviceID, "android:") {
		deviceID = "android:" + deviceID
	}

	return &SteamClient{
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

func (sc *SteamClient) apiCall(method, service, endpoint, version string, params map[string]string) (*http.Response, error) {
	apiURL := fmt.Sprintf("%s/%s/%s/%s", SteamAPIURL, service, endpoint, version)

	req, err := http.NewRequest(method, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", SteamCommunityURL+"/")
	req.Header.Set("Origin", SteamCommunityURL)
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

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

func (sc *SteamClient) GenerateTOTPCode() (string, error) {
	secret := strings.TrimSpace(sc.SharedSecret)
	if len(secret)%4 != 0 {
		secret += strings.Repeat("=", 4-len(secret)%4)
	}

	key, err := base64.StdEncoding.DecodeString(secret)
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

func (sc *SteamClient) SignTradeOffer(tradeURL string) string {
	secret := strings.TrimSpace(sc.IdentitySecret)
	if len(secret)%4 != 0 {
		secret += strings.Repeat("=", 4-len(secret)%4)
	}

	key, _ := base64.StdEncoding.DecodeString(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(tradeURL))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (sc *SteamClient) GetParamsForTrade(tradeURL string, itemID string, message string, apiKey string) (string, url.Values) {
	signature := sc.SignTradeOffer(tradeURL)

	apiURL := "https://api.steampowered.com/IEconService/CreateTradeOffer/v1/"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("tradeoffer_message", message)
	params.Add("trade_offer_access_token", sc.AccessToken)
	params.Add("trade_url", tradeURL)
	params.Add("items_to_give[]", itemID)
	params.Add("signature", signature)

	return apiURL, params
}

func (sc *SteamClient) fetchRSAParams() (*RSAParams, error) {
	if len(sc.Client.Jar.Cookies(&url.URL{Scheme: "https", Host: "steamcommunity.com"})) == 0 {
		_, err := sc.Client.Get(SteamCommunityURL)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize cookies: %v", err)
		}
	}

	params := map[string]string{
		"account_name": sc.Username,
	}

	// for range 5 {
	resp, err := sc.apiCall("GET", "IAuthenticationService", "GetPasswordRSAPublicKey", "v1", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not obtain RSA")
		// continue
	}

	type RSAResponse struct {
		Response struct {
			PublicKeyMod string `json:"publickey_mod"`
			PublicKeyExp string `json:"publickey_exp"`
			Timestamp    string `json:"timestamp"`
		} `json:"response"`
	}
	var rsaResp RSAResponse
	if err := json.NewDecoder(resp.Body).Decode(&rsaResp); err != nil {
		// continue
		return nil, fmt.Errorf("could not obtain RSA")
	}

	if rsaResp.Response.PublicKeyMod == "" || rsaResp.Response.PublicKeyExp == "" {
		// continue
		return nil, fmt.Errorf("could not obtain RSA")
	}

	mod := new(big.Int)
	exp := new(big.Int)

	mod.SetString(rsaResp.Response.PublicKeyMod, 16)
	exp.SetString(rsaResp.Response.PublicKeyExp, 16)

	publicKey := &rsa.PublicKey{
		N: mod,
		E: int(exp.Int64()),
	}

	return &RSAParams{
		PublicKey: publicKey,
		Timestamp: rsaResp.Response.Timestamp,
	}, nil
	// }

	// return nil, fmt.Errorf("could not obtain RSA key after 5 attempts")
}

func (sc *SteamClient) encryptPassword(rsaParams *RSAParams) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaParams.PublicKey, []byte(sc.Password))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (sc *SteamClient) startAuth() (*LoginResponse, error) {
	rsaParams, err := sc.fetchRSAParams()
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := sc.encryptPassword(rsaParams)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("account_name", sc.Username)
	data.Set("persistence", "1")
	data.Set("client_id", "")
	data.Set("encrypted_password", encryptedPassword)
	data.Set("encryption_timestamp", rsaParams.Timestamp)

	req, err := http.NewRequest(
		"POST",
		SteamAPIURL+"/IAuthenticationService/BeginAuthSessionViaCredentials/v1/",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", SteamCommunityURL+"/")
	req.Header.Set("Origin", SteamCommunityURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sc.Client.Do(req)
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

func (sc *SteamClient) submitTOTP(clientID string) error {
	code, err := sc.GenerateTOTPCode()
	if err != nil {
		return err
	}
	fmt.Printf("Generated TOTP: %s\n", code)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("steamid", sc.SteamID)
	data.Set("code", code)
	data.Set("code_type", "3")

	// if sc.DeviceID != "" {
	// 	data.Set("device_details", fmt.Sprintf(`{"device_id":"%s"}`, sc.DeviceID))
	// }

	req, err := http.NewRequest("POST",
		SteamAPIURL+"/IAuthenticationService/UpdateAuthSessionWithSteamGuardCode/v1/",
		strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Origin", SteamCommunityURL)
	req.Header.Set("Referer", SteamCommunityURL+"/mobilelogin?oauth_client_id=DE45CD61")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("SubmitTOTP response:", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK status: %d", resp.StatusCode)
	}

	return nil
}

func (sc *SteamClient) getTokens(clientID, requestID string) error {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("request_id", requestID)

	for range 2 {
		req, err := http.NewRequest(
			"POST",
			SteamAPIURL+"/IAuthenticationService/PollAuthSessionStatus/v1/",
			strings.NewReader(data.Encode()),
		)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		req.Header.Set("Referer", SteamCommunityURL+"/")
		req.Header.Set("Origin", SteamCommunityURL)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := sc.Client.Do(req)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		defer resp.Body.Close()
		var pollResp struct {
			Response struct {
				AccessToken          string `json:"access_token"`
				RefreshToken         string `json:"refresh_token"`
				AccountName          string `json:"account_name"`
				HadRemoteInteraction bool   `json:"had_remote_interaction"`
			} `json:"response"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&pollResp); err != nil {
			return nil
		}

		fmt.Println("Access Token:", pollResp.Response.AccessToken)

		if pollResp.Response.AccessToken != "" {
			sc.AccessToken = pollResp.Response.AccessToken
			sc.RefreshToken = pollResp.Response.RefreshToken
			fmt.Printf("Success! Got tokens\n")
			return nil
		}

		if pollResp.Response.AccessToken == "" {
			return fmt.Errorf("steam error fetch access_token")
		}
	}
	return nil
}

func (sc *SteamClient) GetSteamTime() (time.Time, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ITwoFactorService/QueryTime/v1/?t=%d", time.Now().UnixNano())

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return time.Time{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")

	resp, err := sc.Client.Do(req)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("invalid status code steam time: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		body, _ := io.ReadAll(resp.Body)
		os.WriteFile("steam_time_error.html", body, 0644)
		return time.Time{}, fmt.Errorf("invalid content type: %s", contentType)
	}

	var response struct {
		Response struct {
			ServerTime interface{} `json:"server_time"`
		} `json:"response"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return time.Time{}, fmt.Errorf("json unmarshal error: %v\nBody: %s", err, string(body))
	}

	var serverTime int64

	switch v := response.Response.ServerTime.(type) {
	case float64:
		serverTime = int64(v)
	case string:
		t, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse server_time string: %v", err)
		}
		serverTime = t
	case nil:
		return time.Time{}, fmt.Errorf("server_time is missing in response")
	default:
		return time.Time{}, fmt.Errorf("unexpected type for server_time: %T", v)
	}

	if serverTime == 0 {
		return time.Time{}, fmt.Errorf("invalid server time value: 0")
	}
	slog.Info("time", "is", time.Unix(serverTime, 0))

	return time.Unix(serverTime, 0), nil
}

func (sc *SteamClient) Login() error {
	fmt.Printf("Starting authentication")

	loginResp, err := sc.startAuth()
	if err != nil {
		return fmt.Errorf("start auth failed: %v", err)
	}

	if loginResp.Response.ClientID == "" {
		return fmt.Errorf("login failed: no client_id received")
	}
	fmt.Printf("Got ClientID: %s, RequestID: %s\n", loginResp.Response.ClientID, loginResp.Response.RequestID)

	if err := sc.submitTOTP(loginResp.Response.ClientID); err != nil {
		return fmt.Errorf("TOTP submit failed: %v", err)
	}

	if err := sc.getTokens(loginResp.Response.ClientID, loginResp.Response.RequestID); err != nil {
		return fmt.Errorf("get tokens failed: %v", err)
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

// func (sc *SteamClient) RefreshTokens() error {
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

// func (sc *SteamClient) syncCookies() {
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
