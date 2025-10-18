package service

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-mockingcode/models"
)

type DataGenerator struct {
	fake *gofakeit.Faker
}

func NewDataGenerator(seed uint64) *DataGenerator {
	faker := gofakeit.New(seed)
	return &DataGenerator{
		fake: faker,
	}
}

// GenerateDocuments генерирует массив документов по шаблону
func (g *DataGenerator) GenerateDocuments(fields []models.FieldTemplate, count int) []map[string]interface{} {
	documents := make([]map[string]interface{}, count)

	for i := 0; i < count; i++ {
		documents[i] = g.GenerateDocument(fields)
	}

	return documents
}

// GenerateDocument генерирует один документ по шаблону полей
func (g *DataGenerator) GenerateDocument(fields []models.FieldTemplate) map[string]interface{} {
	doc := make(map[string]interface{})

	for _, field := range fields {
		doc[field.Name] = g.generateField(field)
	}

	return doc
}

func (g *DataGenerator) generateField(field models.FieldTemplate) interface{} {
	switch field.Type {
	case "string":
		return g.generateString(field)
	case "number":
		return g.generateNumber(field)
	case "boolean":
		return g.generateBoolean(field)
	case "date":
		return g.generateDateTime(field)
	default:
		return g.generateString(field)
	}
}

func (g *DataGenerator) generateString(field models.FieldTemplate) string {
	switch field.Format {
	case "email":
		return g.fake.Email()
	case "phone":
		return g.fake.Phone()
	case "name":
		return g.fake.Name()
	case "url":
		return g.fake.URL()
	case "username":
		return g.fake.Username()
	case "address":
		return g.fake.Address().Address
	case "city":
		return g.fake.Address().City
	case "country":
		return g.fake.Address().Country
	case "uuid":
		return g.fake.UUID()
	default:
		if len(field.Options) > 0 {
			return g.fake.RandomString(field.Options)
		}
		return g.fake.Word()
	}
}

func (g *DataGenerator) generateNumber(field models.FieldTemplate) float64 {
	min := 0.0
	max := 100.0

	if field.Min != nil {
		min = *field.Min
	}
	if field.Max != nil {
		max = *field.Max
	}

	return g.fake.Float64Range(min, max)
}

func (g *DataGenerator) generateDateTime(_ models.FieldTemplate) time.Time {
	return g.fake.Date()
}

func (g *DataGenerator) generateBoolean(_ models.FieldTemplate) bool {
	return g.fake.Bool()
}
