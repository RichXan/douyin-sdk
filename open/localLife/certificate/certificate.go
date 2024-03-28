package certificate

import (
	"encoding/json"
	"fmt"

	"github.com/RichXan/douyin-sdk/open/localLife/context"
	"github.com/RichXan/douyin-sdk/response"
	"github.com/RichXan/douyin-sdk/util"
	"github.com/google/go-querystring/query"
)

const (
	getCertificateURL = "https://open.douyin.com/goodlife/v1/fulfilment/certificate/get/"
)

type Certificate struct {
	*context.Context
}

type GetCertificateRequest struct {
	EncryptedData string `json:"encrypted_data,omitempty" form:"encrypted_data,omitempty" url:"encrypted_data,omitempty"`
}

type CanNoUseDate struct {
	EndTime   int64 `json:"end_time"`
	StartTime int64 `json:"start_time"`
}
type NotAvailableTimeInfo struct {
	FulfilEnable    bool           `json:"fulfil_enable"`
	CanNoUseWeekDay []int64        `json:"can_no_use_week_day"`
	CanNoUseDate    []CanNoUseDate `json:"can_no_use_date"`
}
type VerifyRecords struct {
	VerifyType         int    `json:"verify_type"`
	VerifierUniqueID   string `json:"verifier_unique_id"`
	PoiID              int64  `json:"poi_id"`
	TimesCardSerialNum int64  `json:"times_card_serial_num"`
	VerifyID           string `json:"verify_id"`
	CertificateID      string `json:"certificate_id"`
	VerifyTime         int64  `json:"verify_time"`
	CanCancel          bool   `json:"can_cancel"`
}
type Sku struct {
	SkuID         string `json:"sku_id"`
	Title         string `json:"title"`
	GrouponType   int    `json:"groupon_type"`
	MarketPrice   int64  `json:"market_price"`
	SoldStartTime int64  `json:"sold_start_time"`
	ThirdSkuID    string `json:"third_sku_id"`
	AccountID     string `json:"account_id"`
}
type Verify struct {
	CanCancel          bool   `json:"can_cancel"`
	VerifyType         int    `json:"verify_type"`
	VerifierUniqueID   string `json:"verifier_unique_id"`
	PoiID              int64  `json:"poi_id"`
	TimesCardSerialNum int64  `json:"times_card_serial_num"`
	VerifyID           string `json:"verify_id"`
	CertificateID      string `json:"certificate_id"`
	VerifyTime         int64  `json:"verify_time"`
}
type Amount struct {
	OriginalAmount         int64 `json:"original_amount"`
	PayAmount              int64 `json:"pay_amount"`
	MerchantTicketAmount   int64 `json:"merchant_ticket_amount,omitempty"`
	PaymentDiscountAmount  int64 `json:"payment_discount_amount,omitempty"`
	CouponPayAmount        int64 `json:"coupon_pay_amount"`
	ListMarketAmount       int64 `json:"list_market_amount"`
	PlatformDiscountAmount int64 `json:"platform_discount_amount,omitempty"`
}
type CertificateModel struct {
	Code                 string               `json:"code"`
	EncryptedCode        string               `json:"encrypted_code"`
	StartTime            int64                `json:"start_time"`
	CertificateID        int64                `json:"certificate_id"`
	NotAvailableTimeInfo NotAvailableTimeInfo `json:"not_available_time_info"`
	UsedStatusType       int                  `json:"used_status_type"`
	VerifyRecords        []VerifyRecords      `json:"verify_records"`
	Sku                  Sku                  `json:"sku"`
	Verify               Verify               `json:"verify"`
	TimeCard             TimeCard             `json:"time_card"`
	Status               int                  `json:"status"`
	ExpireTime           int64                `json:"expire_time"`
	Amount               Amount               `json:"amount"`
}

func NewCertificate(ctx *context.Context) *Certificate {
	return &Certificate{ctx}
}

func (certificate *Certificate) GetCertificate(in *GetCertificateRequest) (*CertificateModel, error) {
	accessToken, err := certificate.GetAccessToken()
	// fmt.Println(accessToken)
	if err != nil {
		return nil, err
	}
	params, err := query.Values(in)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v?%v", getCertificateURL, params.Encode())
	header := map[string]string{
		"access-token": accessToken,
	}
	res, err := util.HTTPGet(url, header)
	if err != nil {
		return nil, err
	}
	rep := response.Response{}
	rep.Data = CertificateModel{}
	err = util.DecodeWithError(res, &rep, "GetCertificate")
	if err != nil {
		return nil, fmt.Errorf("decodeWithError is invalid %v", err)
	}
	repData, err := json.Marshal(rep.Data)
	if err != nil {
		return nil, fmt.Errorf("rep data encode valid %v", err)
	}
	var certModel CertificateModel
	err = json.Unmarshal(repData, &certModel)
	if err != nil {
		return nil, fmt.Errorf("rep data decode valid %v", err)
	}
	return &certModel, err
}
