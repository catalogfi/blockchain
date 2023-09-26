package btc_test

import (
	"bytes"
	"testing/quick"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/catalogfi/blockchain/btc"
	"github.com/tyler-smith/go-bip39"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin", func() {
	Context("keys", func() {
		It("should generate deterministic keys from mnemonic and user public key", func() {
			test := func() bool {
				entropy, err := bip39.NewEntropy(256)
				Expect(err).To(BeNil())
				mnemonic, err := bip39.NewMnemonic(entropy)
				Expect(err).To(BeNil())
				key, err := btcec.NewPrivateKey()
				Expect(err).To(BeNil())

				extendedKey1, err := btc.GenerateSystemPrivKey(mnemonic, key.PubKey().SerializeCompressed())
				Expect(err).To(BeNil())
				extendedKey2, err := btc.GenerateSystemPrivKey(mnemonic, key.PubKey().SerializeCompressed())
				Expect(err).To(BeNil())

				Expect(bytes.Equal(extendedKey1.Serialize(), extendedKey2.Serialize())).Should(BeTrue())
				return true
			}

			Expect(quick.Check(test, nil)).NotTo(HaveOccurred())
		})
	})
})
