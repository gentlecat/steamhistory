package main

import (
	"fmt"
	"log"
	"os"

	"code.google.com/p/goconf/conf"
	"github.com/rdwilliamson/aws"
	"github.com/rdwilliamson/aws/glacier"
)

func main() {
	c, err := conf.ReadConfigFile("aws.config")
	if err != nil {
		log.Fatal("Failed to load config file! ", err)
	}

	secret, err := c.GetString("default", "secret")
	if err != nil {
		log.Fatal("Failed to get AWS secret key from config! ", err)
	}
	access, err := c.GetString("default", "access")
	if err != nil {
		log.Fatal("Failed to get AWS access key from config! ", err)
	}

	// TODO: Allow server customization.
	conn := glacier.NewConnection(secret, access, aws.USEast)

	// Opening backup file
	name, err := c.GetString("glacier", "backup-file")
	if err != nil {
		log.Fatal(err)
	}
	archive, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	vault, err := c.GetString("glacier", "vault")
	if err != nil {
		log.Fatal(err)
	}
	id, err := conn.UploadArchive(vault, archive, "backup")
	if err != nil {
		log.Fatal("Failed to upload archive! ", err)
	}

	fmt.Printf("Archive uploaded. ID: %s", id)
}
