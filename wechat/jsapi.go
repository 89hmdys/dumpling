package wechat

import (
	"encoding/xml"
	"errors"
	"fmt"
	"noodles/http"
	"time"
)

func NewSubscriptionClient(appId, appKey, mchId, notifyUrl string) Pay {
	return &jsapi{wechat: wechat{
		AppId:     appId,
		AppKey:    appKey,
		MchId:     mchId,
		NotifyUrl: notifyUrl}}
}

type jsapi struct {
	wechat
}

func (this *jsapi) Pay(order Order) (map[string]string, error) {

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

func (this *jsapi) prepareOrder(src string) (map[string]string, error) {
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
