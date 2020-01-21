package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/9z25/go-bitcoind"
	"github.com/gorilla/mux"
)

const (
	SERVER_HOST = "127.0.0.1"
	SERVER_PORT = 15151
	USER        = "testuser"
	PASSWD      = "test"
	USESSL      = false
)

//TaoNode struct for RPC Node
type TaoNode struct {
	Rpc *bitcoind.Bitcoind
}

// Result
type FmTao struct {
	Result string `json:"result"`
}

// LastTx
type LastTx struct {
	Type      string `json:"type"`
	Addresses string `json:"addresses"`
}

// TaoExplorer
type TaoExplorer struct {
	Address  string   `json:"address"`
	Sent     int      `json:"sent"`
	Received string   `json:"received"`
	Balance  string   `json:"balance"`
	lastTxs  []LastTx `json:"last_txs"`
}

//Book struct for handling json response
type Book struct {
	Result string `json:"result"`
}

//SendTo struct for handling json post data
type SendTo struct {
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

//RawTx struct for handling json post data
type RawTx struct {
	Tx string `json:"tx"`
}

// A ScriptSig represents a script signature
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// Vin represent an IN value
type Vin struct {
	Coinbase  string    `json:"coinbase"`
	Txid      string    `json:"txid"`
	Vout      int       `json:"vout"`
	ScriptSig ScriptSig `json:"scriptSig"`
	Sequence  uint32    `json:"sequence"`
}

// ScriptPubKey respresents the Script Public Key of the Vout
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

// Vout represent an OUT value
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// RawTransaction represents the structure of a raw transaction
type RawTransaction struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Version       uint32 `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

// TransactionDetails represents details about a transaction
type TransactionDetails struct {
	Account  string  `json:"account"`
	Address  string  `json:"address,omitempty"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Fee      float64 `json:"fee,omitempty"`
}

//Book map
var books []Book

//Node variable represents node
var Node *bitcoind.Bitcoind

// Connect to RPC
func (taoN TaoNode) Connect() *bitcoind.Bitcoind {

	connection, err := bitcoind.New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		log.Fatalln(err)
	}

	taoN.Rpc = connection
	return taoN.Rpc

}

func authorized(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("X-Csrf-Token") != "123" {
		json.NewEncoder(w).Encode(Book{Result: "Access denied."})
		return false
	}

	return true
}

//GetTransaction decode raw transaction, send back json
func GetTransaction(w http.ResponseWriter, r *http.Request) {

	a := authorized(w, r)
	if a != true {
		return
	}

	params := mux.Vars(r)
	txid := params["txid"]
	fmt.Println(txid)

	res, err := Node.GetTransaction(txid)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(res)
	json.NewEncoder(w).Encode(res)
}

//GetRawTransaction decode raw transaction, send back json
func GetRawTransaction(w http.ResponseWriter, r *http.Request) {

	a := authorized(w, r)
	if a != true {
		return
	}

	params := mux.Vars(r)
	txid := params["txid"]
	fmt.Println(txid)

	res, err := Node.GetRawTransaction(txid, true)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(res)
	json.NewEncoder(w).Encode(res)
}

//DecodeRawTransaction decode raw transaction, send back json
func DecodeRawTransaction(w http.ResponseWriter, r *http.Request) {

	a := authorized(w, r)
	if a != true {
		return
	}

	var hash RawTx

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(d, &hash); err != nil {
		panic(err)
	}

	res, err := Node.DecodeRawTransaction(hash.Tx)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

//SendRawTransaction broadcast transaction
func SendRawTransaction(w http.ResponseWriter, r *http.Request) {

	a := authorized(w, r)
	if a != true {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var hash RawTx

	_ = json.NewDecoder(r.Body).Decode(&hash)

	res, err := Node.SendRawTransaction(hash.Tx)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	fmt.Println(res)

	var book Book
	book.Result = res

	json.NewEncoder(w).Encode(book)

}

// GetAddress : get current address
func GetAddress(w http.ResponseWriter, r *http.Request) {

	a := authorized(w, r)
	if a != true {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	address, err := Node.GetAccountAddress("")
	if err != nil {
		fmt.Println(err)
	}

	var page Book
	page.Result = address

	json.NewEncoder(w).Encode(page)

}

//GetNewAddress get a new address for user
func GetNewAddress(w http.ResponseWriter, r *http.Request) {

	a := authorized(w, r)
	if a != true {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	newAddress, err := Node.GetNewAddress("")

	if err != nil {
		fmt.Println(err)
	}

	var book Book
	book.Result = newAddress

	json.NewEncoder(w).Encode(book)
}

//SendToAddress send Tao to external address
func SendToAddress(w http.ResponseWriter, r *http.Request) {
	return
	a := authorized(w, r)
	if a != true {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var withdraw SendTo
	if err := json.Unmarshal(d, &withdraw); err != nil {
		panic(err)
	}

	txid, err := Node.SendToAddress(withdraw.Recipient, withdraw.Amount, "tao-rolls alpha", "tao-rolls alpha")
	log.Println(err, txid)

	var book Book
	book.Result = txid
	json.NewEncoder(w).Encode(book)

}



func main() {

	var t TaoNode
	Node = TaoNode.Connect(t)

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/api/getnewaddress/", GetNewAddress).Methods("GET")
	r.HandleFunc("/api/getaddress/", GetAddress).Methods("GET")
	r.HandleFunc("/api/sendtoaddress/", SendToAddress).Methods("POST")
	r.HandleFunc("/api/gettransaction/{txid}", GetTransaction).Methods("GET")
	r.HandleFunc("/api/getrawtransaction/{txid}", GetTransaction).Methods("GET")
	r.HandleFunc("/api/sendrawtransaction/", SendRawTransaction).Methods("POST")
	r.HandleFunc("/api/decoderawtransaction/", DecodeRawTransaction).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":8000", r))
}
