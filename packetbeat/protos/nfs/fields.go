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

package nfs

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("packetbeat", "nfs", asset.ModuleFieldsPri, AssetNfs); err != nil {
		panic(err)
	}
}

// AssetNfs returns asset data.
// This is the base64 encoded gzipped contents of protos/nfs.
func AssetNfs() string {
	return "eJykjzFP8zAQhnf/iledv+Yb6JSBBdSNFKliRsa5VCecO8t2IuXfo1RJG5RSBjxZ7/l53vMWnzSUkCYZIHP2VGJT7Y8bA9SUXOSQWaVEtT+i3/1/QArkuGEH6kkyGiZfp8JgupUGALYQ29LsHU8eApU4Re3ClCzfL5meYmKVSz6zXuW0CFfbhahZnfpZAOnaD4qFWTW0LBrf/9xz1vzelu1dXb/D0+Hl9fBWPUMDRTuORuiGSoPTmn6yPS4GOLuvvpH/B41oLcsqZ4GziaDNZZdvLme9Tzf2SdnmLt373bUpUvDDRBTGmK8AAAD//3Ajtg4="
}
