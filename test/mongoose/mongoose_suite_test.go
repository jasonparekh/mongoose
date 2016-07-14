package mongoose_test

import (
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/mongoose/mock"
	"github.com/xeger/mongoose/test/mongoose"
)

var _ = Describe("mongoose dialect", func() {
	var v mongoose.Vehicle
	var w mongoose.Wheel

	URL, err := url.Parse("http://null-island.com")
	if err != nil {
		panic("bad text fixture")
	}

	BeforeEach(func() {
		v = &mongoose.MockVehicle{}
		w = &mongoose.MockWheel{}
		Â(v).On("Attach").With([]mongoose.Wheel{w})
		Â(v).On("Drive").With("north", 42.0).Return(*URL)
		Â(w).On("Diameter").Panic("big wheel keep on turnin'")
	})

	It("matches calls", func() {
		Expect(v.Drive("north", 42.0)).To(Equal(*URL))
		Expect(func() {
			w.Diameter()
		}).To(Panic())
	})

	It("matches basic-type params using equivalence", func() {
		Expect(v.Drive("north", 42)).To(Equal(*URL))
	})

	It("panics on unmatched calls", func() {
		Expect(func() {
			v.Drive("southeast", 12)
		}).To(Panic())
	})
})

func TestMongoose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongoose Suite")
}