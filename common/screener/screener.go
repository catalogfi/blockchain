package screener

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/catalogfi/blockchain"
)

type request struct {
	Chain   string `json:"chain"`
	Address string `json:"address"`
}

type response struct {
	Address       string `json:"address"`
	Chain         string `json:"chain"`
	IsBlacklisted bool   `json:"is_blacklisted"`
}

// Screener sends a list of addresses to the screener server, it returns the blacklisted address.
type Screener interface {
	IsBlacklisted(ctx context.Context, addrs map[string]blockchain.Chain) (map[string]bool, error)
}

// screener implements the Screener interface by querying our screener server. It will return false for all addresses if
// the url is empty.
type screener struct {
	url string
}

func New(url string) Screener {
	return screener{
		url: url,
	}
}

func (screener screener) IsBlacklisted(ctx context.Context, addrs map[string]blockchain.Chain) (map[string]bool, error) {
	res := map[string]bool{}

	// Return false for all addrs if the url is empty
	if screener.url == "" || len(addrs) == 0 {
		for addr := range addrs {
			res[addr] = false
		}
		return res, nil
	}

	// Generate the request body
	client := new(http.Client)
	requestData := make([]request, 0, len(addrs))
	for addr, chain := range addrs {
		// todo : what about l2 chains, should we use the specific l2 chain name or ethereum ?
		requestData = append(requestData, request{
			Address: addr,
			Chain:   string(chain.Name()),
		})
	}
	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("[screener] unable to marshal request, err = %v", err)
	}
	input := bytes.NewBuffer(data)

	// Construct the request
	tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(tctx, "POST", screener.url, input)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	// Send the request and parse the response
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("[screener] error sending request, err = %v", err)
	}

	// from screener status code 201 or 200 is expected both are handled
	if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
		// get error and convert to string
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("[screener] error reading response body, err = %v", err)
		}
		errormsg := string(bodyBytes)
		return nil, fmt.Errorf("[screener] invalid status code, expect 201, got %v with error : %v", resp.StatusCode, errormsg)
	}

	// Parse the response
	var resps []response
	if err := json.NewDecoder(resp.Body).Decode(&resps); err != nil {
		return nil, fmt.Errorf("[screener] unexpected response, %v", err)
	}
	defer resp.Body.Close()

	for _, r := range resps {
		res[r.Address] = r.IsBlacklisted
	}
	return res, nil
}
