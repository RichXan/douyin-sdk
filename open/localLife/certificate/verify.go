package certificate

import (
	"encoding/json"
	"fmt"

	"github.com/RichXan/douyin-sdk/response"
	"github.com/RichXan/douyin-sdk/util"
)

const (
	certificateVerifyUrl = "https://open.douyin.com/goodlife/v1/fulfilment/certificate/verify/"
)

type CertVerifyRequest struct {
	VerifyToken    string   `json:"verify_token" form:"verify_token"`
	EncryptedCodes []string `json:"encrypted_codes" form:"encrypted_codes[]"`
	PoiID          string   `json:"poi_id" form:"poi_id"`
}

type VerifyResults struct {
	VerifyID      string `json:"verify_id"`
	AccountID     string `json:"account_id"`
	CertificateID string `json:"certificate_id"`
	Code          string `json:"code"`
	Msg           string `json:"msg"`
	OrderID       string `json:"order_id"`
	OriginCode    string `json:"origin_code"`
	Result        int    `json:"result"`
}

type VerifyResultData struct {
	VerifyResults []VerifyResults `json:"verify_results"`
	ErrorCode     int             `json:"error_code"`
	Description   string          `json:"description"`
}

func (certificate *Certificate) CertificateVerify(in *CertVerifyRequest) (*VerifyResultData, error) {
	accessToken, err := certificate.GetAccessToken()
	// fmt.Println(accessToken)
	if err != nil {
		return nil, err
	}
	header := map[string]string{
		"access-token": accessToken,
	}
	res, err := util.PostJSON(certificateVerifyUrl, in, header)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(res))
	rep := response.Response{}
	rep.Data = VerifyResultData{}
	err = util.DecodeWithError(res, &rep, "CertVerify")
	if err != nil {
		return nil, fmt.Errorf("decodeWithError is invalid %v", err)
	}
	repData, err := json.Marshal(rep.Data)
	if err != nil {
		return nil, fmt.Errorf("rep data encode valid %v", err)
	}
	var verifyData VerifyResultData
	err = json.Unmarshal(repData, &verifyData)
	if err != nil {
		return nil, fmt.Errorf("rep data decode valid %v", err)
	}
	return &verifyData, err
}
