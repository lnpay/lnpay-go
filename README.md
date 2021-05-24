# lnpay [![](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/fiatjaf/lnpay-go)

Golang API wrapper for https://lnpay.co/.

## Install

```
go get github.com/fiatjaf/lnpay-go
```

## Usage

```go
lnpaySecretKey := lnpay.TEST_KEY
// use your key here: "sak_..."
// you can find it at https://lnpay.co/developers/dashboard

client := lnpay.NewClient(lnpaySecretKey)
wallet, _ := client.CreateWallet("first wallet")

lntx, _ := wallet.Invoice(lnpay.InvoiceParams{
    Memo: "wallet funding",
    NumSatoshis: 1000,
    PassThru: map[string]interface{}{
        "useless_data": 123,
    },
})

fmt.Printf("created invoice with lntx_id %s and payment hash %s.\n", lntx.ID, lntx.RHashDecoded)

details, _ := wallet.Details()
fmt.Printf("wallet %s (%s) has a balance of %d satoshis.\n", details.ID, details.UserLabel, details.Balance)

wtx, _ := wallet.Pay(PayParams{PaymentRequest: "lnbc1..."})
fmt.Printf("sent payment of %d satoshis to node %s.\n", wtx.NumSatoshis, wtx.LnTx.DestPubKey)
```
