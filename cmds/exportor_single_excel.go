// 导出结果为 excel

package cmds

import (
	"context"
	"errors"
	"time"

	"github.com/283713406/freedom/logging"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (

	// HeaderStyle 表头样式
	Header = &excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"FFCCCC"},
			Shading: 0,
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal:      "center",
			JustifyLastLine: true,
			Vertical:        "center",
			WrapText:        true,
		},
	}
	// BodyStyle 表格Style
	Body = &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "left",
			JustifyLastLine: true,
			Vertical:        "center",
			WrapText:        true,
		},
	}
)

// ExportExcel 导出 excel
func (e ExportorRNg) ExportRNgExcel(ctx context.Context, filename string) (result []byte, err error) {
	stocksCount := len(e.Stocks)
	if stocksCount == 0 {
		err = errors.New("no stocks data")
		return
	}
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
	headersLen := len(headers)
	headerStyle, err := f.NewStyle(Header)
	if err != nil {
		logging.Error(ctx, "New HeaderStyle error:"+err.Error())
	}
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

		// 设置表头样式
		hcell, err := excelize.CoordinatesToCellName(1, 1)
		if err != nil {
			logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
			continue
		}
		vcell, err := excelize.CoordinatesToCellName(headersLen, 1)
		if err != nil {
			logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
			continue
		}
		f.SetCellStyle(sheet, hcell, vcell, headerStyle)

		// 设置表格样式
		hcell, err = excelize.CoordinatesToCellName(1, 2)
		if err != nil {
			logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
			continue
		}
		vcell, err = excelize.CoordinatesToCellName(headersLen, stocksCount+3)
		if err != nil {
			logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
			continue
		}
		f.SetCellStyle(sheet, hcell, vcell, bodyStyle)
	}

	for _, sheet := range sheets {
		// 写 header
		for i, header := range headers {
			axis, err := excelize.CoordinatesToCellName(i+1, 1)
			if err != nil {
				logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
				continue
			}
			f.SetCellValue(sheet, axis, header)
		}
	}

	// 写 body
	line := []string{"毛利率", "三项费用率", "销售费用率", "管理费用率", "财务费用率", "净利润率", "资产负债率",
		"固定资产比重", "净资产收益率", "总资产周转率", "经营性现金流净额比净利润", "营业收入增长率", "扣非净利润增长率"}
	for _, sheet := range sheets {
		row := 2
		for _, stock := range e.Stocks {
			headerValueMap := stock.GetHeaderValueMap()
			logging.Infof(ctx, "headerValueMap is 111111111111111111111111111111111 %s", headerValueMap)
			for k, line := range line {
				col := k + 1
				axis, err := excelize.CoordinatesToCellName(col, row)
				if err != nil {
					logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
					continue
				}
				value := headerValueMap[line]
				f.SetCellValue(sheet, axis, value)
			}
			row++
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