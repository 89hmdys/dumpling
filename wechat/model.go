package wechat

import (
	"fmt"
	"github.com/89hmdys/dumpling/utils"
	"reflect"
	"strings"
	"time"
)

const (
	formatDatetime string = "20060102150405"
)

type Notification struct {
	AppId         string `xml:"appid,emitempty"`        //应用ID
	Attach        string `xml:"attach,emitempty"`       //自定义数据，支付时提交，原文返回
	BankType      string `xml:"bank_type,emitempty"`    //付款银行
	CashFee       string `xml:"cash_fee,emitempty"`     //现金支付金额
	FeeType       string `xml:"fee_type,emitempty"`     //货币类型
	IsSubscribe   string `xml:"is_subscribe,emitempty"` //是否关注微信号
	MchId         string `xml:"mch_id,emitempty"`       //商户ID
	NonceStr      string `xml:"nonce_str,emitempty"`    //随机数
	OpenId        string `xml:"openid,emitempty"`       //用户OPENID
	OutTradeNo    string `xml:"out_trade_no,emitempty"` //商户订单号
	ReturnCode    string `xml:"return_code,emitempty"`
	ResultCode    string `xml:"result_code,emitempty"`
	Sign          string `xml:"sign,emitempty"`
	TimeEnd       string `xml:"time_end,emitempty"`       //支付结束时间
	TotalFee      string `xml:"total_fee,emitempty"`      //支付总金额
	TradeType     string `xml:"trade_type,emitempty"`     //交易类型 JSAPI APP这些
	TransactionId string `xml:"transaction_id,emitempty"` //微信流水号
}

func (this Notification) ToMap() map[string]string {
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

type Order map[string]string

type OrderBuilder map[string]string

func (this OrderBuilder) hasSet(key string) {
	if _, ok := this[key]; ok {
		delete(this, key)
	}
}

func (this OrderBuilder) SetDevice(device string) OrderBuilder {
	key := "device_info"
	this.hasSet(key)
	this[key] = device
	return this
}

func (this OrderBuilder) SetDetail(detail string) OrderBuilder {
	key := "detail"
	this.hasSet(key)
	this[key] = detail
	return this
}

func (this OrderBuilder) SetAttach(attach string) OrderBuilder {
	key := "attach"
	this.hasSet(key)
	this[key] = attach
	return this
}

func (this OrderBuilder) SetFeeType(feeType string) OrderBuilder {
	key := "fee_type"
	this.hasSet(key)
	this[key] = feeType
	return this
}

func (this OrderBuilder) SetTime(startTime time.Time, duration int64) OrderBuilder {
	key := "time_start"
	this.hasSet(key)
	this[key] = startTime.Format(formatDatetime)
	key = "time_expire"
	this.hasSet(key)
	this[key] = startTime.Add(time.Duration(duration) * time.Minute)
	return this
}

func (this OrderBuilder) SetLimitPay(rule string) OrderBuilder {
	key := "limit_pay"
	this.hasSet(key)
	this[key] = rule
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
