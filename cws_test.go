package tune

import (
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestCSWStatter(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Buckets", func() {
		g.It("should split a bucket string if multiple segments exist", func() {
			b := "get.endpoint-name.subaction"

			dn, dv, m := metricNames(b)
			Expect(dn).To(Equal("get"))
			Expect(dv).To(Equal("endpoint-name"))
			Expect(m).To(Equal("subaction"))
		})

		g.It("should return the same string if multiple segments do not exist", func() {
			b := "endpoint-name"

			dn, dv, m := metricNames(b)
			Expect(dn).To(Equal("endpoint-name"))
			Expect(dv).To(Equal("endpoint-name"))
			Expect(m).To(Equal("endpoint-name"))
		})

		g.It("should handle a really long set of segments", func() {
			b := "get.some.endpoint.subaction.count"

			dn, dv, m := metricNames(b)
			Expect(dn).To(Equal("get"))
			Expect(dv).To(Equal("some.endpoint.subaction"))
			Expect(m).To(Equal("count"))
		})
	})
}
