package wechat

import (
	"bytes"
	"crypto/md5"
	. "dumpling/utils"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

type wechat struct {
	AppId     string //应用ID
	AppKey    string //API密钥，用于签名验签
	MchId     string //商户号
	NotifyUrl string //回调地址
}

func (this *wechat) Notify(notificationXml string) (*Receipt, error) {

	notification := notification{}

	err := xml.Unmarshal([]byte(notificationXml), &notification)
	if err != nil {
		return nil, err
	}

	notificationMap := notification.toMap()

	returnCode := notificationMap["return_code"]

	if returnCode != success {
		return nil, errors.New("pay fail")
	}
	sign := notificationMap["sign"]
	mySign, errSign := this.sign(notificationMap)
	if errSign != nil {
		return nil, errors.New("verify sign fail")
	}

	if strings.ToLower(mySign) == strings.ToLower(sign) {
		return &Receipt{OrderNo: notificationMap["out_trade_no"],
			TradeNo:    notificationMap["transaction_id"],
			TotalPrice: notificationMap["total_fee"]}, nil
	} else {
		return nil, errors.New("verify sign fail")
	}
}

func (this *wechat) sign(src map[string]string) (string, error) {
	ParaFilter(src)
	waitToSign := SortAndJoin(src)
	waitToSign = fmt.Sprintf("%s&key=%s", waitToSign, this.AppKey)

	hash := md5.New()
	hash.Write([]byte(waitToSign))
	sum := md5.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(sum[:])), nil
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
