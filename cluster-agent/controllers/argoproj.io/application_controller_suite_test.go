package argoprojio_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApplicationController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Controller Suite")
}
