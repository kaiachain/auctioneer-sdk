package sdk

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
)

var SEARCHER_SUBSCRIPTION_HEADER = "0xauc710neer"

// GetSearcherSubscriptionNonceMsg returns subscription fingerprint
func GetSearcherSubscriptionNonceMsg(nonce uint64) common.Hash {
	msg := fmt.Sprintf("%s-%x", SEARCHER_SUBSCRIPTION_HEADER, nonce)
	return common.HexToHash(msg)
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
func GetDialUrl(host, path string, nonce uint64, searcherKey *ecdsa.PrivateKey) (string, error) {
	searcherAddr := crypto.PubkeyToAddress(searcherKey.PublicKey)
	nonce, err := GetSubscriptionNonce(host, searcherAddr)
	if err != nil {
		return "", err
	}
	h, err := crypto.Sign(GetSearcherSubscriptionNonceMsg(nonce).Bytes(), searcherKey)
	if err != nil {
		return "", err
	}
	u := url.URL{
		Scheme:   "ws",
		Host:     host,
		Path:     path,
		RawQuery: fmt.Sprintf("sig=%s&nonce=%d", hexutil.Encode(h), nonce),
	}
	return u.String(), nil
}
