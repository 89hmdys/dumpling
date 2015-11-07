package wechat

import (
	"bytes"
	"crypto/md5"
	"dumpling/utils"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"noodles/http"
	"strings"
	"time"
)

/*
==================================父client，所有的支付方式都从这里派生=====================================
*/

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
	utils.ParaFilter(src)
	waitToSign := utils.SortAndJoin(src)
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

/*
==================================app支付client=====================================
*/

type app struct {
	wechat
}

func (this *app) Pay(order Order) (map[string]string, error) {

	order["appid"] = this.AppId
	order["mch_id"] = this.MchId
	order["notify_url"] = this.NotifyUrl

	sign, err := this.sign(order)

	if err != nil {
		return nil, err
	}
	order["sign"] = sign

	orderInfoXml := this.toXml(order)

	header := map[string]string{
		"Content-Type": "text/xml;charset=UTF-8"}

	returnXml, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", header, orderInfoXml)

	if err != nil {
		return nil, err
	}

	return this.prepareOrder(returnXml)
}

func (this *app) prepareOrder(src string) (map[string]string, error) {
	prepareOrder := &prepareOrder{}
	err := xml.Unmarshal([]byte(src), prepareOrder)
	if err != nil {
		return nil, err
	}

	if prepareOrder.ReturnCode == success {
		if prepareOrder.ResultCode == success {
			payInfo := make(map[string]string)
			payInfo["appid"] = this.AppId
			payInfo["partnerid"] = this.MchId
			payInfo["prepayid"] = prepareOrder.PrepayId
			payInfo["package"] = "Sign=WXPay"
			now := time.Now()
			payInfo["noncestr"] = utils.UUID()
			payInfo["timestamp"] = fmt.Sprintf("%d", now.UnixNano()/int64(time.Second))
			sign, err := this.sign(payInfo)
			if err != nil {
				return nil, err
			}
			payInfo["sign"] = sign
			return payInfo, nil
		} else {
			return nil, errors.New("支付失败,ResultCode : fail")
		}
	} else {
		return nil, errors.New("支付失败，,ReturnCode : fail")
	}
}

/*
==================================公众号支付client=====================================
*/

type subscription struct {
	wechat
}

func (this *subscription) Pay(order Order) (map[string]string, error) {

	order["appid"] = this.AppId
	order["mch_id"] = this.MchId
	order["notify_url"] = this.NotifyUrl

	sign, err := this.sign(order)

	if err != nil {
		return nil, err
	}
	order["sign"] = sign

	orderInfoXml := this.toXml(order)

	header := map[string]string{
		"Content-Type": "text/xml;charset=UTF-8"}

	returnXml, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", header, orderInfoXml)

	if err != nil {
		return nil, err
	}

	return this.prepareOrder(returnXml)
}

func (this *subscription) prepareOrder(src string) (map[string]string, error) {
	prepareOrder := &prepareOrder{}
	err := xml.Unmarshal([]byte(src), prepareOrder)
	if err != nil {
		return nil, err
	}

	if prepareOrder.ReturnCode == success {
		if prepareOrder.ResultCode == success {
			payInfo := make(map[string]string)
			payInfo["appId"] = this.AppId
			now := time.Now()
			payInfo["timeStamp"] = fmt.Sprintf("%d", now.UnixNano()/int64(time.Second))
			payInfo["package"] = fmt.Sprintf("prepay_id=%s", prepareOrder.PrepayId)
			payInfo["nonceStr"] = fmt.Sprintf("%d", int64(now.UnixNano()))
			payInfo["signType"] = "MD5"
			sign, err := this.sign(payInfo)
			if err != nil {
				return nil, err
			}
			payInfo["paySign"] = sign
			return payInfo, nil
		} else {
			return nil, errors.New("支付失败,ResultCode : fail")
		}
	} else {
		return nil, errors.New("支付失败，,ReturnCode : fail")
	}
}
