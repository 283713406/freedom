# FREEDOM

FREEDOM 项目是使用go语言实现的对股票基本面进行分析的程序，通过命令行方式进行操作。数据来源于东方财富网、亿牛网、新浪财经、天天基金。

该程序不构成任何投资建议，程序只是个人辅助工具，具体分析仍然需要自己判断。


# FREEDOM 解决什么问题？

FREEDOM 要解决的问题是，在使用东方财富选股器按设置的条件筛选出股票后，通常筛选出的股票数量较多，需要人工对每一支股票的财报指标进行分析并判断是否可以长期持有。

需要分析的指标较多，并且有一些指标无法直接获取，需要进行计算甚至需要根据历史财务数据进行计算，在大量股票需要分析的时候这是一个繁琐的工作，因此开发了 FREEDOM 来让这个过程自动化。


# 功能

当前已实现的功能：

- 按指定条件的默认值自动筛选满足条件的公司
- 按指定条件的自定义值自动筛选满足条件的公司
- 实现股票检测器
- 支持 ROE、EPS、营收、利润、整体质地、估值、合理价、负债率、历史波动率、市值 检测
- 将筛选结果导出为 JSON 文件
- 将筛选结果导出为 CSV 文件
- 将筛选结果导出为 EXCEL 文件，并按行业、价格、历史波动率分工作表
- 支持关键词搜索股票并对其进行评估
- 检测器支持对银行股按不同规则进行检测
- 支持净利率和毛利率稳定性判断
- 支持 PEG 检测
- 支持营收本业比检测
- 支持财报审计意见检测
- 支持负债流动比检测
- 支持现金流检测


# 使用方法

git clone https://github.com/283713406/freedom.git

go build -o freedom main.go

查看使用说明：

```
$ ./freedom -h
NAME:
   freedom - mamba 的股票筛选和检测程序

USAGE:
   该程序不构成任何投资建议，程序只是个人辅助工具，具体分析仍然需要自己判断。

VERSION:
   v1

AUTHOR:
   mamba <283713406@qq.com>

COMMANDS:
   exportor  股票筛选导出器
   checker   股票检测器
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --loglevel value, -l value  cmd 日志级别 [debug|info|warn|error] (default: info) [$FREEDOM_CMD_LOGLEVEL]
   --help, -h                  show help (default: false)
   --version, -v               show the version (default: false)

COPYRIGHT:
   (c) 2021 mamba
```
   
## exportor

exportor 是数据导出器，不使用参数默认导出 EXCEL 文件。

查看使用说明：

```
./freedom exportor -h
NAME:
   freedom exportor - 股票筛选导出器

USAGE:
   将按条件筛选出的股票导出到文件，根据文件后缀名自动判断导出类型。支持的后缀名：[xlsx|csv|json|all]，all 表示导出全部支持的类型。

OPTIONS:
   --filename value, -f value                      指定导出文件名 (default: ./result/freedom-20211110.xlsx) [$FREEDOM_EXPORTOR_FILENAME]
   --disable_check, -C                             关闭基本面检测，导出所有原始筛选结果 (default: false) [$FREEDOM_EXPORTOR_DISABLE_CHECK]
   --filter.min_roe value                          最低净资产收益率 (%) (default: 1.0)
   --filter.min_netprofit_yoy_ratio value          最低净利润增长率 (%) (default: 0.0)
   --filter.min_toi_yoy_ratio value                最低营收增长率 (%) (default: 0.0)
   --filter.min_zxgxl value                        最低最新股息率 (%) (default: 0.5)
   --filter.min_netprofit_growthrate_3_y value     最低净利润 3 年复合增长率（%） (default: 0.0)
   --filter.min_income_growthrate_3_y value        最低营收 3 年复合增长率（%） (default: 0.0)
   --filter.min_listing_yield_year value           最低上市以来年化收益率（%） (default: 0.0)
   --filter.min_pb_new_mrq value                   最低市净率 (default: 0.1)
   --filter.max_debt_asset_ratio value             最大资产负债率 (%) (default: 95.0)
   --filter.min_predict_netprofit_ratio value      最低预测净利润同比增长（%） (default: 0.0)
   --filter.min_predict_income_ratio value         最低预测营收同比增长（%） (default: 0.0)
   --filter.min_total_market_cap value             最低总市值（亿） (default: 50.0)
   --filter.industry_list value                    行业名
   --filter.min_price value                        股价范围最小值（元） (default: 0.0)
   --filter.max_price value                        股价范围最大值（元） (default: 0.0)
   --filter.listing_over_5_y                       上市时间是否超过 5 年 (default: false)
   --filter.min_listing_volatility_year value      最低上市以来年化波动率 (default: 0.0)
   --filter.exclude_cyb                            排除创业板 (default: true)
   --filter.exclude_kcb                            排除科创板 (default: true)
   --filter.special_security_name_abbr_list value  查询指定名称
   --filter.special_security_code_list value       查询指定代码
   --filter.min_roa value                          最小总资产收益率 ROA (default: 0.0)
   --checker.min_roe value                         最新一期 ROE 不低于该值 (default: 10)
   --checker.check_years value                     连续增长年数 (default: 5)
   --checker.no_check_years_roe value              ROE 高于该值时不做连续增长检查 (default: 15)
   --checker.max_debt_asset_ratio value            最大资产负债率百分比(%) (default: 60)
   --checker.max_hv value                          最大历史波动率 (default: 0.99)
   --checker.min_total_market_cap value            最小市值（亿） (default: 50)
   --checker.bank_min_roa value                    银行股最小 ROA (default: 0.5)
   --checker.bank_min_zbczl value                  银行股最小资本充足率 (default: 8)
   --checker.bank_max_bldkl value                  银行股最大不良贷款率 (default: 3)
   --checker.bank_min_bldkbbfgl value              银行股最低不良贷款拨备覆盖率 (default: 100)
   --checker.is_check_mll_stability                是否检测毛利率稳定性 (default: false)
   --checker.is_check_jll_stability                是否检测净利率稳定性 (default: false)
   --checker.is_check_price_by_calc                是否使用估算合理价进行检测，高于估算价将被过滤 (default: false)
   --checker.max_peg value                         最大 PEG (default: 2.5)
   --checker.min_byys_ratio value                  最小本业营收比 (default: 0.5)
   --checker.max_byys_ratio value                  最大本业营收比 (default: 1.5)
   --checker.min_fzldb value                       最小负债流动比 (default: 1)
   --checker.is_check_cashflow                     是否检测现金流量 (default: false)
   --checker.is_check_mll_grow                     是否检测毛利率逐年递增 (default: false)
   --checker.is_check_jll_grow                     是否检测净利率逐年递增 (default: false)
   --checker.is_check_eps_grow                     是否检测EPS逐年递增 (default: false)
   --checker.is_check_rev_grow                     是否检测营收逐年递增 (default: false)
   --checker.is_check_netprofit_grow               是否检测净利润逐年递增 (default: false)
   --checker.min_gxl value                         最低股息率 (default: 0.5)
   --help, -h                                      show help (default: false)
```

命令行使用示例：

- 导出 JSON 文件：
```
./freedom -l error exportor -f ./stocks.json
```

- 导出 CSV 文件：
```
./freedom -l error exportor -f ./stocks.csv
```

- 导出 EXCEL 文件：
```
./freedom -l error exportor -f ./stocks.xlsx
```

- 导出全部支持的类型：
```
./freedom -l error exportor -f ./stocks.all
```

- 自定义筛选、检测参数
```
./freedom -l error exportor -f ./stocks.xlsx --filter.min_roe=6 --checker.min_roe=6
```

## checker

给定关键词/股票代码搜索股票进行评估检测

查看使用说明：

```
./freedom checker -h
NAME:
   freedom checker - 股票检测器

USAGE:
   freedom checker [command options] [arguments...]

OPTIONS:
   --keyword value, -k value             检给定股票名称或代码，多个股票批量检测使用/分割。如: 兴业银行/分众传媒/601166
   --checker.min_roe value               最新一期 ROE 不低于该值 (default: 10)
   --checker.check_years value           连续增长年数 (default: 5)
   --checker.no_check_years_roe value    ROE 高于该值时不做连续增长检查 (default: 15)
   --checker.max_debt_asset_ratio value  最大资产负债率百分比(%) (default: 60)
   --checker.max_hv value                最大历史波动率 (default: 0.99)
   --checker.min_total_market_cap value  最小市值（亿） (default: 50)
   --checker.bank_min_roa value          银行股最小 ROA (default: 0.5)
   --checker.bank_min_zbczl value        银行股最小资本充足率 (default: 8)
   --checker.bank_max_bldkl value        银行股最大不良贷款率 (default: 3)
   --checker.bank_min_bldkbbfgl value    银行股最低不良贷款拨备覆盖率 (default: 100)
   --checker.is_check_mll_stability      是否检测毛利率稳定性 (default: false)
   --checker.is_check_jll_stability      是否检测净利率稳定性 (default: false)
   --checker.is_check_price_by_calc      是否使用估算合理价进行检测，高于估算价将被过滤 (default: false)
   --checker.max_peg value               最大 PEG (default: 2.5)
   --checker.min_byys_ratio value        最小本业营收比 (default: 0.5)
   --checker.max_byys_ratio value        最大本业营收比 (default: 1.5)
   --checker.min_fzldb value             最小负债流动比 (default: 1)
   --checker.is_check_cashflow           是否检测现金流量 (default: false)
   --checker.is_check_mll_grow           是否检测毛利率逐年递增 (default: false)
   --checker.is_check_jll_grow           是否检测净利率逐年递增 (default: false)
   --checker.is_check_eps_grow           是否检测EPS逐年递增 (default: false)
   --checker.is_check_rev_grow           是否检测营收逐年递增 (default: false)
   --checker.is_check_netprofit_grow     是否检测净利润逐年递增 (default: false)
   --checker.min_gxl value               最低股息率 (default: 0.5)
   --help, -h                            show help (default: false)
```

命令行使用示例：

```
./freedom -l error checker -k 兴业银行
```
