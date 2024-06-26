package main

import (
	"context"
	"os"

	gemeni "changer/pkg/gemini"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	key := os.Getenv("API_KEY")
	ctx := context.Background()

	data := `
	Can your write a concise profissional Change Log for these changes?
	Make it clear and concise.
	Don't include the code snippets.
	Dont add any extra information.
	Only seperate the changes with a new line and '-'.
	and include the functions and models that were changed in that line.
	
`
	// read old swagger file
	old,err := os.ReadFile("swagger.old.json")
	if err != nil {
		log.Fatal(err)
	}
	// read new swagger file
	new,err := os.ReadFile("swagger.new.json")
	if err != nil {
		log.Fatal(err)
	}
	
	data +=`old:
	` + string(old) +`
	
	new:
	` + string(new) +`
	`

	client, err := gemeni.NewGemini(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	r, err := client.Do(data)
	if err != nil {
		log.Fatal(err)
	}

	// write to file
	filename := "CHANGELLOG.md"
	err = writeToFile(filename, r)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Change log written to %s", filename)
}

func writeToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// check if its empty
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if stat.Size() > 0 {
		data = "\n" + data
	}

	// Write to the file
	_, err = file.WriteString(data)
	return err
}
