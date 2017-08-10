package template

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/go-ucfg/yaml"
)

func TestField(t *testing.T) {
	esVersion2, err := common.NewVersion("2.0.0")
	assert.NoError(t, err)

	falseVar := false
	trueVar := true

	tests := []struct {
		field  Field
		method func(f Field) common.MapStr
		output common.MapStr
	}{
		{
			field:  Field{Type: "long"},
			method: func(f Field) common.MapStr { return f.other() },
			output: common.MapStr{
				"type": "long",
			},
		},
		{
			field:  Field{Type: "scaled_float"},
			method: func(f Field) common.MapStr { return f.scaledFloat() },
			output: common.MapStr{
				"type":           "scaled_float",
				"scaling_factor": 1000,
			},
		},
		{
			field:  Field{Type: "scaled_float", ScalingFactor: 100},
			method: func(f Field) common.MapStr { return f.scaledFloat() },
			output: common.MapStr{
				"type":           "scaled_float",
				"scaling_factor": 100,
			},
		},
		{
			field:  Field{Type: "scaled_float", esVersion: *esVersion2},
			method: func(f Field) common.MapStr { return f.scaledFloat() },
			output: common.MapStr{
				"type": "float",
			},
		},
		{
			field:  Field{Type: "object", Enabled: &falseVar},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"type":    "object",
				"enabled": false,
			},
		},
		{
			field:  Field{Type: "object", Enabled: &falseVar},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"type":    "object",
				"enabled": false,
			},
		},
		{
			field:  Field{Type: "text", Analyzer: "autocomplete"},
			method: func(f Field) common.MapStr { return f.text() },
			output: common.MapStr{
				"type":     "text",
				"analyzer": "autocomplete",
				"norms":    false,
			},
		},
		{
			field:  Field{Type: "text", Analyzer: "autocomplete", Norms: true},
			method: func(f Field) common.MapStr { return f.text() },
			output: common.MapStr{
				"type":     "text",
				"analyzer": "autocomplete",
			},
		},
		{
			field:  Field{Type: "text", SearchAnalyzer: "standard", Norms: true},
			method: func(f Field) common.MapStr { return f.text() },
			output: common.MapStr{
				"type":            "text",
				"search_analyzer": "standard",
			},
		},
		{
			field:  Field{Type: "text", Analyzer: "autocomplete", SearchAnalyzer: "standard", Norms: true},
			method: func(f Field) common.MapStr { return f.text() },
			output: common.MapStr{
				"type":            "text",
				"analyzer":        "autocomplete",
				"search_analyzer": "standard",
			},
		},
		{
			field:  Field{Type: "text", MultiFields: Fields{Field{Name: "raw", Type: "keyword"}}, Norms: true},
			method: func(f Field) common.MapStr { return f.text() },
			output: common.MapStr{
				"type": "text",
				"fields": common.MapStr{
					"raw": common.MapStr{
						"type":         "keyword",
						"ignore_above": 1024,
					},
				},
			},
		},
		{
			field: Field{Type: "text", MultiFields: Fields{
				Field{Name: "raw", Type: "keyword"},
				Field{Name: "indexed", Type: "text"},
			}, Norms: true},
			method: func(f Field) common.MapStr { return f.text() },
			output: common.MapStr{
				"type": "text",
				"fields": common.MapStr{
					"raw": common.MapStr{
						"type":         "keyword",
						"ignore_above": 1024,
					},
					"indexed": common.MapStr{
						"type":  "text",
						"norms": false,
					},
				},
			},
		},
		{
			field:  Field{Dynamic: dynamicType{false}},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"dynamic": false, "type": "object",
			},
		},
		{
			field:  Field{Dynamic: dynamicType{true}},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"dynamic": true, "type": "object",
			},
		},
		{
			field:  Field{Dynamic: dynamicType{"strict"}},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"dynamic": "strict", "type": "object",
			},
		},
		{
			field:  Field{Type: "long", Index: &falseVar},
			method: func(f Field) common.MapStr { return f.other() },
			output: common.MapStr{
				"type": "long", "index": false,
			},
		},
		{
			field:  Field{Type: "text", Index: &trueVar},
			method: func(f Field) common.MapStr { return f.other() },
			output: common.MapStr{
				"type": "text", "index": true,
			},
		},
	}

	for _, test := range tests {
		output := test.method(test.field)
		assert.Equal(t, test.output, output)
	}
}

func TestDynamicTemplate(t *testing.T) {
	tests := []struct {
		field  Field
		method func(f Field) common.MapStr
		output common.MapStr
	}{
		{
			field: Field{
				Type: "object", ObjectType: "keyword",
				Name: "context",
			},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"context": common.MapStr{
					"mapping":            common.MapStr{"type": "keyword"},
					"match_mapping_type": "string",
					"path_match":         "context.*",
				},
			},
		},
		{
			field: Field{
				Type: "object", ObjectType: "long",
				path: "language", Name: "english",
			},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"language.english": common.MapStr{
					"mapping":            common.MapStr{"type": "long"},
					"match_mapping_type": "long",
					"path_match":         "language.english.*",
				},
			},
		},
		{
			field: Field{
				Type: "object", ObjectType: "text",
				path: "language", Name: "english",
			},
			method: func(f Field) common.MapStr { return f.object() },
			output: common.MapStr{
				"language.english": common.MapStr{
					"mapping":            common.MapStr{"type": "text"},
					"match_mapping_type": "string",
					"path_match":         "language.english.*",
				},
			},
		},
	}

	for _, test := range tests {
		dynamicTemplates = nil
		test.method(test.field)
		assert.Equal(t, test.output, dynamicTemplates[0])
	}
}

func TestDynamicYaml(t *testing.T) {
	tests := []struct {
		input  []byte
		output Field
		error  bool
	}{
		{
			input: []byte(`
name: test
dynamic: true
`),
			output: Field{
				Name:    "test",
				Dynamic: dynamicType{true},
			},
		},
		{
			input: []byte(`
name: test
dynamic: "true"
`),
			output: Field{
				Name:    "test",
				Dynamic: dynamicType{true},
			},
		},
		{
			input: []byte(`
name: test
dynamic: "blue"
`),
			error: true,
		},
		{
			input: []byte(`
name: test
dynamic: "strict"
`),
			output: Field{
				Name:    "test",
				Dynamic: dynamicType{"strict"},
			},
		},
	}

	for _, test := range tests {
		keys := Field{}

		cfg, err := yaml.NewConfig(test.input)
		assert.NoError(t, err)
		err = cfg.Unpack(&keys)

		if err != nil {
			assert.True(t, test.error)
		} else {
			assert.Equal(t, test.output.Dynamic, keys.Dynamic)
		}
	}
}
