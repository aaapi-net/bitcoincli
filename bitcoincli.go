package bitcoincli

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pebbe/zmq4"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// https://developer.bitcoin.org/reference/rpc/index.html

type BitcoinCli struct {
	isTest bool
	client *RpcClient
}

func NewBitcoinCli(config BitcoinCliConfig, walletTransactionListeners []func(trans RawTransaction)) *BitcoinCli {
	rpcClient, err := newRpcClient(config.Host, config.Port, config.User, config.Password, config.UseSsl, config.Timeout)
	if err != nil {
		panic(err)
	}

	result := &BitcoinCli{isTest: config.IsTest, client: rpcClient}
	if walletTransactionListeners != nil && len(walletTransactionListeners) > 0 {
		if config.WalletNotifyConfig == nil {
			log.Println("no WalletNotifyConfig")
		} else {
			result.startWalletTransactionsListener(config.WalletNotifyConfig.Host, config.WalletNotifyConfig.Port, walletTransactionListeners)
		}
	}
	_, err = result.LoadAllWallets()
	if err != nil {
		panic(err)
	}
	return result
}

// Handle wallet transactions.
// You need to add at bitcoin.conf:
// "walletnotify=curl -d %s http://[YOUR_HOST_HERE]:[YOUR_PORT_HERE]"
func (b *BitcoinCli) startWalletTransactionsListener(host string, port int, listeners []func(trans RawTransaction)) error {
	if len(listeners) == 0 {
		return errors.New("no listeners")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		txid := string(body)

		rawTrans, err := b.GetRawTransaction(txid)
		if err == nil {
			for _, lis := range listeners {
				lis(rawTrans)
			}
		} else {
			log.Println("startWalletTransactionsListener", err)
		}
	})
	go http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	return nil
}

func (b *BitcoinCli) IsTest() bool {
	return b.isTest
}

func (b *BitcoinCli) Rpc() *RpcClient {
	return b.client
}

// Handle all transactions.
// You need to enable zmq hashtx, for example at bitcoin.conf:
// "zmqpubhashtx=tcp://127.0.0.1:7334"  (zmqpubhashtx=tcp://[YOUR_HOST_HERE]:[YOUR_PORT_HERE])
func (b *BitcoinCli) startAllTransactionsListener(host string, port int, listeners []func(trans RawTransaction)) {
	go func () {
		xsub, _ := zmq4.NewSocket(zmq4.SUB)
		err := xsub.Connect("tcp://127.0.0.1:17334")
		if err != nil {
			panic(err)
		}
		err = xsub.SetSubscribe("hashtx")

		if err != nil {
			panic(err)
		}

		for {

			msg, err := xsub.RecvMessageBytes(0)
			if err != nil {
				fmt.Println("transaction listener error: ", err)
				continue
			}

			//msgType := string(msg[0])
			msgBody := hex.EncodeToString(msg[1])

			trans, err := b.GetRawTransaction(msgBody)

			if err == nil {
				for _, lis := range listeners {
					lis(trans)
				}
			} else {
				fmt.Println(msgBody)
			}
		}
	}()
}

func (b *BitcoinCli) SendToAddress(walletName string, addr string, amount float64) (txid string, err error) {
	r, err := b.client.CallWallet(walletName,"sendtoaddress", []interface{}{addr, amount})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &txid)
	return
}

func (b *BitcoinCli) SendMany(walletName string, addrToAmountMap map[string]float64) (txid string, err error) {
	r, err := b.client.CallWallet(walletName,"sendmany", []interface{}{"", addrToAmountMap})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &txid)
	return
}

func (b *BitcoinCli) SendToAddressWithInfo(walletName string, addr string, amount float64, info string, blockchainInfo string) (txid string, err error) {
	r, err := b.client.CallWallet(walletName,"sendtoaddress", []interface{}{addr, amount, info, blockchainInfo})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &txid)
	return
}

// getbalance :
// Returns the total available balance.
func (b *BitcoinCli) GetBalance(walletName string, minConfirms int) (balance float64, err error) {
	r, err := b.client.CallWallet(walletName,"getbalance", []interface{}{"*", minConfirms})
	if err = handleError(err, &r); err != nil {
		return
	}
	balance, err = strconv.ParseFloat(string(r.Result), 64)
	return
}


func (b *BitcoinCli) CreateWallet(id string) (walletName WalletWithWarning, err error) {
	r, err := b.client.Call("createwallet", []string{id})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &walletName)
	return
}

// listwallets :
// Returns a list of currently loaded wallets.
func (b *BitcoinCli) ListWallets() (wallets []string, err error) {
	r, err := b.client.Call("listwallets", []string{})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &wallets)

	return
}

// listwallets :
// Returns a list of wallets in the wallet directory.
func (b *BitcoinCli) ListWalletDir() (wallets []string, err error) {
	r, err := b.client.Call("listwalletdir", []string{})
	if err = handleError(err, &r); err != nil {
		return
	}

	var walletsList ListWalletDir
	err = json.Unmarshal(r.Result, &walletsList)

	for _, k := range walletsList.Wallets {
		wallets = append(wallets, k.Name)
	}
	return
}

// loadwallet :
// Loads a wallet from a wallet file or directory.
func (b *BitcoinCli) LoadWallet(name string)  (wallet WalletWithWarning, err error) {
	r, err := b.client.Call("loadwallet", []string{name})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &wallet)
	return
}

func (b *BitcoinCli) LoadAllWallets() (warnings []string, err error) {
	allWallets, err := b.ListWalletDir()


	if err != nil {
		return
	}

	loadedWallets, err := b.ListWallets()

	walletsToLoad := diffStringSlices(allWallets, loadedWallets)
	fmt.Println(walletsToLoad)

	if err != nil {
		return
	}

	var wallet WalletWithWarning
	for _, w := range walletsToLoad {
		wallet, err = b.LoadWallet(w)
		if err != nil {
			return
		}

		if len(wallet.Warning) != 0 {
			warnings = append(warnings, wallet.Warning)
		}
	}

	return
}

func diffStringSlices(source, with []string) []string {
	target := map[interface{}]bool{}
	for _, x := range with{
		target[x] = true
	}

	var result []string
	for _, x := range source {
		if _, ok := target[x]; !ok {
			result = append(result, x)
		}
	}

	return result
}

// unloadwallet :
// Unloads the wallet referenced by the request endpoint otherwise unloads the wallet specified in the argument.
func (b *BitcoinCli) UnloadWallet(name string)  (wallet WalletWithWarning, err error) {
	r, err := b.client.Call("unloadwallet", []string{name})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &wallet)
	return
}


// gettransaction :
// Get detailed information about in-wallet transaction <txid>
func (b *BitcoinCli) GetTransaction(walletName string, txid string) (trans Transaction, err error) {
	r, err := b.client.CallWallet(walletName,"gettransaction", []string{txid})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &trans)

	return
}

// getrawtransaction :
// By default this function only works for mempool transactions. When called with a blockhash argument, getrawtransaction will return the transaction if the specified block is available and the transaction
//is found in that block. When called without a blockhash argument, getrawtransaction will return the transaction if it is in the mempool, or if -txindex is enabled and the transaction is in a block in the blockchain.
func (b *BitcoinCli) GetRawTransaction(txid string) (trans RawTransaction, err error) {
	r, err := b.client.Call("getrawtransaction", []interface{}{txid, true})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &trans)

	return
}


func (b *BitcoinCli) GetNewAddress(walletName string, label ...string) (addr string, err error) {
	if len(label) > 1 {
		err = errors.New("Bad parameters for GetNewAddress: you can set 0 or 1 label")
		return
	}
	r, err := b.client.CallWallet(walletName,"getnewaddress", label)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &addr)
	return
}

// GetWalletAddress Returns bitcoin address for receiving
// payments to this wallet.
func (b *BitcoinCli) GetWalletAddress(walletName string) (address AddressInfo, err error) {
	r, err := b.ListPublicAddressesInfoByWallet(walletName)

	if len(r) > 0 {
		return r[0], err
	}
	return AddressInfo{}, err
}

// listaddressgroupings :
// Lists groups of addresses which have had their common ownership made public by common use as inputs or as the resulting change in past transactions
func (b *BitcoinCli) ListPublicAddressesInfoByWallet(walletName string) (result []AddressInfo, err error) {
	r, err := b.client.CallWallet(walletName,"listaddressgroupings", []string{})
	if err = handleError(err, &r); err != nil {
		return
	}

	var addresses [][][]interface{}
	err = json.Unmarshal(r.Result, &addresses)

	if err == nil && len(addresses) > 0 {
		for _, addrgroup := range addresses {
			for _, v := range addrgroup {
				wi := AddressInfo{
					Address: v[0].(string),
					Amount:  v[1].(float64),
				}

				if len(v) > 2 {
					wi.Label = v[2].(string)
				}

				result = append(result, wi)
			}
		}

	}

	return
}

// listaddressgroupings :
// Lists groups of addresses which have had their common ownership made public by common use as inputs or as the resulting change in past transactions
func (b *BitcoinCli) ListPublicAddressesMapByWallet(walletName string) (result map[string]float64, err error) {
	r, err := b.client.CallWallet(walletName,"listaddressgroupings", []string{})
	if err = handleError(err, &r); err != nil {
		return
	}

	var addresses [][][]interface{}
	err = json.Unmarshal(r.Result, &addresses)

	result = make(map[string]float64)
	if err == nil && len(addresses) > 0 {
		for _, v := range addresses[0] {
			result[v[0].(string)] = v[1].(float64)
		}
	}

	return
}

// listlabels :
// LReturns the list of all labels, or labels that are assigned to addresses with a specific purpose.
func (b *BitcoinCli) ListLabels(walletName string) (labels []string, err error) {
	r, err := b.client.CallWallet(walletName,"listlabels", []string{})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &labels)

	return
}

// getaddressesbylabel :
// Returns the list of addresses assigned the specified label.
func (b *BitcoinCli) GetAddressesByLabel(walletName, label string) (addresses []string, err error) {
	r, err := b.client.CallWallet(walletName,"getaddressesbylabel", []string{label})
	if err = handleError(err, &r); err != nil {
		return
	}

	var addressesMap map[string]interface{}
	err = json.Unmarshal(r.Result, &addressesMap)

	for addr, _ := range addressesMap {
		addresses = append(addresses, addr)
	}
	return
}

func (b *BitcoinCli) GetBalanceAddress(walletName, address string) (balance float64, err error) {
	list, err := b.ListPublicAddressesInfoByWallet(walletName)
	if err != nil {
		return
	}

	for _, w := range list {
		if w.Address == address {
			balance = w.Amount
			return
		}
	}
	err = errors.New("not found")
	return
}

func (b *BitcoinCli) GetReceivedByAddress(walletName string, addr string, minConfirms int) (amount float64, err error) {
	r, err := b.client.CallWallet(walletName,"getreceivedbyaddress", []interface{}{addr, minConfirms})
	if err = handleError(err, &r); err != nil {
		return
	}
	amount, err = strconv.ParseFloat(string(r.Result), 64)
	return
}
