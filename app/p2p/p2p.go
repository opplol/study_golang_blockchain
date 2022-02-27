package p2p

import (
	"fmt"
	"go_crypo_coin/blockchain"
	"go_crypo_coin/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	fmt.Printf("%s wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)

}

func AddPeer(address, port, openport string, broadcast bool) {
	fmt.Printf("%s wants to connect to port %s \n", openport, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openport), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	if broadcast {
		broadcastNewPeer(p)
		return
	}
	sendNewestBlock(p)

}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	fmt.Printf("broadcastNewTX: running :: peer %v \n", Peers.v)

	for _, p := range Peers.v {
		fmt.Printf("%s: broadcast New TX \n", p.port)
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}

}
