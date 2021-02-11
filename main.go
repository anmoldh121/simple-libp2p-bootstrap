package main

import (
	"context"
	"fmt"
	mrand "math/rand"

	libp2p "github.com/libp2p/go-libp2p"
	// circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
	net "github.com/libp2p/go-libp2p-core/network"

	// "github.com/libp2p/go-libp2p-core/peerstore"
	crypto "github.com/libp2p/go-libp2p-crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	noise "github.com/libp2p/go-libp2p-noise"
	// quic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
	// peer "github.com/libp2p/go-libp2p-core/peer"
)

var kahdemia *dht.IpfsDHT

type notifier struct {
	DHT  *dht.IpfsDHT
	host *host.Host
}

func (n *notifier) Connected(netw net.Network, conn net.Conn) {
	// n.host.Peerstore().AddAddr(conn.RemotePeer(), conn.RemoteMultiaddr(), peerstore.PermanentAddrTTL)
	fmt.Printf("Connection established with %v\n", conn.RemotePeer().Pretty())
	fmt.Println(n.DHT.RoutingTable().ListPeers())
	fmt.Println((*n.host).Peerstore().Peers())
}

func (n *notifier) Listen(netw net.Network, addr ma.Multiaddr)       {}
func (n *notifier) ListenClose(netw net.Network, addr ma.Multiaddr)  {}
func (n *notifier) Disconnected(netw net.Network, conn net.Conn)     {}
func (n *notifier) OpenedStream(netw net.Network, stream net.Stream) {}
func (n *notifier) ClosedStream(netw net.Network, stream net.Stream) {}

func main() {
	ctx := context.Background()
	// sourceMultiAddr := ma.StringCast("/ip4/0.0.0.0/tcp/4000")

	r := mrand.New(mrand.NewSource(int64(10)))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/4000", "/ip4/0.0.0.0/udp/4000/quic"),
		// libp2p.Transport(quic.NewTransport),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Identity(prvKey),
		// libp2p.EnableRelay(circuit.OptHop),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("This node: ", host.ID().Pretty(), " ", host.Addrs())

	kahdemia, err = dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	ntf := notifier{kahdemia, &host}

	host.Network().Notify(&ntf)

	select {}
}
