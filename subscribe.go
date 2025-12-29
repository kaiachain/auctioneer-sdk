package sdk

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
)

var SEARCHER_SUBSCRIPTION_HEADER = "0xauc710neer"

// GetSearcherSubscriptionNonceMsg returns subscription fingerprint
func GetSearcherSubscriptionNonceMsg(nonce uint64) common.Hash {
	msg := fmt.Sprintf("%s-%d", SEARCHER_SUBSCRIPTION_HEADER, nonce)
	return common.BytesToHash(crypto.Keccak256([]byte(msg)))
}

// GetSubscriptionNonce retrieves searcher's subscription nonce
func GetSubscriptionNonce(host string, addr common.Address) (uint64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/subscribe/nonce?address=%s", host, addr.String()))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	num, err := strconv.ParseUint(string(body), 0, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// GetDialUrl returns a websocket url which targets subscription request
func GetDialUrl(host string, searcherKey *ecdsa.PrivateKey) (string, error) {
	searcherAddr := crypto.PubkeyToAddress(searcherKey.PublicKey)
	nonce, err := GetSubscriptionNonce(host, searcherAddr)
	if err != nil {
		return "", err
	}
	h, err := crypto.Sign(GetSearcherSubscriptionNonceMsg(nonce).Bytes(), searcherKey)
	if err != nil {
		return "", err
	}

	// trim `http` and `https` prefix
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	u := url.URL{
		Scheme:   "wss",
		Host:     host,
		Path:     "api/v1/subscribe/pendingtxs",
		RawQuery: fmt.Sprintf("sig=%s&nonce=%d", hexutil.Encode(h), nonce),
	}
	return u.String(), nil
}
