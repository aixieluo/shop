package biz_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"shop/app/user/internal/mocks/mrepo"
	"shop/app/user/internal/testdata"

	"shop/app/user/internal/biz"
)

var _ = Describe("User", func() {
	var userCase *biz.UserUsecase
	var mUserRepo *mrepo.MockUserRepo
	BeforeEach(func() {
		mUserRepo = mrepo.NewMockUserRepo(ctl)
		userCase = biz.NewUserUsecase(mUserRepo, nil)
	})
	It("Create", func() {
		info := testdata.User()
		mUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(info, nil)
		l, err := userCase.Create(ctx, info)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(l.ID).To(Equal(uint64(1)))
	})
	It("List", func() {
		info := testdata.User()
		var users []*biz.User
		users = append(users, info)
		mUserRepo.EXPECT().List(ctx, 1, 10).Return(users, 1, nil)
		us, total, err := userCase.List(ctx, 1, 10)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(total).Should(Equal(1))
		Ω(us).Should(HaveLen(1))
		Ω(us[0].ID).Should(Equal(info.ID))
	})
	It("Get", func() {
		info := testdata.User()
		mUserRepo.EXPECT().FindByID(ctx, uint(info.ID)).Return(info, nil)
		u, err := userCase.Get(ctx, uint(info.ID))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(u.ID).Should(Equal(info.ID))
	})
	It("Update", func() {
		info := testdata.User()
		mUserRepo.EXPECT().Update(ctx, uint(info.ID), info).Return(true, nil)
		b, err := userCase.Update(ctx, uint(info.ID), info)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(b).Should(BeTrue())
	})
})
