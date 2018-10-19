package cxtonsms

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lvzhihao/goutils"
)

const (
	CTXON_SEED_TIME = "20060102150405" // YYYYMMDDHHMISS
)

var (
	TimeZone                           string         = "Asia/Shanghai"  // 时区设置
	TimeLocation                       *time.Location                    // 当前时区
	DefaultTimeout                     time.Duration  = 60 * time.Second // 默认超时
	DefaultTransportInsecureSkipVerify bool           = true             // 是否跳过ssl证书验证
	DefaultTransportDisableCompression bool           = true             // 不使用压缩
)

func init() {
	SetTimeZone(TimeZone)
}

// 设置时区
func SetTimeZone(zone string) error {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return err
	} else {
		TimeZone = zone
		TimeLocation = loc
		return nil
	}
}

type Client struct {
	ApiPrefix  string       // Api 域名
	Header     Header       // 头信息
	httpClient *http.Client // http client
}

type Api interface {
	Map() (map[string]string, error)
}

func NewClient(apiPrefix, name, password string) *Client {
	c := &Client{
		ApiPrefix: apiPrefix,
	}
	c.Init()
	c.Header.Set(name, password)
	return c
}

// 初始化client
func (c *Client) Init() {
	c.httpClient = &http.Client{
		Timeout: DefaultTimeout,
	}
	URL, err := url.Parse(c.ApiPrefix)
	if err != nil && strings.ToLower(URL.Scheme) == "https" {
		c.httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: DefaultTransportInsecureSkipVerify,
			},
			DisableCompression: DefaultTransportDisableCompression,
		}
	}
}

// 请求
func (c *Client) Action(action string, ctx interface{}) ([]byte, error) {
	req, err := c.request(action, ctx)
	if err != nil {
		return nil, err
	}
	rsp, err := c.Do(req)
	return c.Scan(rsp, err)
}

func (c *Client) request(action string, ctx interface{}) (*http.Request, error) {
	var err error
	p := url.Values{}
	switch ctx.(type) {
	case nil:
		// nil
	case map[string]interface{}:
		for k, v := range ctx.(map[string]interface{}) {
			p.Set(k, goutils.ToString(v))
		}
	case map[string]string:
		for k, v := range ctx.(map[string]string) {
			p.Set(k, v)
		}
	default:
		if _, ok := ctx.(Api); ok {
			var maps map[string]string
			maps, err = ctx.(Api).Map()
			if err == nil {
				for k, v := range maps {
					p.Set(k, v)
				}
			}
		} else {
			err = errors.New("params error")
		}
	}
	if err != nil {
		return nil, err
	}
	p.Set("name", c.Header.Name)
	seed, key := c.Header.GenerateKey()
	p.Set("seed", seed)
	p.Set("key", key)
	req, err := http.NewRequest(
		"POST",
		strings.TrimRight(c.ApiPrefix, "/")+"/"+strings.TrimLeft(action, "/"),
		bytes.NewBufferString(p.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", goutils.ToString(len(p.Encode())))
	log.Println(req)
	return req, nil
}

// 审查请求结果
func (c *Client) Scan(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(resp)
	log.Printf("%s\n", b)
	if err, ok := Errors[string(b)]; ok {
		return b, errors.New(err)
	} else {
		return b, nil
	}
}

// 发送请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// 发送短信
func (c *Client) SendStrongGBK(strong *Strong) ([]byte, error) {
	return c.Action("eums/send_strong.do", strong)
}

// 发送短信 utf8
func (c *Client) SendStrongUTF8(strong *Strong) ([]byte, error) {
	return c.Action("eums/utf8/send_strong.do", strong)
}

type Header struct {
	Name     string // 帐号，由网关分配
	Password string // 密码
	//Seed     string // 当前时间，格式: YYYYMMDDHHMISS
	//Key      string // md5(md5(password) + seed))
}

func (c *Header) Set(name, password string) *Header {
	c.SetName(name)
	c.SetPassword(password)
	return c
}

// 设置帐号
func (c *Header) SetName(name string) *Header {
	c.Name = name
	return c
}

// 设置密码
func (c *Header) SetPassword(password string) *Header {
	c.Password = password
	return c
}

func (c *Header) GenerateKey() (string, string) {
	seed := time.Now().Format(CTXON_SEED_TIME)
	p1 := fmt.Sprintf("%x", md5.Sum([]byte(c.Password)))
	return seed, fmt.Sprintf("%x", md5.Sum([]byte(p1+seed)))
}
