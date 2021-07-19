package template

import (
	"strings"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/def"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
)

func CreateToml() string {
	return `
[service]
name = "gen"

[api]
[api.structs]
GetPreOrderRequest = [
    "name;string;用户姓名",
    "dId;int64;用户dId",
]
CreatePreOrderRequest = [
    "name;string;用户姓名",
    "dId;int64;用户dId",
    "create_pre_order_content;struct:CreatePreOrderContent:详情"
]
CreatePreOrderContent = [
    "get_pre_order_content;*GetPreOrderContent;用户姓名",
    "list;[]*GetPreOrderContent;用户dId",
]
GetPreOrderContent = [
    "name;string;用户姓名",
    "dId;int64;用户dId",
]

[[api.handlers]]
name = "PublicOrder"
prefix = "/public/v1/order"
routes = [
    "createPreOrder;CreatePreOrderRequest;创建订单",
    "getPreOrder;GetPreOrderRequest;获取订单详情",
]



[[api.handlers]]
name = "PrivateOrder"
prefix = "/private/v1/order"
routes = [
    "createPrivatePreOrder;CreatePreOrderRequest;创建私有订单",
    "getPrivatePreOrder;GetPreOrderRequest;获取私有订单"
    ]
 
[cronjob]
StatisticDataByDay = "30 1 0 * * *"
LoopCSync = "30 1 0 * * *"

[job]
LoopOrder = "loopOrder"
LoopInvoice = "loopOrder"

[mysql]
[mysql.order]
order = [
    "column:id;primary_key;type:int AUTO_INCREMENT",
    "column:order_no;type:varchar(100) NOT NULL;default:'';comment:'订单号';unique_index",
    "column:uId;type:int NOT NULL;default:0;comment:'用户ID号';index",
    ]
goods = [
    "column:id;primary_key;type:int AUTO_INCREMENT",
    "column:order_no;type:varchar(100) NOT NULL;default:'';comment:'订单号';unique_index",
    "column:uId;type:int NOT NULL;default:0;comment:'用户ID号';index",
    "column:goods_id;type:varchar(50) NOT NULL;default:'';comment:'商品id';index",
    ]
[mysql.wallet]
record = [
    "column:id;primary_key;type:int AUTO_INCREMENT",
    "column:order_no;type:varchar(100) NOT NULL;default:'';comment:'订单号';unique_index",
    "column:uId;type:int NOT NULL;default:0;comment:'用户ID号';index",
    "column:goods_id;type:varchar(50) NOT NULL;default:'';comment:'商品id';index",
    "column:goods_num;type:int NOT NULL;default:'0';comment:'数量(sku属性)'",
    ]

[redis]
[redis.order]
`
}

func CreateField(field string) string {
	tags := strings.Split(field, ";")
	if len(tags) == 0 {
		return ""
	}

	fieldMap := make(map[string]string)
	for _, v1 := range tags {
		attributes := strings.Split(v1, ":")
		if len(attributes) < 2 {
			continue
		}
		fieldMap[attributes[0]] = attributes[1]
	}
	fieldName := fieldMap["column"]
	upFieldName := file.UnderlineToCamel(fieldName)
	fieldType := def.GetTypeName(fieldMap["type"])
	return upFieldName + "  " + fieldType + " `json:\"" + fieldName + "\" gorm:\"" + field + "\"`\n"
}
