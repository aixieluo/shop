package data_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"shop/app/user/internal/biz"
	"shop/app/user/internal/data"
)

var _ = Describe("User CURD", func() {
	var ro biz.UserRepo
	var uD *biz.User
	BeforeEach(func() {
		ro = data.NewUserRepo(Db, nil)
		uD = &biz.User{
			ID:       1,
			Mobile:   "17521066239",
			Password: "123123",
			Nickname: "39",
			Username: "aixieluo",
			Email:    "admin@aixieluo.com",
		}
	})
	It("Create User", func() {
		u, err := ro.Create(ctx, uD)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(u.Mobile).Should(Equal("17521066239"))
	})
	It("User List", func() {
		users, total, err := ro.List(ctx, 1, 10)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(users).Should(HaveLen(1))
		Ω(total).Should(Equal(1))
	})
	It("Get User", func() {
		u, err := ro.FindByID(ctx, uint(uD.ID))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(u.ID).Should(Equal(uD.ID))
	})
	It("Update User", func() {
		b, err := ro.Update(ctx, uint(uD.ID), &biz.User{
			Nickname: "3939",
		})
		Ω(err).ShouldNot(HaveOccurred())
		Ω(b).Should(BeTrue())
	})
})
