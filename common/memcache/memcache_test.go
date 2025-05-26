package memcache_test

import (
	"time"

	"github.com/catalogfi/blockchain/common/memcache"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemCache", func() {
	var cache memcache.Cache[string]
	var err error

	BeforeEach(func() {
		cache, err = memcache.NewMemCache[string](5 * time.Second)
		Expect(err).Should(BeNil())
	})

	Context("when setting and getting a value", func() {
		It("should store and retrieve the value", func() {
			result := cache.Set("test", "test")
			Expect(result).To(BeTrue())

			value, found := cache.Get("test")
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("test"))
		})
	})

	Context("when the TTL expires", func() {
		It("should not retrieve the value after expiration", func() {
			result := cache.Set("test", "test")
			Expect(result).To(BeTrue())

			time.Sleep(6 * time.Second)

			_, found := cache.Get("test")
			Expect(found).To(BeFalse())
		})
	})
})
