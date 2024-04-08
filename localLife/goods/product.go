package goods

import (
	"encoding/json"
	"fmt"

	"github.com/RichXan/douyin-sdk/response"
	"github.com/RichXan/douyin-sdk/util"
	"github.com/google/go-querystring/query"
)

type GoodQuery struct {
	Cursor           string `json:"cursor" form:"cursor" url:"cursor,omitempty"`                                     // 第一页不传，之后用前一次返回的next_cursor传入进行翻页
	Count            int    `json:"count" form:"count" url:"count,omitempty"`                                        // 分页数量，不传默认为5，最大上限50
	Status           int    `json:"status" form:"status" url:"status,omitempty"`                                     // 过滤在线状态 1-在线 2-下线 3-封禁
	AccountID        string `json:"account_id" form:"account_id" url:"account_id,required"`                          // 商家ID，传入时服务商须与该商家满足授权关系
	GoodsCreatorType int    `json:"goods_creator_type" form:"goods_creator_type" url:"goods_creator_type,omitempty"` // 区分商品创建者的查询方式 0-查询服务商/开发者创建的商品（默认） 1-查询商家（account_id）创建的商品
	QueryAllPoi      bool   `json:"query_all_poi" form:"query_all_poi" url:"query_all_poi,omitempty"`                // 是否查询商品全量门店，仅开放给主营类目为餐饮的商家 设置为true时，分页数量最大上限为20
	GoodsQueryType   int    `json:"goods_query_type" form:"goods_query_type" url:"goods_query_type,omitempty"`       // 新商品查询参数，仅支持餐饮行业的ka自研开发者使用，当前参数生效时goods_creator_type参数不再生效 1-可查询归属于ka自研的所有商品，包括非ka自研创建的商品
}

func (good *Good) GetGoodList(param *GoodQuery) (*[]ProductOnline, error) {
	clientToken, err := good.GetClientToken()
	fmt.Println("clientToken:", clientToken)
	if err != nil {
		return nil, err
	}
	params, err := query.Values(param)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v?%v", getOnlineProductsURL, params.Encode())
	fmt.Println(url)
	header := map[string]string{
		"access-token": clientToken,
	}
	res, err := util.HTTPGet(url, header)
	if err != nil {
		return nil, err
	}
	// if res != nil {
	// 	err = fmt.Errorf(string(res))
	// 	return nil, err
	// }
	rep := response.Response{}
	rep.Data = ResProducts{}
	// fmt.Println(res)
	err = util.DecodeWithError(res, &rep, "GetGoodList")
	if err != nil {
		return nil, fmt.Errorf("decodeWithError is invalid %v", err)
	}
	// err = json.Unmarshal(res, &rep)
	// if err != nil {
	// 	return nil, fmt.Errorf("json Unmarshal Error, err=%v", err)
	// }
	// list := rep.Data.(GoodList)
	// fmt.Println(rep)
	repData, err := json.Marshal(rep.Data)
	// err = mapstructure.Decode(rep.Data, &goodList)
	if err != nil {
		return nil, fmt.Errorf("rep data encode valid %v", err)
	}
	var list ResProducts
	err = json.Unmarshal(repData, &list)
	if err != nil {
		return nil, fmt.Errorf("rep data decode valid %v", err)
	}
	return &list.Products, err
}

type Product struct {
	CategoryFullName string `json:"category_full_name"` // 品类全名，保存时不必填写
	CategoryID       string `json:"category_id"`        // 品类ID
	ProductType      int    `json:"product_type"`       // 商品类型 1-团购 11-代金券 15-次卡

}

// ResProducts 是返回的响应体结构体
type ResProducts struct {
	// NextCursor 用于查询下一页的游标
	NextCursor string `json:"next_cursor"`
	// HasMore 表示是否还有下一页数据
	HasMore bool `json:"has_more"`
	// Products 是线上数据列表
	Products []ProductOnline `json:"products"`
	// ErrorCode 表示错误码
	ErrorCode int `json:"error_code"`
	// Description 是错误描述
	Description string `json:"description"`
}

// ProductOnline 表示线上商品数据
type ProductOnline struct {
	// OnlineStatus 表示在线状态：1-在线，2-下线，3-封禁
	OnlineStatus int `json:"online_status"`
	// Product 是商品结构体
	Product ProductStruct `json:"product"`
	SKU     SKU           `json:"sku"`
}

// ProductStruct 表示商品结构体
type ProductStruct struct {
	// ProductID 是商品ID
	ProductID string `json:"product_id"`
	// OutID 是外部商品ID
	OutID string `json:"out_id"`
	// ProductName 是商品名
	ProductName string `json:"product_name"`
	// CategoryFullName 是品类全名
	CategoryFullName string `json:"category_full_name"`
	// CategoryID 是品类ID
	CategoryID int `json:"category_id"`
	// ProductType 是商品类型：1-团购，11-代金券，15-次卡
	ProductType int `json:"product_type"`
	// BizLine 是业务线：1-闭环自研开发者，3-直连服务商，5-小程序
	BizLine int `json:"biz_line"`
	// AccountName 是商家名
	AccountName string `json:"account_name"`
	// SoldStartTime 是售卖开始时间
	SoldStartTime int `json:"sold_start_time,omitempty"`
	// SoldEndTime 是售卖结束时间
	SoldEndTime int `json:"sold_end_time,omitempty"`
	// CreateTime 是创建时间
	CreateTime int `json:"create_time,omitempty"`
	// UpdateTime 是更新时间
	UpdateTime int `json:"update_time,omitempty"`
	// OutURL 是第三方跳转链接，小程序商品必填
	OutURL string `json:"out_url,omitempty"`
	// POIs 是店铺列表
	POIs []PoiStruct `json:"pois"`
	// AttrKeyValueMap 是商品属性键值对
	AttrKeyValueMap map[string]string `json:"attr_key_value_map"`
}

// PoiStruct 表示店铺结构体
type PoiStruct struct {
	// SupplierExtID 是接入方店铺ID
	SupplierExtID string `json:"supplier_ext_id"`
	// POIID 是POI ID
	POIID int `json:"poi_id,omitempty"`
	// SupplierID 是店铺ID
	SupplierID string `json:"supplier_id"`
}

// SKU 表示售卖单元
type SKU struct {
	// SkuID 是SKU ID
	SkuID string `json:"sku_id"`
	// SkuName 是SKU名
	SkuName string `json:"sku_name"`
	// OriginAmount 是原价，团购创建时可不填，会根据商品搭配计算原价
	OriginAmount int `json:"origin_amount"`
	// ActualAmount 是实际支付价格
	ActualAmount int `json:"actual_amount"`
	// Stock 是库存信息
	Stock StockStruct `json:"stock"`
	// OutSkuID 是第三方ID
	OutSkuID string `json:"out_sku_id"`
	// Status 是状态：1-在线，-1-删除
	Status int `json:"status"`
}

// StockStruct 表示库存信息结构体
type StockStruct struct {
	// LimitType 是库存上限类型：1-有限库存，2-无限库存
	LimitType int `json:"limit_type"`
	// StockQty 是总库存，当 LimitType 为 2 时无意义
	StockQty int `json:"stock_qty"`
	// AvailQty 是可用库存，当 LimitType 为 2 时无意义
	AvailQty int `json:"avail_qty"`
	// FrozenQty 是冻结库存，保存SKU时不填
	FrozenQty int `json:"frozen_qty"`
	// SoldQty 是售卖库存，退款回滚，保存SKU时不填
	SoldQty int `json:"sold_qty"`
	// SoldCount 是销量，退款不回滚，保存SKU时不填
	SoldCount int `json:"sold_count"`
}
