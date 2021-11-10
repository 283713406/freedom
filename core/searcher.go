// 关键词搜索股票

package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/283713406/freedom/datacenter"
	"github.com/283713406/freedom/datacenter/eastmoney"
	"github.com/283713406/freedom/datacenter/sina"
	"github.com/283713406/freedom/logging"
	"github.com/283713406/freedom/models"
	"sync"
)

// Searcher 搜索器实例
type Searcher struct{}

// NewSearcher 创建搜索器实例
func NewSearcher(ctx context.Context) Searcher {
	return Searcher{}
}

// SearchStocks 按股票名或代码搜索股票
func (s Searcher) SearchStocks(ctx context.Context, keywords []string) (map[string]models.Stock, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	kLen := len(keywords)
	if kLen == 0 {
		return nil, errors.New("empty keywords")
	}
	// 根据关键词匹配股票代码
	matchedResults := []sina.SearchResult{}
	for _, kw := range keywords {
		wg.Add(1)
		go func(kw string) {
			defer func() {
				wg.Done()
			}()
			searchResults, err := datacenter.Sina.KeywordSearch(ctx, kw)
			if err != nil {
				logging.Errorf(ctx, "search %s error:%s", kw, err.Error())
				return
			}
			if len(searchResults) == 0 {
				logging.Warnf(ctx, "search %s no data", kw)
				return
			}
			logging.Infof(ctx, "search keyword:%s results:%+v, %+v matched", kw, searchResults, searchResults[0])
			mu.Lock()
			matchedResults = append(matchedResults, searchResults[0])
			mu.Unlock()
		}(kw)
	}
	wg.Wait()
	if len(matchedResults) == 0 {
		return nil, fmt.Errorf("无法获取对应数据 %v", keywords)
	}
	// 查询匹配到的股票代码的股票信息
	filter := eastmoney.Filter{}
	for _, result := range matchedResults {
		filter.SpecialSecurityCodeList = append(filter.SpecialSecurityCodeList, result.SecurityCode)
	}
	stocks, err := datacenter.EastMoney.QuerySelectedStocksWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	results := map[string]models.Stock{}
	for _, stock := range stocks {
		wg.Add(1)
		go func(stock eastmoney.StockInfo) {
			defer func() {
				wg.Done()
			}()
			mstock, err := models.NewStock(ctx, stock)
			if err != nil {
				logging.Errorf(ctx, "%s new models stock error:%v", stock.SecurityCode, err.Error())
				return
			}
			mu.Lock()
			results[stock.SecurityCode] = mstock
			mu.Unlock()
		}(stock)
	}
	wg.Wait()
	return results, nil
}