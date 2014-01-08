package gogb2312

import (
	"fmt"
	"testing"
)

func print_utf8(u8 int) {
	fmt.Printf("%x %x %x\n", u8>>16, (u8>>8)&0xff, u8&0xff)
}

func Test_unicode2utf8(t *testing.T) {
	print_utf8(unicode2utf8(0x90ed))
}

func Test_RE(t *testing.T) {
	var ss = []string{
		"who  \t are u \t\t  ? ",
		"i \t \t am\t\t\t   boy.  ",
		"i \t \t am\t\t\t 90  boy.  ",
	}
	var gb = []byte("\x90\xae")
	var b = byte(98)

	for _, s := range ss {
		ret := re_space.ReplaceAllString(s, " ")
		_ = ret
		//fmt.Println(s, " - ", ret)
	}
	_ = gb
	_ = b
}

func Test_Convert(t *testing.T) {
	readcp936("./CP936.TXT")
}
