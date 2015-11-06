package wechat

type Pay interface {
	Notify(v map[string]string) (*Receipt, error)
	Pay(command interface{}) (map[string]string, error)
}
