package lnpay

import (
	"github.com/imroc/req"
	"reflect"
	"time"
)

const (
	BASE_URL = "https://api.lnpay.co/v1"
	TEST_KEY = "pak_O0iUMxk8kK_qUzkT4YKFvp1ZsUtp"
)

type Client struct {
	header req.Header
}

// NewClient is the first function you must call. Pass your main API key here.
// It will return a client you can later use to access wallets and transactions.
// You can find it at https://lnpay.co/developers/dashboard
func NewClient(key string) *Client {
	return &Client{
		header: req.Header{
			"X-Api-Key":    key,
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

// Transaction
func (c *Client) Transaction(lntxId string) (lnTx LnTx, err error) {
	resp, err := req.Get(BASE_URL+"/lntx/"+lntxId, c.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&lnTx)
	return
}

// CreateWallet creates a new wallet with a given descriptive label.
// It will return the wallet object which you can use to create invoices and payments.
// https://docs.lnpay.co/wallet/create-wallet
func (c *Client) CreateWallet(label string) (wal Wallet, err error) {
	resp, err := req.Post(BASE_URL+"/wallet", c.header, req.BodyJSON(struct {
		UserLabel string `json:"user_label"`
	}{label}))
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wal)
	return
}

// Wallet returns a wallet that was already created.
// Pass the wallet key you got when creating it.
// I can be either the admin, invoice or read-only key.
func (c *Client) Wallet(key string) *Wallet {
	return &Wallet{
		c,
		key,
		BASE_URL + "/wallet/" + key,
	}
}

type Wallet struct {
	*Client
	Key      string
	BASE_URL string
}

// Details returns basic information about a wallet, such as its id, label or balance.
// https://docs.lnpay.co/wallet/get-balance
func (w *Wallet) Details() (wal Wal, err error) {
	resp, err := req.Get(w.BASE_URL, w.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wal)
	return
}

// Transactions returns a list of the transactions associated with the wallet.
// https://docs.lnpay.co/wallet/get-transactions
func (w *Wallet) Transactions() (txs []Wtx, err error) {
	resp, err := req.Get(w.BASE_URL+"/transactions", w.header)
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&txs)
	return
}

type InvoiceParams struct {
	Memo        string `json:"memo"`         // the invoice description.
	NumSatoshis int64  `json:"num_satoshis"` // the invoice amount.
	Expiry      int64  `json:"expiry"`       // seconds, default 86400 (1 day).

	// custom data you may want to associate with this invoice. optional.
	PassThru map[string]interface{} `json:"passThru"`

	// base64-encoded. If this is provided, memo will be ignored.
	// don't use this if you don't know what it means.
	DescriptionHash string `json:"description_hash"`
}

// Invoice creates an invoice associated with this wallet.
// https://docs.lnpay.co/wallet/generate-invoice
func (w *Wallet) Invoice(params InvoiceParams) (lntx LnTx, err error) {
	resp, err := req.Post(w.BASE_URL+"/invoice", w.header, req.BodyJSON(&params))
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&lntx)
	return
}

type PayParams struct {
	// the BOLT11 payment request you want to pay.
	PaymentRequest string `json:"payment_request"`

	// custom data you may want to associate with this invoice. optional.
	PassThru map[string]interface{} `json:"passThru"`
}

// Pay pays a given invoice with funds from the wallet.
// https://docs.lnpay.co/wallet/pay-invoice
func (w *Wallet) Pay(params PayParams) (wtx Wtx, err error) {
	resp, err := req.Post(w.BASE_URL+"/withdraw", w.header, req.BodyJSON(&params))
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wtx)
	return
}

type TransferParams struct {
	Memo         string `json:"memo"`           // the transfer description.
	NumSatoshis  int64  `json:"num_satoshis"`   // the transfer amount.
	DestWalletId string `json:"dest_wallet_id"` // the key or id of the destination
}

// Transfer transfers between two lnpay.co wallets.
// https://docs.lnpay.co/wallet/transfers-between-wallets
func (w *Wallet) Transfer(params TransferParams) (wtx Wtx, err error) {
	resp, err := req.Post(w.BASE_URL+"/transfer", w.header, req.BodyJSON(&params))
	if err != nil {
		return
	}

	if resp.Response().StatusCode >= 300 {
		var reqErr Error
		resp.ToJSON(&reqErr)
		err = reqErr
		return
	}

	err = resp.ToJSON(&wtx)
	return
}

// Update the receiver transaction using lntx api.
// It will compare the receiver tx with the updated tx and replace the receiver,
// if any changes are detected.
func (tx *LnTx) Update(client *Client) (updated bool, err error) {
	lntx, err := client.Transaction(tx.ID)
	if err != nil {
		return false, err
	}
	if !tx.eq(lntx) {
		*tx = lntx
		return true, nil
	}
	return false, nil
}

// Expired checks if lntx is expired
func (tx *LnTx) Expired() bool {
	if tx.ExpiresAt < int(time.Now().Unix()) {
		return true
	}
	return false
}

// IsSettled checks if receiver lntx is settled yet.
// It will only update the tx if currently not settled.
func (tx *LnTx) IsSettled(client *Client) bool {
	if tx.Settled == 0 {
		updated, err := tx.Update(client)
		if err != nil {
			return false
		}
		if updated {
			// check if updated transaction settled
			if tx.Settled == 1 {
				return true
			}
		}
	}
	return false
}

// eq will compare the receiver transaction with another transaction
func (tx *LnTx) eq(lntx LnTx) (updated bool) {
	return reflect.DeepEqual(tx, &lntx)
}
