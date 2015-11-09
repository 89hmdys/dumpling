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
	key := "device_info"
	if _, ok := this[key]; ok {
		delete(this, key)
	}
	this[key] = device
	return this
}

func (this OrderBuilder) SetDetail(detail string) OrderBuilder {
	key := "detail"
	if _, ok := this[key]; ok {
		delete(this, key)
	}
	this[key] = detail
	return this
}

func (this OrderBuilder) SetAttach(attach string) OrderBuilder {
	key := "attach"
	if _, ok := this[key]; ok {
		delete(this, key)
	}
	this[key] = attach
	return this
}

func (this OrderBuilder) SetFeeType(feeType string) OrderBuilder {
	key := "fee_type"
	if _, ok := this[key]; ok {
		delete(this, key)
	}
	this[key] = feeType
	return this
}

func (this OrderBuilder) SetStartTime(startTime time.Time) OrderBuilder {
	key := "time_start"
	if _, ok := this[key]; ok {
		delete(this, key)
	}
	this[key] = startTime.Format(formatDatetime)
	return this
}

func (this OrderBuilder) SetExpireTime(expireTime time.Time) OrderBuilder {
	key := "time_expire"
	if _, ok := this[key]; ok {
		delete(this, key)
	}
	this[key] = expireTime.Format(formatDatetime)
	return this
}

func (this OrderBuilder) SetLimitPay(rule string) OrderBuilder {
	key := "limit_pay"
	if _, ok := this[key]; ok {
		delete(this, key)
	}
	this[key] = rule
	return this
}

func (this Order) BuildForApp(orderNo, description string, price int64) Order {

	this["nonce_str"] = utils.UUID()
	this["out_trade_no"] = orderNo
	this["body"] = description
	this["total_fee"] = fmt.Sprintf("%d", price)
	this["spbill_create_ip"] = utils.LocalIP()
	this["trade_type"] = "APP"

	return Order(this)
}

func (this Order) BuildForSubscription(orderNo, description, openId string, price float64) Order {

	this["nonce_str"] = utils.UUID()
	this["out_trade_no"] = orderNo
	this["body"] = description
	this["total_fee"] = fmt.Sprintf("%d", price)
	this["spbill_create_ip"] = utils.LocalIP()
	this["trade_type"] = "JSAPI"
	this["openid"] = openId

	return Order(this)
}
