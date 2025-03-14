/**
 * Created by GoLand.
 * User: buzzlight.frank@qq.com
 * Date: 2025/3/14
 * Time: 12:10
 */
package utils

import "time"

func GoMerge(arr1 []interface{}, arr2 []interface{}) []interface{} {
	for _, val := range arr2 {
		arr1 = append(arr1, val)
	}
	return arr1
}

func GoRepeat(str string, num int) string {
	var i int
	newStr := ""
	if num != 0 {
		for i = 0; i < num; i++ {
			newStr += str
		}
	}
	return newStr
}

func round(a int, b int) int {
	rem := a % b
	dis := a / b
	if rem > 0 {
		return dis + 1
	} else {
		return dis
	}
}

func MDate(times time.Time) string {
	return times.Format("2006-01-02 15:04:05")
}
