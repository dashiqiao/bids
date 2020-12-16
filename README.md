# BIDS 报表

## 内容列表

- [背景](#背景)
- [安装](#安装)
- [使用说明](#使用说明)
- [数据库说明](#数据库说明)
- [维护者](#维护者)
- [捐赠](#捐赠)
- [使用许可](#使用许可)

## 背景 
> 构建低代码报表平台，可以快速建立各部门需要的数据报表，减少代码开发。

## 安装

## 使用说明

> 创建报表，填写SQL、选择需要的检索条件，完成配置。

> 图表与报表创建模式一致。

> 数据指标看板。

> 用户分群等。

----------------------------
## 数据库说明
### bids
#### 

|  表名   | 意义 
|  ----   | ----  
| report_defines  | 报表定义表 |
| report_conditions  | 条件表 |
| report_header  | 表头 |
| report_support  | 报表与条件关联表 |
| report_graph  | 图表定义表 |
| report_graph_relation  | 图表与报表关联表 |
| report_graph_support  | 图表与条件关联表 |
| report_target_defines  | 指标定义表 |
| report_board  | 看板表 |
| report_board_card  | 看板卡片表 |
| report_layout  | 看板布局表 |
| report_field_defines  | “动作”表 |
| report_field  | “动作”详情表 |
| report_field_relation  | “动作”和报表关联表 |
| report_analysis  | 用户分群表 |
| report_action_log  | “动作”日志表 |
| report_log  | “浏览”日志表 |

## 维护者

> 刘亮、宋善强

## 捐赠
> Buy Me A Coffee

## 使用许可
 采用 [Apache License, Version 2.0](https://github.com/denverdino/aliyungo/blob/master/LICENSE.txt)许可证授权原则。