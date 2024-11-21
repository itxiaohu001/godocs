package example

// User 表示系统中的用户实体
// 包含用户的基本信息
type User struct {
	// ID 是用户的唯一标识
	ID int `json:"id"`
	// Name 是用户的显示名称
	Name string `json:"name"` // 用户名称
	// Email 是用户的邮箱地址
	Email   string   `json:"email"` // 电子邮箱
	age     int      // 私有字段：年龄
	Address *Address `json:"address"` // 用户地址信息
}

// Address 表示地址信息
type Address struct {
	Street  string `json:"street"`  // 街道
	City    string `json:"city"`    // 城市
	Country string `json:"country"` // 国家
	ZipCode string `json:"zipCode"` // 邮政编码
	Contact Contact `json:"contact"` // 联系信息
}

// Contact 表示联系人信息
type Contact struct {
	Phone   string `json:"phone"`   // 电话号码
	Mobile  string `json:"mobile"`  // 手机号码
	WeChat  string `json:"wechat"`  // 微信号
	Primary bool   `json:"primary"` // 是否为主要联系方式
}

type privateStruct struct {
	field string
}
