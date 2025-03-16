package appTypes

import "encoding/json"

type Register int

const (
	Email Register = iota // 邮箱验证码注册
	QQ                    // QQ登录注册
)

// MarshalJSON 实现了 json.Marshaler 接口
func (r Register) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToString())
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口
func (r *Register) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*r = ToRegister(str)
	return nil
}

// // String 枚举值转中文
func (r Register) ToString() string {
	var str string
	switch r {
	case Email:
		str = "邮箱"
	case QQ:
		str = "QQ"
	default:
		str = "未知"
	}
	return str
}

// ToRegister 中文转枚举值
func ToRegister(str string) Register {
	switch str {
	case "邮箱":
		return Email
	case "QQ":
		return QQ
	default:
		return -1
	}
}
