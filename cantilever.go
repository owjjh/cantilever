package main

import (
		"os"
		"fmt"
		"io/ioutil"
		"encoding/xml"
		"encoding/json"
		"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
// new comment
	type Dependency struct {
		GroupId string `xml:"groupId"`
		ArtifactId string `xml:"artifactId"`
		Version string `xml:"version"`
	}
	type ProjectXml struct {
		Dependencies []Dependency `xml:"dependencies>dependency"`
	}

	type ProjectJson struct {
		Dependencies interface{} `json:"dependencies"`
	}

	// collect all dependencies
	var allDepend []Dependency

	xmlFile, err := os.Open("pom.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	var pom ProjectXml
	xml.Unmarshal(b, &pom)

	for _, depend := range pom.Dependencies {
		allDepend = append(allDepend, depend)
		// fmt.Println(depend)
	}

	jsonFile, err := os.Open("package.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer jsonFile.Close()

	c, _ := ioutil.ReadAll(jsonFile)

	var pack ProjectJson

	json.Unmarshal(c, &pack)

	m := pack.Dependencies.(map[string]interface{})

	for k,v := range m {
		switch vv := v.(type) {
		case string:
			allDepend = append(allDepend, Dependency{"no group",k, vv})
			// fmt.Println(k, vv)
    	default:
    	}
	}


	// Load everything into dynamodb
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))

	for _,value := range allDepend {
		fmt.Println(value.GroupId, value.ArtifactId, value.Version)

		params := &dynamodb.PutItemInput {
			Item: map[string]*dynamodb.AttributeValue {
				"packageID": {S: aws.String(value.ArtifactId)},
				"version": {S: aws.String(value.Version) },
			},
			TableName: aws.String("owjjh-test"),
		}

		_, err := svc.PutItem(params)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// fmt.Println(resp)
	}

}
