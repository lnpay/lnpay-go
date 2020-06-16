package lnpay

// Webhook types

// "wallet_created" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WalletCreatedWebhookEvent struct {
	CreatedAt int    `json:"created_at"`
	ID        string `json:"id"`
	Event     Event  `json:"event"`
	Data      struct {
		Wal Wal `json:"wal"`
	} `json:"data"`
}

// "wallet_send" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WalletSendWebhookEvent struct {
	CreatedAt int    `json:"created_at"`
	ID        string `json:"id"`
	Event     Event  `json:"event"`
	Data      struct {
		Wtx Wtx `json:"wtx"`
	} `json:"data"`
}

// "wallet_receive" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WalletInternalTransferWebhookEvent struct {
	ID        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	Event     Event  `json:"event"`
	Data      struct {
		Wtx Wtx `json:"wtx"`
	} `json:"data"`
}

// "wallet_transfer_IN/OUT" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WalletReceiveWebhookEvent struct {
	CreatedAt int    `json:"created_at"`
	ID        string `json:"id"`
	Event     Event  `json:"event"`
	Data      struct {
		Wtx Wtx `json:"wtx"`
	} `json:"data"`
}

// "paywall_created" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#paywalls
type PaywallCreatedWebhookEvent struct {
	ID        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	Event     Event  `json:"event"`
	Data      struct {
		Pywl Pywl `json:"pywl"`
	} `json:"data"`
}

// "paywall_conversion" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#paywalls
type PaywallConversionWebhookEvent struct {
	ID        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	Event     Event  `json:"event"`
	Data      struct {
		Pywl Pywl `json:"pywl"`
		Wtx  Wtx  `json:"wtx"`
	} `json:"data"`
}

// primary types

type Wal struct {
	ID         string      `json:"id"`
	UserLabel  string      `json:"user_label"`
	CreatedAt  int         `json:"created_at"`
	UpdatedAt  int         `json:"updated_at"`
	Balance    int64       `json:"balance"`
	StatusType StatusType  `json:"statusType"`
	AccessKeys *AccessKeys `json:"accessKeys"`
}

type Wtx struct {
	NumSatoshis int                    `json:"num_satoshis"`
	UserLabel   string                 `json:"user_label"`
	CreatedAt   int                    `json:"created_at"`
	ID          string                 `json:"id"`
	Wal         Wal                    `json:"wal"`
	WtxType     WtxType                `json:"wtxType"`
	LnTx        LnTx                   `json:"lnTx"`
	PassThru    map[string]interface{} `json:"passThru"`
}

type LnTx struct {
	ID              string                 `json:"id"`
	CreatedAt       int                    `json:"created_at"`
	DestPubkey      string                 `json:"dest_pubkey"`
	PaymentRequest  string                 `json:"payment_request"`
	RHashDecoded    string                 `json:"r_hash_decoded"`
	Memo            string                 `json:"memo"`
	DescriptionHash string                 `json:"description_hash"`
	NumSatoshis     int64                  `json:"num_satoshis"`
	Expiry          int                    `json:"expiry"`
	ExpiresAt       int                    `json:"expires_at"`
	PaymentPreimage string                 `json:"payment_preimage"`
	Settled         int                    `json:"settled"`
	SettledAt       int                    `json:"settled_at"`
	IsKeysend       bool                   `json:"is_keysend"`
	CustomRecords   map[string]interface{} `json:"custom_records"`
}

type Pywl struct {
	DestinationURL string                 `json:"destination_url"`
	Memo           string                 `json:"memo"`
	ShortURL       string                 `json:"short_url"`
	NumSatoshis    int64                  `json:"lnd_value"`
	CreatedAt      int                    `json:"created_at"`
	UpdatedAt      int                    `json:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata"`
	ID             string                 `json:"id"`
	PaywallLink    string                 `json:"paywall_link"`
	CustyDomain    struct {
		DomainName string `json:"domain_name"`
	} `json:"custyDomain"`
	StatusType  StatusType `json:"statusType"`
	PaywallType struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
	} `json:"paywallType"`
	Template struct {
		Layout string `json:"layout"`
	} `json:"template"`
	LinkExpRule struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		TimeMinutes int    `json:"time_minutes"`
	} `json:"linkExpRule"`
}

// Secondary/helper types

type Event struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type AccessKeys struct {
	WalletAdmin   []string `json:"Wallet Admin"`
	WalletInvoice []string `json:"Wallet Invoice"`
	WalletRead    []string `json:"Wallet Read"`
}

type StatusType struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type WtxType struct {
	Layer       string `json:"layer"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
