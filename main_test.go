package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	var routerJson OVNRawJson
	ReadOVNRawJson("testdata/logical_router.json", &routerJson)
	t.Run("Parse ovn-nbctled json file", func(t *testing.T) {
		var routers []LogicalRouter
		ParseNetworkDevice(&routers, routerJson)
	})
}
