package certificate

import (
	"encoding/json"
	"fmt"

	"github.com/RichXan/douyin-sdk/response"
	"github.com/RichXan/douyin-sdk/util"
)

const (
	certificateCancelUrl = "https://open.douyin.com/goodlife/v1/fulfilment/certificate/cancel/"
)

type CertCancelRequest struct {
	// 代表券码一次核销的唯一标识（验券时返回）(次卡撤销多次时请填0)
	VerifyID string `json:"verify_id" form:"verify_id" url:"verify_id"`
	// 代表一张券码的标识（验券时返回）
	CertificateID string `json:"certificate_id" form:"certificate_id" url:"certificate_id"`
	// 取消核销总次数（多次卡商品可传，优先级低于verify_id） 注意：如果是分门店结算，此字段不要传！！！
	CancelToken string `json:"cancel_token,omitempty" form:"cancel_token,omitempty" url:"cancel_token,omitempty"`
	// 撤销核销幂等操作，主要针对次卡，避免因超时等原因在短时间内重复请求导致撤销多次（幂等有效期1小时） 注意：如果是分门店结算，此字段不要传！！！
	TimesCardCancelCount int `json:"times_card_cancel_count" form:"times_card_cancel_count,omitempty" url:"times_card_cancel_count,omitempty"`
}

type CancelData struct {
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type CancelExtra struct {
	ErrorCode      int    `json:"error_code"`
	Description    string `json:"description"`
	SubErrorCode   int    `json:"sub_error_code"`
	SubDescription string `json:"sub_description"`
	Logid          string `json:"logid"`
	Now            int    `json:"now"`
}

func (certificate *Certificate) CertificateCancel(in *CertCancelRequest) (*CancelData, error) {
	clientToken, err := certificate.GetClientToken()
	if err != nil {
		return nil, err
	}
	header := map[string]string{
		"content-type": "application/json",
		"access-token": clientToken,
	}
	res, err := util.PostJSON(certificateCancelUrl, in, header)
	if err != nil {
		return nil, err
	}
	rep := response.Response{}
	rep.Data = CancelData{}
	err = util.DecodeWithError(res, &rep, "CertificateCancel")
	if err != nil {
		return nil, fmt.Errorf("decodeWithError is invalid %v", err)
	}
	repData, err := json.Marshal(rep.Data)
	if err != nil {
		return nil, fmt.Errorf("rep data encode valid %v", err)
	}
	var cancelData CancelData
	err = json.Unmarshal(repData, &cancelData)
	if err != nil {
		return nil, fmt.Errorf("rep data decode valid %v", err)
	}
	return &cancelData, err
}
