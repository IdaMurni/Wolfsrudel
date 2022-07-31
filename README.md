# PoW
Pow Implementation <br/>
run Blockchain: <br/>
target dir blockchain_server http://0.0.0.0:5000
````
go run *go
````
available BlockchainServer API<br>
````
func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	http.HandleFunc("/mine", bcs.Mine)
	http.HandleFunc("/mine/start", bcs.StartMine)
	http.HandleFunc("/amount", bcs.Amount)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}
````
<br/>

run Wallet: <br/>
target dir wallet_server  http://0.0.0.0:8080
````
go run *go
````

available wallet API: <br/>
````
func (waletServer *WalletServer) Run() {
	http.HandleFunc("/", waletServer.Index)
	http.HandleFunc("/wallet", waletServer.Wallet)
	http.HandleFunc("/wallet/amount", waletServer.WalletAmount)
	http.HandleFunc("/transaction", waletServer.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(waletServer.Port())), nil))
}
````


