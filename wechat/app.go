package wechat

import (
	"dumpling/utils"
	"encoding/xml"
	"errors"
	"fmt"
	"noodles/http"
	"time"
)

//const tradeType string = "APP"
//const minInterval time.Duration = 5 * time.Minute //最短间隔5分钟
type AppPayCommand struct {
	OrderNo string  //订单编号
	Subject string  //商品名称
	Price   float64 //支付价格
}

func NewForAPP(appId, appKey, mchId, notifyUrl string) Pay {
	return &app{wechat: wechat{
		AppId:     appId,
		AppKey:    appKey,
		MchId:     mchId,
		NotifyUrl: notifyUrl}}
}

type app struct {
	wechat
}

func (this *app) Pay(v interface{}) (map[string]string, error) {
	command := v.(AppPayCommand)

	orderInfo := make(map[string]string)
	orderInfo["appid"] = this.AppId
	orderInfo["mch_id"] = this.MchId
	orderInfo["nonce_str"] = utils.UUID()
	orderInfo["out_trade_no"] = command.OrderNo
	orderInfo["body"] = command.Subject
	orderInfo["total_fee"] = fmt.Sprintf("%.0f", command.Price)
	orderInfo["spbill_create_ip"] = utils.LocalIP()
	orderInfo["notify_url"] = this.NotifyUrl
	orderInfo["trade_type"] = "APP"

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

func (this *app) prepareOrder(src string) (map[string]string, error) {
	prepareOrder := &PrepareOrder{}
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
