package alipay

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"net/url"
)

var (
	appID         = "9021000140600989"
	privateKey    = "MIIEpAIBAAKCAQEAruInhU3A3tibCtZQLNeVKkTSo66w8vazPS/PySrNVFNv26phNuli9IIYx06/EW6YvXYh/1VNHo9Cv1jA/7jKmGIY6gM8klDTgU3241kpfGy8sAwAYCs/cdd9laF5ofpw0Q9dAiAFNAEb9c+EGx60VyXQg8D6ClPtg7rlginXLGuUsfH3DamF8XI2BqQFS2WLSrLYMkVmCT5t1AuL1yf/uGdTaYdGi+q4GHq9orIOw/ShfynHfNfILmbJ0AIr/jKpKAvJwmRssVV9Yfxp54EgQAD8pk+RAIW43VlhV8xw45d2THP7VwtlaaaoKjfaxZ1mqwW2/pqqjS3K6FMseEGk1wIDAQABAoIBAG8k9f8HamN2gBrEF9JX9MoXUVOLq5PObB2f7EuqczJ7kKSnxP70FtrUb9EDX/VBN7t5f6PZ35hjbgVT79zV4ZQ1DCZ1hiJrUfBkz1qwCOi15rlv9zEnazv7uynEpRvnoZmMTQ0TsprZdZ6kkw31VqHoO3vwwjicHGBTAJfX5ZOpgf1b4ugFYRJD0cdGrEH7mqByZQJBEnSt4rg+UslhvKhxVafMh+act9jRcStT7vebeqDfQ5zQAGl1UNiryDnLNMWH+QLOBcR04ETH/WfCSmVBQTWr//t4O5vPbNj5givRLojohJ2KibzxL3nezoJyVTBF/phrJVfX4r1Mv2zXOEkCgYEA45/0x36Q2Mnim7cPSWa58U+6eSVqS3tLYYOOVvDLPhofUrXHGafeom9fd4FrNt6v577Fl5t9geI79Bkb1qcPGSvtnIs/eKNu0NzvwK1Fe9kQKNTWDMKE/3eOTd54O2TZ7FpCQ/tNDngDaOyrqHvYCfnlBO/CIJfwiEO+XRhasW0CgYEAxK8YetulneuGQpcSdalVra9jjdNn4wMFsY0RHg/dem/PHdKSwxmBBxQGpzTL+g2vEU9JoDdHypA1arnHS41Kma9xuFg8Yr2Cb3bhcGTiVEaD02zjlejy6CVty60qLGcWU4g/+oZflhzb6EUQFnl8WJz5C+b4Xa+nzycO7swOCNMCgYEAu9pYowM6+w6x65yKCyOyNQp9dFmCfcTFEzcFE48pzJi2XQYTyIKX5CpR+UhfeSsStQjl/Raf378bh3npVZ8NgNKWCGmK+j62x7xuSO82tt0OzwPHm0Q1irfaQz1ksG+swbhDk+MjVtuIxOD9UdDTHHiVnxtXdJqwMWTnB+F/h4kCgYAJu4xfj/zrCpuTMfyU2/NEa/hmLT7nyd9/QLbHIQvZoizCkgf3JYzv97q4jXFGh2TRW3YOOo4P5QDvrg/BmlVFs5vR/nPGgxAwSdawBB37A55EWRAN/AABItEDEieTGOrO6WAZGosiV30+SiSYqBxSGjpsr1o88JNCOfGQOYK8RwKBgQCzl4oAOeRDTTPgohXAHcwTLE6mhTBPcDbpVdXBa1cjZFE9xnhGB/qI89X04ZhHiEaAlNyyMUOofDoO/l3/aoFLTY9SMDn23tGe5CaODmGfNpugiKTyp4z7ynEqk0zsDJ5WPG3gFR1SmBDNtw0ic+wDtPGEoGFAYyRvsYff0hnN0g=="
	aliPhublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAruInhU3A3tibCtZQLNeVKkTSo66w8vazPS/PySrNVFNv26phNuli9IIYx06/EW6YvXYh/1VNHo9Cv1jA/7jKmGIY6gM8klDTgU3241kpfGy8sAwAYCs/cdd9laF5ofpw0Q9dAiAFNAEb9c+EGx60VyXQg8D6ClPtg7rlginXLGuUsfH3DamF8XI2BqQFS2WLSrLYMkVmCT5t1AuL1yf/uGdTaYdGi+q4GHq9orIOw/ShfynHfNfILmbJ0AIr/jKpKAvJwmRssVV9Yfxp54EgQAD8pk+RAIW43VlhV8xw45d2THP7VwtlaaaoKjfaxZ1mqwW2/pqqjS3K6FMseEGk1wIDAQAB"
)

func Pay() {
	var client, err = alipay.New(appID, privateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(aliPhublicKey)
	if err != nil {
		panic(err)
	}
	//var p = alipay.TradeWapPay{}
	var p = alipay.TradePagePay{}
	p.NotifyURL = "https://www.baidu.com" //支付宝回调
	p.ReturnURL = "https://www.baidu.com" //支付后调转页面
	p.Subject = "eugene-订单支付"             //标题
	p.OutTradeNo = "eugeng"               //传递一个唯一单号
	p.TotalAmount = "10.00"               //金额
	//p.ProductCode = "QUICK_WAP_WAY"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY" //网页支付
	var url2 *url.URL
	url2, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url2.String()
	fmt.Println(payURL)
	// 这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
}
