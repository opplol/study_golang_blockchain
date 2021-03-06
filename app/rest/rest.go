package rest

import (
	"encoding/json"
	"fmt"
	"go_crypo_coin/blockchain"
	"go_crypo_coin/p2p"
	"go_crypo_coin/utils"
	"go_crypo_coin/wallet"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type uRLDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func (u uRLDescription) String() string {
	return "URL Description in String()"
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type errorRepsonse struct {
	ErrorMessage string `json:"errormessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type addPeerPayload struct {
	Address string
	Port    string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []uRLDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentaion",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{Hash}"),
			Method:      "GET",
			Description: "Get Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSocket",
		},
	}

	//hard way for render json
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)

	// simple way for render json
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		newBlock := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorRepsonse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})

}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		utils.HandleErr(json.NewEncoder(rw).Encode(balanceResponse{address, amount}))
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UnspandTxOutsByAddress(address, blockchain.Blockchain())))
	}

}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempoll().Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	fmt.Printf("Running transaction\n")
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	fmt.Printf("Running transaction:: Add MemPool\n")
	tx, err := blockchain.Mempoll().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorRepsonse{err.Error()})
		return
	}
	fmt.Printf("Running transaction:: BroadCast %v\n", tx)
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{address})
}

func addPeers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}

}

func Start(aPort int) {
	handler := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	handler.Use(jsonContentTypeMiddleware, loggerMiddleware)
	handler.HandleFunc("/", documentation).Methods("GET")
	handler.HandleFunc("/status", status).Methods("GET")
	handler.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	handler.HandleFunc("/balance/{address}", balance).Methods("GET")
	handler.HandleFunc("/mempool", mempool).Methods("GET")
	handler.HandleFunc("/wallet", myWallet).Methods("GET")
	handler.HandleFunc("/transactions", transactions).Methods("POST")
	handler.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	handler.HandleFunc("/peers", addPeers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
