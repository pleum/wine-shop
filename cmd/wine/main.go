package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"wineshop/internal/wine"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("file path is required")
		return
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	entryTime := time.Now()
	wines, err := wine.NewWineFromFile(file, entryTime)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(wines, "", "  ")
	if err != nil {
		panic(fmt.Errorf("error marshalling to JSON: %s", err.Error()))
	}

	outFile := fmt.Sprintf("./out/wine-%s.json", entryTime.Format(time.DateTime))
	if err = os.WriteFile(outFile, jsonData, 0644); err != nil {
		panic(fmt.Errorf("error writing to file: %s", err.Error()))
	}
}
