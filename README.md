# bitcoincli
## JSON RPC Bitcoin Core Go Client

```
func main() {
    bitcliconfig := bitcoincli.NewDefaultBitcoinCliConfig().
      WithUser("YOUR_USER_NAME").                               // form bitcoin.conf
      WithPassword("YOUR_PASSWORD").
      WithWalletNotify(&bitcoincli.WalletNotifyConfig{
        Port: 17331,
    })

    client := bitcoincli.NewBitcoinCli(*bitcliconfig, []func(trans bitcoincli.RawTransaction){RawTransactionCallback})

    walletBalance, _ := client.GetBalance("", 1)
    
    // You can create any rpc request from https://developer.bitcoin.org/reference/intro.html
    rpcWalletBalance, _ := client.Rpc().CallWallet("", "getbalance", []interface{}{"*", 1})  

    fmt.Println(walletBalance, rpcWalletBalance)
}


func RawTransactionCallback(trans bitcoincli.RawTransaction) {	// receiving transaction updates
	fmt.Println("TX UPD", trans.TxId, "confirms", trans.Confirmations)
}
```

