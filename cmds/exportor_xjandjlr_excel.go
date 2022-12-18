// 导出结果为 excel

package cmds

import (
	"context"
	"time"

	"github.com/283713406/freedom/logging"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// ExportExcel 导出 excel
func (e ExportorXJAndJLR) ExportXJAndJLRExcel(ctx context.Context, filename string) (result []byte, err error) {
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

	headers := []string{"年份", "净利润", "经营性现金流净额"}
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
		vcell, err := excelize.CoordinatesToCellName(3, 10)
		if err != nil {
			logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
			continue
		}
		f.SetCellStyle(sheet, hcell, vcell, bodyStyle)
	}

	for _, sheet := range sheets {
		// 写 行头
		for i, header := range headers {
			axis, err := excelize.CoordinatesToCellName(i+1, 1)
			if err != nil {
				logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
				continue
			}
			f.SetCellValue(sheet, axis, header)
		}
	}

	years := []string{"2017", "2018", "2019", "2020", "2021"}
	for _, sheet := range sheets {
		// 写 列头
		for i, year := range years {
			axis, err := excelize.CoordinatesToCellName(1, i+2)
			if err != nil {
				logging.Error(ctx, "CoordinatesToCellName error:"+err.Error())
				continue
			}
			f.SetCellValue(sheet, axis, year)
		}
	}

	// 写 body
	for _, sheet := range sheets {
		col := 2
		logging.Debugf(ctx,"sheet is %s", sheet)
		for _, stock := range e.Stocks {
			if stock.Name != sheet {
				continue
			}
			logging.Debugf(ctx,"Stocks name is %s", stock.Name)
			headerValueMap := stock.GetHeaderValueMap()
			for k, line := range headers {
				row := k + 2
				values := headerValueMap[line]
				for i, value := range values.([]string)  {
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
