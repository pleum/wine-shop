package wine

import (
	"bufio"
	"errors"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Wine struct {
	OriginCountry string
	Vineyard      string
	WineType      string
	Winery        string
	Vintage       string
	Rating        float32
	UnitPrice     decimal.Decimal
	EntryTime     time.Time
}

var orderedColumn = []string{
	"Country of Origin",
	"Vineyard",
	"Winery",
	"Wine Type",
	"Vintage",
	"Rating",
	"Price",
}

func NewWineFromFile(file *os.File, entryTime time.Time) ([]Wine, error) {
	var wines []Wine
	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return wines, err
	}

	if len(lines) == 0 {
		return wines, nil
	}

	// -1 = to be check
	//  0 = field ordered
	//  1 = column named
	parsedMode := -1
	columnIdxToValueIdxMap := make(map[int]int)

	for lineIdx, line := range lines {
		values := strings.Split(line, ",")

		if parsedMode == -1 {
			for columnIdx, column := range orderedColumn {
				valueIdx := slices.Index(values, column)

				if valueIdx != -1 {
					// set to field ordered, if found at least 1 column named.
					if parsedMode == -1 {
						parsedMode = 1
					}

					columnIdxToValueIdxMap[columnIdx] = valueIdx
				}
			}

			// set to field ordered.
			if parsedMode == -1 {
				parsedMode = 0
			}
		}

		if parsedMode == 1 {
			// Skip header row.
			if lineIdx == 0 {
				continue
			}

			var orderedValues []string
			for columnIdx := range orderedColumn {
				orderedValues = append(orderedValues, values[columnIdxToValueIdxMap[columnIdx]])
			}

			values = orderedValues
		}

		w, err := createWineFromOrderedFieldValues(values, entryTime)
		if err != nil {
			return wines, err
		}

		wines = append(wines, w)
	}

	return wines, nil
}

func parseRating(rawRating string) (float32, error) {
	ratingParts := strings.Split(rawRating, ": ")
	if len(ratingParts) != 2 {
		return 0.0, errors.New("invalid rating input")
	}

	r, err := strconv.ParseFloat(ratingParts[1], 32)
	if err != nil {
		return 0.0, err
	}

	ratingScale := ratingParts[0]
	if ratingScale == "Parker" {
		return float32(r), nil
	}

	if ratingScale == "Robinson" {
		return (float32(r) / 20.0) * 100.0, nil
	}

	return 0.0, errors.New("invalid rating scale")
}

func createWineFromOrderedFieldValues(values []string, entryTime time.Time) (Wine, error) {
	if len(values) != 7 {
		return Wine{}, errors.New("there's some missing data")
	}

	var formattedValue []string
	for _, v := range values {
		formattedValue = append(formattedValue, strings.TrimSpace(v))
	}

	rating, err := parseRating(formattedValue[5])
	if err != nil {
		return Wine{}, err
	}

	unitPrice, err := decimal.NewFromString(formattedValue[6])
	if err != nil {
		return Wine{}, err
	}

	return Wine{
		OriginCountry: formattedValue[0],
		Vineyard:      formattedValue[1],
		Winery:        formattedValue[2],
		WineType:      formattedValue[3],
		Vintage:       formattedValue[4],
		Rating:        rating,
		UnitPrice:     unitPrice,
		EntryTime:     entryTime,
	}, nil
}
