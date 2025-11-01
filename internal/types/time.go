package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	formatTime = "2006-01-02 15:04:05"
)

// GxTime 自定义时间类型，提供给json和sqlx使用
type GxTime struct {
	time.Time
}

// MarshalJSON实现json.Marshaler接口
func (gx GxTime) MarshalJSON() ([]byte, error) {
	if gx.IsZero() { // 零值返回空，而不是"0001-01-01T00:00:00Z"
		return []byte(`""`), nil
	}

	// 格式化时间
	str := fmt.Sprintf("\"%s\"", gx.Format(formatTime))
	return []byte(str), nil
}

// UnmarshalJSON实现json.Unmarshaler接口
func (gx *GxTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == `""` {
		gx.Time = time.Time{}
		return nil
	}
	parsed, err := time.ParseInLocation(`"`+formatTime+`"`, str, time.Local)
	if err != nil {
		return err
	}
	gx.Time = parsed
	return nil
}

//======================================

// sqlx的底层是基于database/sql,

// Value写入数据库时调用，实现driver.Valuer接口
func (gx GxTime) Value() (driver.Value, error) {
	if gx.IsZero() {
		return nil, nil
	}
	return gx.Time, nil
}

// Scan从数据库读出时调用，实现sql.Scanner接口
func (gx *GxTime) Scan(value interface{}) error {
	if value == nil {
		*gx = GxTime{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		gx.Time = v
	case []byte:
		parsed, err := time.ParseInLocation(formatTime, string(v), time.Local)
		if err != nil {
			return err
		}
		gx.Time = parsed
	case string:
		parsed, err := time.ParseInLocation(formatTime, v, time.Local)
		if err != nil {
			return err
		}
		gx.Time = parsed
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	return nil
}
