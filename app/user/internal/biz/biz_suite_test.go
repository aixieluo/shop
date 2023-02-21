package biz_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBiz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Biz Suite")
}

var ctl *gomock.Controller
var cleaner func()
var ctx context.Context

var _ = BeforeEach(func() {
	ctl = gomock.NewController(GinkgoT())
	cleaner = ctl.Finish
	ctx = context.Background()
})
var _ = AfterEach(func() {
	cleaner()
})
