package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	OwnerSlug string `yaml:"owner-slug"`
	Source    string `yaml:"source"`
	Output    string `yaml:"output"`
}

type Feed struct {
	Owner    string    `json:"owner_slug"`
	Products []Product `json:"products"`
}

type XmlResponse struct {
	Document Document
}

type Document struct {
	Products []Product `xml:"DocDetail" json:"products"`
}

type Product struct {
	Ean      string  `xml:"EAN" json:"ean"`
	Sku      string  `xml:"SenderPrdCode1" json:"sku"`
	Vendor   string  `xml:"Brand" json:"vendor"`
	Category string  `xml:"Category" json:"category"`
	Price    float32 `xml:"Price8" json:"price"`
	Name     string  `xml:"ProductName" json:"name"`
	Stock    float32 `xml:"Quantity" json:"stock"`
}

func main() {
	cfg := Config{}
	err := cleanenv.ReadConfig("config.yml", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	xmlFile, err := os.Open(cfg.Source)
	defer xmlFile.Close()

	fmt.Println("Parse", cfg.Source)

	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var xmlResponse XmlResponse

	xml.Unmarshal(byteValue, &xmlResponse)

	feed := Feed{
		Owner:    cfg.OwnerSlug,
		Products: xmlResponse.Document.Products,
	}

	file, err := json.MarshalIndent(feed, "", " ")

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(cfg.Output, file, 0644)

	fmt.Println("Saved in", cfg.Output)

	if err != nil {
		log.Fatal(err)
	}
}
