package wechat

import (
	"bytes"
	. "dumpling/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"toast/hash/md5"
)

const success string = "SUCCESS"

type PrepareOrder struct {
	ReturnCode string `xml:"return_code"`
	ResultCode string `xml:"result_code"`
	PrepayId   string `xml:"prepay_id,omitempty"`
}

type Receipt struct {
	OrderNo    string
	TradeNo    string
	TotalPrice string
}

type wechat struct {
	AppId     string //应用ID
	AppKey    string //API密钥，用于签名验签
	MchId     string //商户号
	NotifyUrl string //回调地址
}

func (this *wechat) Notify(src map[string]string) (*Receipt, error) {
	returnCode := src["return_code"]

	if returnCode != success {
		return nil, errors.New("pay fail")
	}
	sign := src["sign"]
	mySign, errSign := this.sign(src)
	if errSign != nil {
		return nil, errors.New("verify sign fail")
	}

	if strings.ToLower(mySign) == strings.ToLower(sign) {
		return &Receipt{OrderNo: src["out_trade_no"], TradeNo: src["transaction_id"], TotalPrice: src["total_fee"]}, nil
	} else {
		return nil, errors.New("verify sign fail")
	}
}

func (this *wechat) sign(src map[string]string) (string, error) {
	ParaFilter(src)
	waitToSign := SortAndJoin(src)
	waitToSign = fmt.Sprintf("%s&key=%s", waitToSign, this.AppKey)

	sum := md5.Sum([]byte(waitToSign))

	return strings.ToUpper(hex.EncodeToString(sum)), nil
}

func (this *wechat) toXml(order map[string]string) string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString("<xml>")
	for k, v := range order {
		buffer.WriteString(fmt.Sprintf("<%s>", k))
		buffer.WriteString(fmt.Sprintf("<![CDATA[%s]]>", v))
		buffer.WriteString(fmt.Sprintf("</%s>", k))
	}
	buffer.WriteString("</xml>")
	return buffer.String()
}
