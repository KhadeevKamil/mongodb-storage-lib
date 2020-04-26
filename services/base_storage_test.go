package services_test

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rekamarket/mongodb-storage-lib/models"
	. "github.com/rekamarket/mongodb-storage-lib/services"

	timecop "github.com/bluele/go-timecop"
)

type ExampleModel struct {
	models.BaseModel `bson:"mb"`
	Title            string `json:"title,omitempty" bson:"title,omitempty"`
}

func NewExampleModel() *ExampleModel {
	return &ExampleModel{Title: "Test title " + uuid.New().String()}
}

// nolint
var _ = Describe("BaseStorage", func() {
	var (
		mongoURI, mongoDBName, mongoCollectionName string
		storage                                    *BaseStorage
		connErr                                    error
	)

	BeforeEach(func() {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://0.0.0.0:27017"
		}

		mongoDBName = "base_models_db_test"
		mongoCollectionName = "base_model_test"

		storage, connErr = NewBaseStorage(context.TODO(), mongoURI, mongoDBName, mongoCollectionName)

		if connErr != nil {
			GinkgoT().Fatal(connErr)
		}
	})

	Describe("consturctor", func() {
		It("should create a new storage", func() {
			Expect(connErr).To(BeNil())
			Expect(storage).NotTo(BeNil())
		})
	})

	Describe("methods", func() {
		AfterEach(func() {
			storage.DropAll(context.TODO())
		})

		Describe("save methods", func() {
			It("InsertOne", func() {
				curTime := time.Now().Add(time.Hour * 1)

				timecop.Freeze(curTime)
				defer timecop.Return()

				m := NewExampleModel()

				_, saveErr := storage.InsertOne(context.TODO(), m)

				Expect(saveErr).To(BeNil())
				Expect(m.CreatedAt).To(Equal(curTime.Unix()))
				Expect(m.UpdatedAt).To(Equal(curTime.Unix()))

			})
		})

	})

})
