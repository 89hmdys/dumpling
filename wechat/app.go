package wechat

import (
	"dumpling/utils"
	"encoding/xml"
	"errors"
	"fmt"
	"noodles/http"
	"time"
)

func NewAppClient(appId, appKey, mchId, notifyUrl string) Pay {
	return &app{wechat: wechat{
		AppId:     appId,
		AppKey:    appKey,
		MchId:     mchId,
		NotifyUrl: notifyUrl}}
}

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
