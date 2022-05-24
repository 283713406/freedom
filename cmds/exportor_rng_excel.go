// 导出结果为 excel

package cmds

import (
	"context"
	"time"

	"github.com/283713406/freedom/logging"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	// BodyStyle 表格Style
	Body = &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "center",
			JustifyLastLine: true,
			Vertical:        "center",
			WrapText:        true,
		},
	}
)

// ExportExcel 导出 excel
func (e ExportorRNg) ExportRNgExcel(ctx context.Context, filename string) (result []byte, err error) {
	f := excelize.NewFile()

	// 创建全部数据表
	sheets := []string{}
	// 添加股票名
	Sheet1Name := ""
	for i, name := range e.Stocks.GetName() {
		sheets = append(sheets, name)
		if i == 0 {
			Sheet1Name = name
		}
	}

	headers := []string{"观察项目", "具体数据"}
	bodyStyle, err := f.NewStyle(Body)
	if err != nil {
		logging.Error(ctx, "New BodyStyle error:"+err.Error())
	}
	// 创建 sheet
	for _, sheet := range sheets {
		if sheet == Sheet1Name {
			f.SetSheetName("Sheet1", Sheet1Name)
		}
		f.NewSheet(sheet)
		width := 25.0
		f.SetColWidth(sheet, "A", "F", width)
		// 设置表头行高
		rowNum := 1
		height := 25.0
		f.SetRowHeight(sheet, rowNum, height)

		// 设置表格样式
		hcell, err := excelize.CoordinatesToCellName(1, 1)
		if err != nil {
			logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
			continue
		}
		vcell, err := excelize.CoordinatesToCellName(6, 15)
		if err != nil {
			logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
			continue
		}
		f.SetCellStyle(sheet, hcell, vcell, bodyStyle)
	}

	years := []string{"2021", "2020", "2019", "2018", "2017"}
	for _, sheet := range sheets {
		f.MergeCell(sheet, "B1", "F1")
		f.MergeCell(sheet, "A1", "A2")
		// 写 列头
		for i, header := range headers {
			axis, err := excelize.CoordinatesToCellName(i+1, 1)
			if err != nil {
				logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
				continue
			}
			f.SetCellValue(sheet, axis, header)
		}

		for i, year := range years {
			axis, err := excelize.CoordinatesToCellName(i+2, 2)
			if err != nil {
				logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
				continue
			}
			f.SetCellValue(sheet, axis, year)
		}
	}

	// 写 body
	lines := []string{"毛利率", "三项费用率", "销售费用率", "管理费用率", "财务费用率", "净利润率", "资产负债率",
		"固定资产比重", "净资产收益率", "总资产周转率", "经营性现金流净额比净利润", "营业收入增长率", "扣非净利润增长率"}
	for _, sheet := range sheets {
		// 写 行头
		for i, line := range lines {
			axis, err := excelize.CoordinatesToCellName(1, i+3)
			if err != nil {
				logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
				continue
			}
			f.SetCellValue(sheet, axis, line)
		}
	}

	for _, sheet := range sheets {
		col := 2
		logging.Debugf(ctx,"sheet is %s", sheet)
		for _, stock := range e.Stocks {
			if stock.Name != sheet {
				continue
			}
			logging.Debugf(ctx,"Stocks name is %s", stock.Name)
			headerValueMap := stock.GetHeaderValueMap()
			for k, line := range lines {
				row := k + 3
				values := headerValueMap[line]
				for i, value := range values.([]float64)  {
					axis, err := excelize.CoordinatesToCellName(i+col, row)
					if err != nil {
						logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
						continue
					}
					f.SetCellValue(sheet, axis, value)
				}
			}
		}
	}
	f.SetDocProps(&excelize.DocProperties{
		Created:     time.Now().Format("2006-01-02 15:04:05"),
		Creator:     "mamba",
		Keywords:    "freedom: https://github.com/283713406/freedom",
	})

	buf, err := f.WriteToBuffer()
	result = buf.Bytes()
	err = f.SaveAs(filename)
	return
}
