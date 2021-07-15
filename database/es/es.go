package es

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
)

type IndexEsData struct {
	Index        string
	AliasIndex   string
	DocumentType string
	Mapping      map[string]interface{}
}

type EM map[string]interface{}

type ConfSet struct {
	EsCli *elastic.Client
}

// New  es client
func New(addr []string, esUser, esPwd string) (*ConfSet, error) {
	cli, err := elastic.NewClient(
		elastic.SetURL(addr...),
		elastic.SetBasicAuth(esUser, esPwd),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}
	if cli == nil {
		return nil, fmt.Errorf("EsClient is nil")
	}
	return &ConfSet{EsCli: cli}, nil
}

// 创建索引 和 索引别名
func (c *ConfSet) createIndexWithAlias(ctx context.Context, index string, aliasIndex ...string) error {
	exist, err := c.EsCli.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	// 索引不存在时创建
	_, err = c.EsCli.CreateIndex(index).Do(ctx)
	if err != nil {
		return err
	}
	if len(aliasIndex) > 0 && len(aliasIndex[0]) > 0 {
		_, err = c.EsCli.Alias().Add(index, aliasIndex[0]).Do(ctx)
		if err != nil {
			c.EsCli.DeleteIndex(index).Do(ctx)
			return err
		}
	}

	return nil
}

// 删除索引
func (c *ConfSet) deleteIndex(ctx context.Context, index string) error {
	_, err := c.EsCli.DeleteIndex(index).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

// 创建映射
func (c *ConfSet) putMapping(ctx context.Context, index, DocumentType string, data map[string]interface{}) error {
	_, err := c.EsCli.PutMapping().Index(index).Type(DocumentType).BodyJson(data).AllowNoIndices(true).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfSet) InitEsIndex(indexList []IndexEsData) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	if len(indexList) == 0 {
		return fmt.Errorf("Es index is empty ")
	}

	for _, ei := range indexList {
		exist, err := c.EsCli.IndexExists(ei.Index).Do(ctx)
		if exist {
			// 如果存在 插入映射自动更新新增字段
			err = c.putMapping(ctx, ei.Index, ei.DocumentType, ei.Mapping)
			if err != nil {
				return err
			}
		} else {
			err = c.createIndexWithAlias(ctx, ei.Index, ei.AliasIndex)
			if err != nil {
				return err
			}

			err = c.putMapping(ctx, ei.Index, ei.DocumentType, ei.Mapping)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
