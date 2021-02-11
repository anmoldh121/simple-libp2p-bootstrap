package main

import (
	"context"
	"fmt"
	mrand "math/rand"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
	noise "github.com/libp2p/go-libp2p-noise"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

func main() {
	ctx := context.Background()
	sourceMultiAddr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/4000")

	r := mrand.New(mrand.NewSource(int64(10)))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Identity(prvKey),
		libp2p.EnableRelay(circuit.OptHop),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("This node: ", host.ID().Pretty(), " ", host.Addrs())

	_, err = dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	select {}
}
