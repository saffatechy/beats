// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package aws

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "aws", asset.ModuleFieldsPri, AssetAws); err != nil {
		panic(err)
	}
}

// AssetAws returns asset data.
// This is the base64 encoded gzipped contents of module/aws.
func AssetAws() string {
	return "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"
}
