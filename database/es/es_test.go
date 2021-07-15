package es

import (
	"context"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestCreatIndexAndAlias(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	type EM map[string]interface{}
	documentType := "_doc" // Es DocumentType
	index := "dev_test"
	aliasIndex := "dev_test_alias"
	data := EM{
		"dynamic": false,
		"properties": EM{
			"id": EM{
				"type": "integer",
			},
			"tag": EM{
				"type":            "text",
				"analyzer":        "ik_max_word",
				"search_analyzer": "ik_smart",
			},
			"description": EM{
				"type":            "text",
				"analyzer":        "ik_max_word",
				"search_analyzer": "ik_smart",
			},
			"assort": EM{
				"type": "keyword",
			},
			"isSearch": EM{
				"type": "integer",
			},
			"city": EM{
				"type": "keyword",
			},
		},
	}

	esConf, err := New([]string{"127.0.0.1"}, "", "")
	assert.NilError(t, err)

	t.Run("creat index and alias", func(t *testing.T) {
		err = esConf.createIndexWithAlias(ctx, index, aliasIndex)
		assert.NilError(t, err)
	})

	t.Run("put mapping", func(t *testing.T) {
		err = esConf.putMapping(ctx, index, documentType, data)
		assert.NilError(t, err)
	})

	t.Run("delete index", func(t *testing.T) {
		err = esConf.deleteIndex(ctx, index)
		assert.NilError(t, err)
	})
}
