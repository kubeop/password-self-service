package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"password-self-service/pkg/config"
	"password-self-service/pkg/logging"

	"github.com/go-ldap/ldap"
	"golang.org/x/text/encoding/unicode"
)

type Attributes struct {
	Username       string `json:"username"`       // 用户名
	Nickname       string `json:"nickname"`       // 昵称
	Email          string `json:"email"`          // 邮箱
	Mobile         string `json:"mobile"`         // 手机
	DN             string `json:"dn"`             // 用户DN
	PasswordExpire string `json:"passwordExpire"` // 密码过期时间
}

type Client struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	TLS       bool   `json:"tls"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Base      string `json:"base"`
	Filter    string `json:"filter"`
	TimeLimit int    `json:"timeLimit"`
	SizeLimit int    `json:"sizeLimit"`
	Conn      *ldap.Conn
}

func NewLdapClient() *Client {
	return &Client{
		Host:      config.Setting.Ldap.Host,
		Port:      config.Setting.Ldap.Port,
		TLS:       config.Setting.Ldap.TLS,
		Username:  config.Setting.Ldap.Username,
		Password:  config.Setting.Ldap.Password,
		Base:      config.Setting.Ldap.Base,
		Filter:    config.Setting.Ldap.Filter,
		TimeLimit: config.Setting.Ldap.TimeLimit,
		SizeLimit: config.Setting.Ldap.SizeLimit,
	}
}

func (l *Client) Connect() error {
	var err error
	if l.TLS {
		l.Conn, err = ldap.DialTLS("tcp", fmt.Sprintf("%v:%v", l.Host, l.Port), &tls.Config{
			InsecureSkipVerify: true,
		})

	} else {
		l.Conn, err = ldap.Dial("tcp", fmt.Sprintf("%v:%v", l.Host, l.Port))
	}

	if err != nil {
		return err
	}

	err = l.Conn.Bind(l.Username, l.Password)
	if err != nil {
		return err
	}

	return nil
}

func (l *Client) Search(username string) (*Attributes, error) {
	defer l.Conn.Close()
	searchRequest := ldap.NewSearchRequest(l.Base,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, l.TimeLimit, false,
		fmt.Sprintf(l.Filter, username),
		[]string{"dn", "displayName", "sAMAccountName", "mail", "mobile", "employeeNumber", "title", "department", "manager", "pwdLastSet"},
		nil)

	sr, err := l.Conn.SearchWithPaging(searchRequest, uint32(l.SizeLimit))
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) == 0 {
		return nil, errors.New("user is not found")
	}

	user := &Attributes{
		Username: sr.Entries[0].GetAttributeValue("sAMAccountName"),
		Nickname: sr.Entries[0].GetAttributeValue("displayName"),
		Email:    sr.Entries[0].GetAttributeValue("mail"),
		Mobile:   sr.Entries[0].GetAttributeValue("mobile"),
		DN:       sr.Entries[0].DN,
	}

	return user, err
}

func (l *Client) Login(username, password string) error {
	user, err := l.Search(username)
	if err != nil {
		return err
	}

	err = l.Conn.Bind(user.DN, password)
	if err != nil {
		return err
	}

	return nil
}

func (l *Client) ModifyPassword(userDn, newPassword string) error {
	defer l.Conn.Close()
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// The password needs to be enclosed in quotes
	pwdEncoded, err := utf16.NewEncoder().String(fmt.Sprintf("\"%s\"", newPassword))
	if err != nil {
		return err
	}

	modifyReq := ldap.NewModifyRequest(userDn, nil)
	modifyReq.Replace("unicodePwd", []string{pwdEncoded})

	err = l.Conn.Modify(modifyReq)
	if err != nil {
		logging.Logger().Sugar().Errorf("password change failed: %s", err.Error())
		return err
	}

	return nil
}

func (l *Client) UnlockAccount(userDn string) error {
	defer l.Conn.Close()
	modifyReq := ldap.NewModifyRequest(userDn, nil)
	modifyReq.Replace("lockoutTime", []string{"0"})

	err := l.Conn.Modify(modifyReq)
	if err != nil {
		return err
	}

	return nil
}

func (l *Client) CheckPasswordExpired() ([]*Attributes, error) {
	// 过滤掉设置密码永不过期，用户状态异常的用户
	//filter := "(sAMAccountName=san.zhang)"
	filter := "(&(objectCategory=person)(!(pwdLastSet=0))(!(userAccountControl:1.2.840.113556.1.4.803:=65536))(!(userAccountControl:1.2.840.113556.1.4.803:=2))(!(userAccountControl:1.2.840.113556.1.4.803:=8388608)))"

	searchRequest := ldap.NewSearchRequest(l.Base,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, l.TimeLimit, false,
		filter,
		[]string{"dn", "displayName", "sAMAccountName", "mail", "mobile", "pwdLastSet"},
		nil)

	sr, err := l.Conn.SearchWithPaging(searchRequest, uint32(l.SizeLimit))
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) == 0 {
		return nil, errors.New("user does not exist")
	}

	user := make([]*Attributes, 0)
	for _, u := range sr.Entries {
		user = append(user, &Attributes{
			Username:       u.GetAttributeValue("sAMAccountName"),
			Nickname:       u.GetAttributeValue("displayName"),
			Email:          u.GetAttributeValue("mail"),
			Mobile:         u.GetAttributeValue("mobile"),
			DN:             sr.Entries[0].DN,
			PasswordExpire: u.GetAttributeValue("pwdLastSet"),
		})
	}

	defer l.Conn.Close()
	return user, err
}
