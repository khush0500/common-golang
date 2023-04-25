package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)
var Data map[string]interface{}

func main() {
	yamlFile, err := ioutil.ReadFile("reqs.yaml")
	if err != nil {
		log.Panicf("%s", err)
	}

	
	err = yaml.Unmarshal(yamlFile, &Data)
	if err != nil {
		log.Fatalf("error unmarshalling YAML data: %v", err)
	}
	reqs := []string{}

	for key, value := range Data {
		if value == true {
			reqs = append(reqs, key)
		}
	}
	log.Println(reqs)
	args := []string{"/_newApp/"}
	err = generateBootstrap(args)
	if err != nil {
		log.Panicf("%s", err)
	}
	

}
