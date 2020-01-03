package main


import (
	      "github.com/9z25/go-bitcoind"
	      "log"
        "github.com/gorilla/mux"
        "io/ioutil"
        "net/http"
        "encoding/json"
        "fmt"
)

const (
	SERVER_HOST        = "127.0.0.1"
	SERVER_PORT        = 15151
	USER               = "testuser"
	PASSWD             = "test"
	USESSL             = false
)



type TaoNode struct {
  Rpc *bitcoind.Bitcoind
}


//example
type Book struct {
  Result string `json:"result"`
}

type SendTo struct {
  Recipient string `json:"recipient"`
  Amount    float64 `json:"amount"`
}

type RawTx struct {
  Tx string `json:"tx"`
}

// A ScriptSig represents a scriptsyg
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

var books []Book
var Node *bitcoind.Bitcoind


// Connect: Connect to RPC
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
json.NewEncoder(w).Encode(Book{Result:"Access denied.",})

return false
}
return true
}

//DecodeRawTransaction: decode raw transaction, send back json
func DecodeRawTransaction(w http.ResponseWriter, r *http.Request) {

  a := authorized(w,r)
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

    

    w.Header().Set("Content-Type","application/json")
    json.NewEncoder(w).Encode(res)
  
  }


//SendRawTransaction: broadcast transaction
func SendRawTransaction(w http.ResponseWriter, r *http.Request) {

  a := authorized(w,r)
  if a != true {
  return
  }
  w.Header().Set("Content-Type","application/json")

  
  var hash RawTx

  _ = json.NewDecoder(r.Body).Decode(&hash)
  
    res, err := Node.SendRawTransaction(hash.Tx)
    if err != nil {
      fmt.Println(err)
      }

      fmt.Println(res)
    
    
      var book Book
      book.Result = res

    
    json.NewEncoder(w).Encode(book)

  
  }


  // GetAddress : get current address
func GetAddress(w http.ResponseWriter, r *http.Request) {

a := authorized(w,r)
if a != true {
return
}
w.Header().Set("Content-Type","application/json")



  address, err := Node.GetAccountAddress("")

  if err != nil {
  fmt.Println(err)
  }

  var page Book
  page.Result = address


  
  json.NewEncoder(w).Encode(page)

}

//GetNewAddress, get a new address for user
func GetNewAddress(w http.ResponseWriter, r *http.Request) {

a := authorized(w,r)
if a != true {
return
}
w.Header().Set("Content-Type","application/json")

  newAddress, err := Node.GetNewAddress("")

  if err != nil {
  fmt.Println(err)
  }


  var book Book
  book.Result = newAddress


  
  json.NewEncoder(w).Encode(book)

}

//SendToAddress: send Tao to external address
func SendToAddress(w http.ResponseWriter, r *http.Request) {

a := authorized(w,r)
if a != true {
return
}

  w.Header().Set("Content-Type","application/json")

  d, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }

    var withdraw SendTo
    if err := json.Unmarshal(d, &withdraw); err != nil {
        panic(err)
    }
  
   txid, err := Node.SendToAddress(withdraw.Recipient, withdraw.Amount,"tao-rolls alpha","tao-rolls alpha")
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
/*
       h := handlers.AllowedHeaders([]string{"Content-Type","Authorization"})
       m := handlers.AllowedMethods([]string{"GET","POST"})
       o := handlers.AllowedOrigins([]string{"https://google.com",})
*/
       // example

       books = append(books, Book{Result: "testing"})
       books = append(books, Book{Result: "test2"})



//    handler := c.Handler(r)

       // Route Handlers / Endpoints

       r.HandleFunc("/api/getnewaddress/", GetNewAddress).Methods("GET")
       r.HandleFunc("/api/getaddress/", GetAddress).Methods("GET")
       r.HandleFunc("/api/sendtoaddress/", SendToAddress).Methods("POST")
       r.HandleFunc("/api/sendrawtransaction/", SendRawTransaction).Methods("POST")
       r.HandleFunc("/api/decoderawtransaction/", DecodeRawTransaction).Methods("POST")
       log.Fatal(http.ListenAndServe(":8000", r))
       //log.Fatal(http.ListenAndServe(":8000",r))


}
