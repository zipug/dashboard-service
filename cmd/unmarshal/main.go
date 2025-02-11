package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var content []map[string]string
	test := "[{\"project_id\": \"2\", \"question\": \"Настройка 1С\", \"telegram_id\": \"13312121\", \"article_name\": \"Настройка 1С по шагам\", \"date\": \"2025-02-11\"}, {\"project_id\": \"1\", \"question\": \"Как написать функцию на Python\", \"telegram_id\": \"13312142\", \"article_name\": \"Основы Python\", \"date\": \"2025-02-11\"}]"
	err := json.Unmarshal([]byte(test), &content)
	if err != nil {
		fmt.Println(err)
	}
	for i, item := range content {
		fmt.Printf("#%d ROW: {\n", i)
		for key, value := range item {
			fmt.Printf("\t%s: %s\n", key, value)
		}
		fmt.Println("}")
	}
}
