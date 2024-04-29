package main

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/common/pinyin"
	"regexp"
)

func ValidateFullName(fullName string) any {
	reg := regexp.MustCompile(`^[\p{Han}\w]+([\p{Han}\w\s.-]*)$`)
	println(fullName)
	if !reg.MatchString(fullName) {
		fmt.Println("名字中含有特殊字符")
		return nil
	}
	return nil
}

func main() {
	//ValidateFullName("panjinhao")
	urlToken("潘金豪pan0320爱你")
	// [pan jin hao pan0320 ai ni]
}

func urlToken(name string) {
	rs := []rune(name)
	var ss []string
	i := 0
	isLetterDigit := func(r rune) bool {
		return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
	}
	for i < len(rs) {
		if isLetterDigit(rs[i]) {
			s := ""
			for isLetterDigit(rs[i]) {
				s += string(rs[i])
				i++
			}
			ss = append(ss, s)
		} else {
			ss = append(ss, pinyin.SinglePinyin(rs[i], pinyin.NewArgs())[0])
			i++
		}
	}
	fmt.Println(ss)
}
