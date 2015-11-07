package wechat

type Pay interface {
	Notify(s string) (*Receipt, error)
	Pay(order Order) (map[string]string, error)
}
