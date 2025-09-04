package bot

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

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

func (sc *SteamBot) ReceiveFromUser(assetID, tradeURL, SellerID string) (string, error) {
	log.Info().Str("assetID", assetID).Str("tradeURL", tradeURL).Msg("RECEIVE FROM START")

	partner, token, err := parseTradeURL(tradeURL)
	if err != nil {
		log.Error().Err(err).Str("tradeURL", tradeURL).Msg("Failed to parse trade URL")
		return "", err
	}
	log.Info().Str("partner", partner).Str("token", token).Msg("Parsed trade URL successfully")

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
		"partner":                   {SellerID},
		"tradeoffermessage":         {""},
		"trade_offer_create_params": {fmt.Sprintf(`{"trade_offer_access_token":"%s"}`, token)},
		"json_tradeoffer":           {toJSON(offer)},
	}

	req, err := http.NewRequest("POST", "https://steamcommunity.com/tradeoffer/new/send", strings.NewReader(form.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("err to create HTTP request for trade offer")
		return "", err
	}

	req.Header.Set("Referer", tradeURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := sc.Client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("err to send trade offer request")
		return "", err
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		log.Info().Interface("response_headers", resp.Header).
			Msg("Received response headers")
		if gzReader, err := gzip.NewReader(resp.Body); err == nil {
			defer gzReader.Close()
			reader = gzReader
		}
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error ReceiveFromUser %d", resp.StatusCode)
	}

	var res struct {
		TradeOfferID string `json:"tradeofferid"`
	}
	bodyBytes, _ := io.ReadAll(reader)
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return "", fmt.Errorf("failed to parse tradeofferid: %w", err)
	}
	log.Info().Str("id", string(bodyBytes)).
		Msg("resp tradeofferid new trade")

	return res.TradeOfferID, nil
}

func (sc *SteamBot) SendToBuyer(assetID, tradeURL, buyerID string) error {
	_, token, err := parseTradeURL(tradeURL)
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
		"partner":                   {buyerID},
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

// func (sc *SteamBot) GetStatus(tradeOfferID string) error {
// 	url := fmt.Sprintf("https://steamcommunity.com/tradeoffer/%s/?json=1", tradeOfferID)

// 	resp, err := sc.apiCall("GET", url, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to get trade status %s: %w", tradeOfferID, err)
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to get trade status %s: unexpected status code %d", tradeOfferID, resp.StatusCode)
// 	}

// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to read body: %w", err)
// 	}

//		log.Info().Msgf("TradeOffer %s status body: %s", tradeOfferID, string(body))
//		return nil
//	}

func (sc *SteamBot) GetStatus(tradeOfferID string) error {
	return nil
}

func (sc *SteamBot) DeclineTrade(tradeOfferID string) error {
	params := map[string]string{
		"sessionid": sc.GetSessionID(),
	}

	resp, err := sc.apiCall("POST",
		fmt.Sprintf("https://steamcommunity.com/tradeoffer/%s/cancel", tradeOfferID), params)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to cancel trade %s, status: %d, body: %s", tradeOfferID, resp.StatusCode, body)
	}
	log.Info().Any("BODY:", resp.Body).Msg("resp:")

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
