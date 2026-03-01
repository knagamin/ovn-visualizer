package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
)

type uuid struct {
	uuid string
}

type DeviceParser interface {
	ParseNetworkDevice()
}

type LogicalRouter struct {
	name  string
	uuid  string
	ports []string
}

type LogicalSwitch struct {
	name  string
	uuid  string
	ports []string
}

type OVNRawJson struct {
	Records [][]any  `json:"data"`
	Header  []string `json:"headings"`
}

func ReadOVNRawJson(path string, raw *OVNRawJson) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Can't read file: %v", err)
	}

	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		log.Fatalf("Can't unmarshal JSON data: %v", err)
	}

}

func ParseNetworkDevice(routers *[]LogicalRouter, devJson OVNRawJson) {
	indexMap := make(map[string]int)
	headers := []string{"_uuid", "name", "ports"}

	for i, v := range devJson.Header {
		if slices.Contains(headers, v) {
			indexMap[v] = i
		}
	}

	for _, r := range devJson.Records {
		var router LogicalRouter
		router.name = r[indexMap["name"]].(string)
		router.uuid = r[indexMap["_uuid"]].([]any)[1].(string)
		for _, port := range r[indexMap["ports"]].([]any)[1].([]any) {
			router.ports = append(router.ports, port.([]any)[1].(string))
		}
		*routers = append(*routers, router)
	}
}

func main() {

	var routerJson OVNRawJson

	ReadOVNRawJson("testdata/logical_router.json", &routerJson)

	var routers []LogicalRouter

	ParseNetworkDevice(&routers, routerJson)

	fmt.Printf("%+v", routers)

}
