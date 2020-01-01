package main


import (
	      "github.com/toorop/go-bitcoind"
	      "log"
        "github.com/gorilla/mux"
        "net/http"
        "encoding/json"
        "fmt"
)

const (
	SERVER_HOST        = "127.0.0.1"
	SERVER_PORT        = 15151
	USER               = ""
	PASSWD             = ""
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
  Recipient string `json:"address"`
  Amount    float64 `json:"amount"`
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

//SendRawTransaction: broadcast transaction
func SendRawTransaction(w http.ResponseWriter, r *http.Request) {

  a := authorized(w,r)
  if a != true {
  return
  }

  fmt.Println(r.Body)
  
  /*
    address, err := Node.SendRawTransaction("")
  
    if err != nil {
    fmt.Println(err)
    }
  
  
  
  
    var page Book
    page.Result = address
  
  
    w.Header().Set("Content-Type","application/json")
    json.NewEncoder(w).Encode(page)
  */
  }


  // GetAddress : get current address
func GetAddress(w http.ResponseWriter, r *http.Request) {

a := authorized(w,r)
if a != true {
return
}


  address, err := Node.GetAccountAddress("")

  if err != nil {
  fmt.Println(err)
  }




  var page Book
  page.Result = address


  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(page)

}

//GetNewAddress, get a new address for user
func GetNewAddress(w http.ResponseWriter, r *http.Request) {

a := authorized(w,r)
if a != true {
return
}

  newAddress, err := Node.GetNewAddress("")

  if err != nil {
  fmt.Println(err)
  }


  var book Book
  book.Result = newAddress


  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(book)

}

//SendToAddress: send Tao to external address
func SendToAddress(w http.ResponseWriter, r *http.Request) {

a := authorized(w,r)
if a != true {
return
}

  w.Header().Set("Content-Type","application/json")
  var book Book
  var withdraw SendTo
  _ = json.NewDecoder(r.Body).Decode(&withdraw)



   txid, err := Node.SendToAddress(withdraw.Recipient, withdraw.Amount,"tao-rolls alpha","tao-rolls alpha")
                log.Println(err, txid)




  book.Result = txid
  //book.Result = "test666"

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
       r.HandleFunc("/api/sendrawtransaction/", SendToAddress).Methods("POST")
       log.Fatal(http.ListenAndServe(":8000", r))
       //log.Fatal(http.ListenAndServe(":8000",r))


}
