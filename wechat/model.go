package wechat

import (
	"dumpling/utils"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	formatDatetime string = "2006-01-02 15:04:05"
)

type notification struct {
	AppId         string `xml:"appid,emitempty"`
	BankType      string `xml:"bank_type,emitempty"`
	CashFee       string `xml:"cash_fee,emitempty"`
	FeeType       string `xml:"fee_type,emitempty"`
	IsSubscribe   string `xml:"is_subscribe,emitempty"`
	MchId         string `xml:"mch_id,emitempty"`
	NonceStr      string `xml:"nonce_str,emitempty"`
	OpenId        string `xml:"openid,emitempty"`
	OutTradeNo    string `xml:"out_trade_no,emitempty"`
	ReturnCode    string `xml:"return_code,emitempty"`
	ResultCode    string `xml:"result_code,emitempty"`
	Sign          string `xml:"sign,emitempty"`
	TimeEnd       string `xml:"time_end,emitempty"`
	TotalFee      string `xml:"total_fee,emitempty"`
	TradeType     string `xml:"trade_type,emitempty"`
	TransactionId string `xml:"transaction_id,emitempty"`
}

func (this notification) toMap() map[string]string {
	t := reflect.TypeOf(this).Elem()
	v := reflect.ValueOf(this).Elem()

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		tags := strings.Split(t.Field(i).Tag.Get("xml"), ",")
		data[string(tags[0])] = v.Field(i).String()
	}

	return data
}

type prepareOrder struct {
	ReturnCode string `xml:"return_code"`
	ResultCode string `xml:"result_code"`
	PrepayId   string `xml:"prepay_id,omitempty"`
}

type Receipt struct {
	OrderNo    string
	TradeNo    string
	TotalPrice string
}

type Order map[string]string

type OrderBuilder map[string]string

func (this OrderBuilder) SetDevice(device string) OrderBuilder {
	this["device_info"] = device
	return this
}

func (this OrderBuilder) SetDetail(detail string) OrderBuilder {
	this["detail"] = detail
	return this
}

func (this OrderBuilder) SetAttach(attach string) OrderBuilder {
	this["attach"] = attach
	return this
}

func (this OrderBuilder) SetFeeType(feeType string) OrderBuilder {
	this["fee_type"] = feeType
	return this
}

func (this OrderBuilder) SetStartTime(startTime time.Time) OrderBuilder {
	this["time_start"] = startTime.Format(formatDatetime)
	return this
}

func (this OrderBuilder) SetExpireTime(expireTime time.Time) OrderBuilder {
	this["time_expire"] = expireTime.Format(formatDatetime)
	return this
}

func (this OrderBuilder) SetLimitPay(rule string) OrderBuilder {
	this["limit_pay"] = rule
	return this
}

func (this OrderBuilder) BuildForApp(orderNo, description string, price int64) Order {

	this["nonce_str"] = utils.UUID()
	this["out_trade_no"] = orderNo
	this["body"] = description
	this["total_fee"] = fmt.Sprintf("%d", price)
	this["spbill_create_ip"] = utils.LocalIP()
	this["trade_type"] = "APP"

	return Order(this)
}

func (this OrderBuilder) BuildForSubscription(orderNo, description, openId string, price float64) Order {

	this["nonce_str"] = utils.UUID()
	this["out_trade_no"] = orderNo
	this["body"] = description
	this["total_fee"] = fmt.Sprintf("%d", price)
	this["spbill_create_ip"] = utils.LocalIP()
	this["trade_type"] = "JSAPI"
	this["openid"] = openId

	return Order(this)
}
