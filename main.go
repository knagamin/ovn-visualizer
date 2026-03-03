package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"slices"
)

var logicalRouters []LogicalRouter

type uuid struct {
	uuid string
}

type DeviceParser interface {
	ParseNetworkDevice()
}

type LogicalRouter struct {
	Name  string   `json:"name"`
	UUID  string   `json:"uuid"`
	Ports []string `json:"ports"`
}

type LogicalSwitch struct {
	Name  string   `json:"name"`
	UUID  string   `json:"uuid"`
	Ports []string `json:"ports"`
}

type OVNRawJson struct {
	Records [][]any  `json:"data"`
	Header  []string `json:"headings"`
}

func ReadOVNRawJson(path string, raw *OVNRawJson) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Can't read file")
	}

	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		slog.Error("Can't unmarshal JSON data")
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
		router.Name = r[indexMap["name"]].(string)
		router.UUID = r[indexMap["_uuid"]].([]any)[1].(string)
		for _, port := range r[indexMap["ports"]].([]any)[1].([]any) {
			router.Ports = append(router.Ports, port.([]any)[1].(string))
		}
		*routers = append(*routers, router)
	}
}

type api struct {
	addr string
}

func (a *api) getLogicalRoutersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(logicalRouters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var routerJson OVNRawJson

	slog.Info("Load OVN Topology info from JSON File")
	ReadOVNRawJson("testdata/logical_router.json", &routerJson)

	var routers []LogicalRouter

	ParseNetworkDevice(&routers, routerJson)

	slog.Info("successfully Loaded")

	logicalRouters = routers

	api := &api{addr: ":8080"}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /routers", api.getLogicalRoutersHandler)

	slog.Info("starting api server...", "addr", api.addr)

	srv.ListenAndServe()

}
