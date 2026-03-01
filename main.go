package main

import (
	"encoding/json"
	"fmt"
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

type OVNRawTable struct {
	Records [][]any  `json:"data"`
	Header  []string `json:"headings"`
}

func ParseNetworkDevice(routers *[]LogicalRouter, table OVNRawTable) {
	indexMap := make(map[string]int)
	headers := []string{"_uuid", "name", "ports"}

	for i, v := range table.Header {
		if slices.Contains(headers, v) {
			indexMap[v] = i
		}
	}

	for _, r := range table.Records {
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
	bytes, err := os.ReadFile("./logical_router.json")
	if err != nil {
		panic(err)
	}

	var lrt OVNRawTable
	err = json.Unmarshal(bytes, &lrt)
	if err != nil {
		panic(err)
	}

	var routers []LogicalRouter

	ParseNetworkDevice(&routers, lrt)

	fmt.Printf("%+v", routers)

}
