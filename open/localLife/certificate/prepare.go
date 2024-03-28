package certificate

import (
	"encoding/json"
	"fmt"

	"github.com/RichXan/douyin-sdk/response"
	"github.com/RichXan/douyin-sdk/util"
	"github.com/google/go-querystring/query"
)

const (
	certificatePrepareUrl = "https://open.douyin.com/goodlife/v1/fulfilment/certificate/prepare/"
)

type CertPrepareRequest struct {
	EncryptedData string `json:"encrypted_data,omitempty" form:"encrypted_data,omitempty" url:"encrypted_data,omitempty"`
	Code          string `json:"code,omitempty" form:"code,omitempty" url:"code,omitempty"`
	PoiId         string `json:"poi_id" form:"poi_id" url:"poi_id"`
}

type TimeCardAmount struct {
	OriginalAmount         int32 `json:"original_amount"`
	PayAmount              int32 `json:"pay_amount"`
	MerchantTicketAmount   int32 `json:"merchant_ticket_amount"`
	ListMarketAmount       int32 `json:"list_market_amount"`
	PlatformDiscountAmount int32 `json:"platform_discount_amount"`
	PaymentDiscountAmount  int32 `json:"payment_discount_amount"`
	CouponPayAmount        int32 `json:"coupon_pay_amount"`
}

type SerialAmountList struct {
	SerialNumb int32          `json:"serial_numb"`
	Amount     TimeCardAmount `json:"amount"`
}
type Certificates struct {
	StartTime     int64  `json:"start_time"`
	Amount        Amount `json:"amount"`
	CertificateID int64  `json:"certificate_id"`
	EncryptedCode string `json:"encrypted_code"`
	ExpireTime    int64  `json:"expire_time"`
	Sku           Sku    `json:"sku"`
}
type CertificatesV2 struct {
	StartTime     int64  `json:"start_time"`
	Amount        Amount `json:"amount"`
	CertificateID int64  `json:"certificate_id"`
	EncryptedCode string `json:"encrypted_code"`
	ExpireTime    int64  `json:"expire_time"`
	Sku           Sku    `json:"sku"`
}

type TimeCard struct {
	TimesCount       int32              `json:"times_count"`
	TimesUsed        int32              `json:"times_used"`
	SerialAmountList []SerialAmountList `json:"serial_amount_list"`
}
type PrepareData struct {
	Certificates   []Certificates   `json:"certificates"`
	CertificatesV2 []CertificatesV2 `json:"certificates_v2"`
	TimeCard       TimeCard         `json:"time_card"`
	OrderID        string           `json:"order_id"`
	VerifyToken    string           `json:"verify_token"`
	ErrorCode      int              `json:"error_code"`
	Description    string           `json:"description"`
}

func (certificate *Certificate) CertificatePrepare(in *CertPrepareRequest) (*PrepareData, error) {
	accessToken, err := certificate.GetAccessToken()
	// fmt.Println(accessToken)
	if err != nil {
		return nil, err
	}
	params, err := query.Values(in)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v?%v", certificatePrepareUrl, params.Encode())
	fmt.Println(url)
	header := map[string]string{
		"access-token": accessToken,
	}
	res, err := util.HTTPGet(url, header)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(res))
	rep := response.Response{}
	rep.Data = PrepareData{}
	// nrep := rep
	err = util.DecodeWithError(res, &rep, "CertificatePrepare")
	if err != nil {
		return nil, fmt.Errorf("decodeWithError is invalid %v", err)
	}
	repData, err := json.Marshal(rep.Data)
	if err != nil {
		return nil, fmt.Errorf("rep data encode valid %v", err)
	}
	var prepareData PrepareData
	err = json.Unmarshal(repData, &prepareData)
	if err != nil {
		return nil, fmt.Errorf("rep data decode valid %v", err)
	}
	// fmt.Println(rep)
	// if len(prepareData.Certificates) == 0 {

	// }
	return &prepareData, err
}
