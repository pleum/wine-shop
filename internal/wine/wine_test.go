package wine_test

import (
	"fmt"
	"os"
	"testing"
	"time"
	"wineshop/internal/wine"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const TestDataPath = "./testdata"

func TestNewWineFromFile(t *testing.T) {
	assert := assert.New(t)

	type expected struct {
		wines []wine.Wine
	}

	testCases := []struct {
		name     string
		expected expected
	}{
		{
			name: "data_without_column_name",
			expected: expected{
				wines: []wine.Wine{
					{
						OriginCountry: "Chile",
						Vineyard:      "Valle Central",
						Winery:        "Baron De Rothschild",
						WineType:      "Chardonnay",
						Vintage:       "2019",
						Rating:        94.0,
						UnitPrice:     decimal.NewFromFloat32(1290.0),
					},
					{
						OriginCountry: "Germany",
						Vineyard:      "Ihringer Winklerberg",
						Winery:        "Dr. Heger",
						WineType:      "Pinot Noir",
						Vintage:       "2014",
						Rating:        95.0,
						UnitPrice:     decimal.NewFromFloat32(2390.0),
					},
				},
			},
		},
		{
			name: "data_with_column_name",
			expected: expected{
				wines: []wine.Wine{
					{
						OriginCountry: "Chile",
						Vineyard:      "Valle Central",
						Winery:        "Baron De Rothschild",
						WineType:      "Chardonnay",
						Vintage:       "2019",
						Rating:        94.0,
						UnitPrice:     decimal.NewFromFloat32(1290.0),
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		now := time.Now()

		t.Run(tC.name, func(t *testing.T) {
			filePath := fmt.Sprintf("./%s/%s.csv", TestDataPath, tC.name)

			file, err := os.Open(filePath)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			wines, err := wine.NewWineFromFile(file, now)
			assert.NoError(err)

			assert.Len(wines, len(tC.expected.wines))

			for idx, expectedWine := range tC.expected.wines {
				assert.Equal(expectedWine.OriginCountry, wines[idx].OriginCountry)
				assert.Equal(expectedWine.Vineyard, wines[idx].Vineyard)
				assert.Equal(expectedWine.Winery, wines[idx].Winery)
				assert.Equal(expectedWine.WineType, wines[idx].WineType)
				assert.Equal(expectedWine.Rating, wines[idx].Rating)
				assert.True(expectedWine.UnitPrice.Equal(wines[idx].UnitPrice))
				assert.Equal(now, wines[idx].EntryTime)
			}
		})
	}
}
