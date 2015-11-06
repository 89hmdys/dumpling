package wechat

import (
	"dumpling/utils"
	"encoding/xml"
	"errors"
	"fmt"
	"noodles/http"
	"time"
)

type JsApiPayCommand struct {
	OrderNo string  //订单编号
	Subject string  //商品名称
	Price   float64 //支付价格
	OpenId  string  //用户OPENID
}

func NewForJSAPI(appId, appKey, mchId, notifyUrl string) Pay {
	return &jsapi{wechat: wechat{
		AppId:     appId,
		AppKey:    appKey,
		MchId:     mchId,
		NotifyUrl: notifyUrl}}
}

type jsapi struct {
	wechat
}

func (this *jsapi) Pay(v interface{}) (map[string]string, error) {
	command := v.(JsApiPayCommand)

	orderInfo := make(map[string]string)
	orderInfo["appid"] = this.AppId
	orderInfo["mch_id"] = this.MchId
	orderInfo["nonce_str"] = utils.UUID()
	orderInfo["out_trade_no"] = command.OrderNo
	orderInfo["body"] = command.Subject
	orderInfo["total_fee"] = fmt.Sprintf("%.0f", command.Price)
	orderInfo["spbill_create_ip"] = utils.LocalIP()
	orderInfo["notify_url"] = this.NotifyUrl
	orderInfo["trade_type"] = "JSAPI"
	orderInfo["openid"] = command.OpenId

	sign, err := this.sign(orderInfo)

	if err != nil {
		return nil, err
	}
	orderInfo["sign"] = sign

	orderInfoXml := this.toXml(orderInfo)

	header := map[string]string{
		"Content-Type": "text/xml;charset=UTF-8"}

	returnXml, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", header, orderInfoXml)

	if err != nil {
		return nil, err
	}

	return this.prepareOrder(returnXml)
}

func (this *jsapi) prepareOrder(src string) (map[string]string, error) {
	prepareOrder := &PrepareOrder{}
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
