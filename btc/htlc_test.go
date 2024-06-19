package btc_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/blockchain/btc"

	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTLC (p2tr)", func() {

	It("should be able to generate HTLC address", func(ctx context.Context) {
		// generate random 32 bytes secret
		secret := make([]byte, 32)
		_, err := rand.Read(secret)
		Expect(err).To(BeNil())

		secretHash := sha256.Sum256(secret)

		// generate random keypair
		initiator, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		initiatorPubkey := initiator.PubKey().X().Bytes()

		redeemer, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())

		redeemerPubkey := redeemer.PubKey().X().Bytes()

		// generate HTLC address
		htlcWallet := btc.NewHTLCWallet(nil, &chaincfg.RegressionNetParams)

		htlcAddr, err := htlcWallet.Address(&btc.HTLC{
			InitiatorPubkey: initiatorPubkey,
			RedeemerPubkey:  redeemerPubkey,
			SecretHash:      secretHash[:],
			LockTime:        100,
		})

		Expect(err).To(BeNil())
		Expect(htlcAddr).ShouldNot(BeNil())

		// log the address
		color.Green("HTLC address: %s", htlcAddr)

	})

})
