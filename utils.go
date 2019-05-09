/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush Limit plugin]
 */

package limit

// Some get or a default value
func Some(t interface{}, i interface{}) interface{} {
	if t != nil && t != "" && t != 0 {
		return t
	}
	return i
}
