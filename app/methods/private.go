package methods

import (
	"fmt"
)

// GetPrivateRelativeKey 获取两个用户私聊之间的关系字符串,用于数据表等其他关联
func GetPrivateRelativeKey(uid1, uid2 uint64, delimiters ...string) string {
	delimiter := "_"
	if len(delimiters) > 0 {
		delimiter = delimiters[0]
	}
	if uid1 > uid2 {
		return fmt.Sprintf("%d%s%d", uid2, delimiter, uid1)
	}
	return fmt.Sprintf("%d%s%d", uid1, delimiter, uid2)
}
