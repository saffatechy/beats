// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package nginx

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("filebeat", "nginx", asset.ModuleFieldsPri, AssetNginx); err != nil {
		panic(err)
	}
}

// AssetNginx returns asset data.
// This is the base64 encoded gzipped contents of module/nginx.
func AssetNginx() string {
	return "eJysl89u4zYQxu9+isGeN3oAHwoUCyywhxZFT70pDDmSp0tx1Bkqrd++oOzNOg4pU7J5ciTN9/3mTyjxCb7jcQ+hp/DfDiBS9LiHT7+nvz/tAByqFRojcdjDLzsAgN/YTR6hY4HRiFLoIR4Q5hDw3ENHHrXZAeiBJbaWQ0f9HqJMuAPoCL3T/Sz1BMEM+NM+rXgccQ+98DSer2QY0vo6C0EnPJQA0rr0u/Q01qLq2+Wc8YJ5Wl84RENBzxZzRX6CnPQTzxtKDucSSXDgiC2NrSeN7x75gWdEzPHqzgJiWr+GUxRwd3aAb3+AcU5QFbWBbxFIwUAyhRe0ZlIEmi9aHgYOEBkoWD85/AwvqORQ50ytJwzXoHAh//md1alXBzQORcHTd4Tnv56+svxrxKFLv56bD2p/ovGgPImdwUlBUCMLusT1fLrT0HgRmq3uC7tjqxhi83KMqPnyejLXd0YTD3s4xDg2gjpyUGySVlZmoF7MqRPneV9o80qEt0w3eE6K0qafKz1TXJOJq/EcMB7YbavzPxNqbLIKVemKX5uo+IaFegrmOrTGMGG3ryhKHLZknA+tm6bTTLaW3druvh9qjSZOmtOp4+hQBOWefhc0auxN/3Enqhnudg5c2/ry/1iZ4/rNAoW3waWkw1ey191YTi2bXlFnqbSXIIP5m6/bsoGjJFONQeEhGAWZWozRRHu4H6MkU4uR2Zc3UBRUaiE4Z7ISgbXpJu9zL6h1KO2D5pT17lFNMI+Z1gRz58Cylmq7nuXeJuV32i0sZaVV+3aPXPgW27JlWw6RAoZ4T8XPn3o9cnNTr7bwlqcQ5diScu4dvwntpmItnGc7P3Y/1IJSLYxgTxwe1L9lsermUTw+aqAWpFZW6HGjdFuwhPYDCUXebZEPPs7P8qtO85ZDQBvntPKnIM+hX3eW//KmCeQwROoI5ca51+Mrrj0Nee6bXFzNl/hYyLbsNgpbVG0+Rtb4xc1+8SBoXLPJdUBV0689cOWjrv3+DwAA///Lb/Q7"
}
