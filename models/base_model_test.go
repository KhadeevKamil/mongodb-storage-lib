package models_test

import (
	"time"

	timecop "github.com/bluele/go-timecop"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rekamarket/mongodb-storage-lib/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("BaseModel", func() {
	Context("with existing model", func() {
		var (
			bs       *BaseModel
			objectID primitive.ObjectID
		)

		BeforeEach(func() {
			objectID = primitive.NewObjectID()
			bs = &BaseModel{ID: objectID}
		})

		Describe("constructor", func() {
			It("should create new model", func() {
				Expect(bs).NotTo(BeNil())
				Expect(bs.ID).To(Equal(objectID))
			})
		})

		Describe("methods", func() {
			Describe("GetID()", func() {
				It("should return ObjectID", func() {
					Expect(bs.GetID()).To(Equal(objectID))
				})
			})

			Describe("GetHexID()", func() {
				It("should return Hex ID", func() {
					Expect(bs.GetHexID()).To(Equal(objectID.Hex()))
				})
			})

			Describe("SetHexID()", func() {
				It("should set Hex ID", func() {
					newObjectID := primitive.NewObjectID()
					err := bs.SetHexID(newObjectID.Hex())

					Expect(err).To(BeNil())
					Expect(bs.GetHexID()).NotTo(Equal(objectID.Hex()))
					Expect(bs.GetHexID()).To(Equal(newObjectID.Hex()))
				})
			})

			Describe("SetupTimestamps()", func() {
				Context("when created_at and updated_at are empty", func() {
					It("should set created_at and updated_at to current time", func() {
						timecop.Freeze(time.Now().Add(time.Hour * 1))
						defer timecop.Return()

						bs.SetupTimestamps()

						Expect(bs.CreatedAt).To(Equal(timecop.Now().Unix()))
						Expect(bs.UpdatedAt).To(Equal(timecop.Now().Unix()))
					})
				})

				Context("when created_at is not empty, updated_at is empty", func() {
					var (
						prevCreatedAt int64
					)

					BeforeEach(func() {
						prevCreatedAt = timecop.Now().Unix()
						bs.CreatedAt = prevCreatedAt
						bs.UpdatedAt = 0
					})

					It("should updated only updated_at field", func() {
						timecop.Freeze(time.Now().Add(time.Hour * 1))
						defer timecop.Return()

						bs.SetupTimestamps()

						Expect(bs.CreatedAt).To(Equal(prevCreatedAt))
						Expect(bs.UpdatedAt).To(Equal(timecop.Now().Unix()))
					})
				})
			})
		})
	})
})
