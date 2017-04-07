package common

import (
	"testing"

	log "github.com/Sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestPowershell(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("powershell_junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Powershell wrapper test suite",
		[]Reporter{junitReporter})
}

var _ = Describe("Controller", func() {

	type TestCase struct {
		message string
	}
	DescribeTable("output encoding works correctly",
		func(t TestCase) {
			//$OutputEncoding = [System.Text.Encoding]::UTF8;
			stdout, stderr, err := CallPowershell("Write-Host", t.message)
			Expect(err).ToNot(HaveOccurred())
			Expect(stdout).To(Equal(t.message))
			Expect(stderr).To(Equal(""))
		},
		Entry("ASCII", TestCase{
			message: "Hello!",
		}),
		FEntry("UTF8", TestCase{
			message: "Héllo!",
		}),
		Entry("UTF8", TestCase{
			message: "Hśllo!",
		}),
	)

	It("trims leading and trailing whitespace from powershell output", func() {
		msg := "  Hello there!  "
		stdout, stderr, err := CallPowershell("Write-Host", msg)
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).To(Equal("Hello there!"))
		Expect(stderr).To(Equal(""))
	})
})
