package cxtonsms

import "strings"

type Strong struct {
	Api
	Dest      []string // 手机号码, 建议不超过 3 万个
	Content   string   // 短信内容。最多 500 个字符。【签名】+内容
	Ext       string   // 扩展号码(视通道是否支持扩展，可以为空或不填)
	Reference string   // 参考信息(最多 50 个字符，在推送状态报告、推送上行时，推送给合作方，可以为空或不填)如果不知道如何使用，请忽略该参数，不能含有半角的逗号和分号。
	Delay     string   // 定时参数;格式: YYYYMMDDHHMISS, 可定时时间范围为一个月。可为空或不填
}

func (c *Strong) Map() (map[string]string, error) {
	maps := make(map[string]string, 0)
	// todo 判断是否大于3W
	maps["dest"] = strings.Join(c.Dest, ",")
	// todo 判断内容大于500
	maps["content"] = c.Content
	if c.Ext != "" {
		maps["ext"] = c.Ext
	}
	if c.Reference != "" {
		// todo 判断长度及是否有非标符号
		maps["reference"] = c.Reference
	}
	if c.Delay != "" {
		maps["delay"] = c.Delay
	}
	return maps, nil
}
