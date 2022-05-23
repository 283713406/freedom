// 导出各类型的数据结果

package cmds

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/283713406/freedom/core"
	"github.com/283713406/freedom/logging"
	"github.com/283713406/freedom/models"
)

// Exportor exportor 实例
type Exportor struct {
	Stocks   models.ExportorDataList
	Selector core.Selector
}

// New 创建要导出的数据列表
func New(ctx context.Context, stocks models.StockList, selector core.Selector) Exportor {
	dlist := models.ExportorDataList{}
	for _, s := range stocks {
		dlist = append(dlist, models.NewExportorData(ctx, s))
	}

	return Exportor{
		Stocks:   dlist,
		Selector: selector,
	}
}

// ExportorRNg exportor 实例
type ExportorRNg struct {
	Stocks   models.ExportorRNgDataList
	Selector core.Selector
}

// NewRNg 创建要导出的数据列表
func NewRNg(ctx context.Context, stocks map[string]models.Stock, selector core.Selector) ExportorRNg {
	dlist := models.ExportorRNgDataList{}
	for _, s := range stocks {
		dlist = append(dlist, models.NewExportorRNgData(ctx, s))
	}

	return ExportorRNg{
		Stocks:   dlist,
		Selector: selector,
	}
}

// Export 导出数据
func Export(ctx context.Context, exportFilename string, selector core.Selector) {
	beginTime := time.Now()
	filedir := path.Dir(exportFilename)
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		os.Mkdir(filedir, 0755)
	}

	logging.Infof(ctx, "freedom exportor start export selected stocks to %s", exportFilename)
	var err error
	// 自动筛选股票
	stocks, err := selector.AutoFilterStocks(ctx)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	e := New(ctx, stocks, selector)

	_, err = e.ExportExcel(ctx, exportFilename)

	if err != nil {
		logging.Fatal(ctx, err.Error())
	}

	fmt.Printf(
		"\nfreedom exportor export succuss, total:%d latency:%#vs\n",
		len(stocks),
		time.Now().Sub(beginTime).Seconds(),
	)
}

// Export 导出数据
func ExportRNg(ctx context.Context, keywords []string, selector core.Selector) {
	exportFilename := fmt.Sprintf("./result/freedom-RNg-%s.xlsx", time.Now().Format("20060102"))
	beginTime := time.Now()
	filedir := path.Dir(exportFilename)
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		os.Mkdir(filedir, 0755)
	}

	logging.Infof(ctx, "freedom exportor start export selected stocks to %s", exportFilename)
	var err error
	// 自动筛选股票
	searcher := core.NewSearcher(ctx)
	stocks, err := searcher.SearchStocks(ctx, keywords)

	e := NewRNg(ctx, stocks, selector)

	_, err = e.ExportRNgExcel(ctx, exportFilename)

	if err != nil {
		logging.Fatal(ctx, err.Error())
	}

	fmt.Printf(
		"\nfreedom exportor export succuss, latency:%#vs\n",
		time.Now().Sub(beginTime).Seconds(),
	)
}
