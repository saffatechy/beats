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

package pgsql

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("packetbeat", "pgsql", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "eJzEk8GO2jAQhu95il976aXwADlUQisqrRS1LHBHJpkkVh07zNgg3r5yEkoKiVRO9QkN/v199kwW+EXXFG0lJ5MAXntDKd42TnzFtPvM3hKgIMlZt147m+JbAgD3DQtpKdelzkFnsh6lJlPIMsHwK+32L2BVQ3dQXP7aUoqKXWiHyjgxTp0C8fVPdVLotvY1gd2lB/VBKAGTKlCya+BrgmdlReUx/yX+dwokfpk8cbUQs+MRoHc+OmdI2X8z+ig75qbafWaDEJMPbAXKogN8ha+19NeHFgh5eAfPgSasusghdwXNCcRHuHeoDyAGlk9XMc5Wc4iGRFT1GmXIzGoLnYm1n+3m5KG30Fi/dSL6aOhwVibQaGh65nq7/bl9qH1f7VfZQ22z+vHx/ixrQ3Pop/GlJu/W2fp9P3RZl5CQ5yRSBjPX4poi60j813GuHL6FYVSomHjQ6Mju8t8MI3vs9zsAAP//6hY8Lw=="
}
