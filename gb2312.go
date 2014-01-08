package gogb2312

// just for convert gb2312 to utf-8

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

const (
	UTF8_B1 = 0x80
	UTF8_B2 = 0x80
	UTF8_B3 = 0xe0
)

// param: input: input bytes array
// return: output: output bytes array
//         err: error if there are errors when convert
//         ic: input has been converted
//         oc: output has been converted
func ConvertGB2312(input []byte) (output []byte, err error, ic int, oc int) {
	ilen := len(input)
	output = make([]byte, (ilen/2)*3+3)
	olen := 0
	for i := 0; i < ilen-1; {
		if input[i] <= 0x7f {
			output[olen] = input[i]
			olen++
			i++
		} else {
			gb := int(input[i])<<8 | int(input[i+1])
			u8, ok := gb2312toutf8[gb]
			if !ok {
				err = fmt.Errorf("gb2312 has no character %x, at %d\n", gb, ilen)
				ic = i
				oc = olen
				return
			}
			output[olen] = byte(u8 >> 16)
			olen++
			output[olen] = byte((u8 >> 8) & 0xff)
			olen++
			output[olen] = byte(u8 & 0xff)
			olen++
			i = i + 2
		}
	}
	output = output[0:olen]
	return output, nil, ilen, olen
}

func unicode2utf8(u int) int {
	if u > 0x10000 {
		panic(fmt.Sprintf("unicode for gb2312 is invalid: 0x%x\n", u))
	}
	if u < 0x7f {
		return u
	}

	b1 := (u & 0x3f) | UTF8_B1
	b2 := ((u >> 6) & 0x3f) | UTF8_B2
	b3 := (u >> 12) | UTF8_B3

	u8 := b1 | (b2 << 8) | (b3 << 16)

	return u8
}

func readcp936(fp string) {
	f, err := os.Open(fp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sf, err := os.Create("./gbk2utf8.tmp")
	if err != nil {
		panic(err.Error())
	}
	sf2, err2 := os.Create("./gbk2unicode.tmp")
	if err2 != nil {
		panic(err2.Error())
	}

	sf.WriteString("package gogb2312\n\nvar gb2312toutf8=map[int]int{")
	sf2.WriteString("package gogb2312\n\nvar gb2312tounicode=map[int]int{")

	rd := bufio.NewReader(f)
	i := 0
	line, _, err := rd.ReadLine()
	for ; err == nil; line, _, err = rd.ReadLine() {
		gb, unicode, u8, e := parseline(line, i)
		if e != nil {
			fmt.Print(e)
			i++
			continue
		}
		i++
		serr := savecode(sf, gb, u8)
		if serr != nil {
			fmt.Println("save code failed", serr)
		}
		serr = savecode(sf2, gb, unicode)
		if serr != nil {
			fmt.Println("save code failed", serr)
		}
	}
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
	fmt.Printf("Lines of %s: %d\n", fp, i)
	sf.WriteString("}\n")
	sf.Close()
	sf2.WriteString("}\n")
	sf2.Close()
}

func savecode(rd *os.File, gb, u8 int) error {
	s := fmt.Sprintf("0x%x:0x%x,\n", gb, u8)
	_, err := rd.WriteString(s)
	return err
}

var re_space = regexp.MustCompile(`\s+`)

func parseline(line []byte, i int) (int, int, int, error) {
	comma := bytes.Index(line, []byte{'#'})
	if comma >= 0 {
		line = line[0:comma]
	}
	if l := len(line); l < 4 {
		return -1, -1, -1, fmt.Errorf("line %d length invalid: %d %q\n", i, l, line)
	}
	rl := bytes.TrimSpace(line)
	rl = re_space.ReplaceAll(rl, []byte{' '})
	chs := bytes.Split(rl, []byte{' '})
	if len(chs) != 2 {
		return -1, -1, -1,
			fmt.Errorf("line %d has %d numbers: %s\n", i, len(chs), rl)
	}
	ret, err := strconv.ParseInt(string(bytes.ToLower(chs[0])), 0, 32)
	if err != nil {
		return -1, -1, -1,
			fmt.Errorf("convert %q to int failed at line %d: %s\n", chs[0], i, err)
	}
	gb2312 := int(ret)
	if gb2312 <= 0x7f {
		return gb2312, gb2312, gb2312, fmt.Errorf("No need convert for ascii 0x%x\n", gb2312)
	}
	ret, err = strconv.ParseInt(string(bytes.ToLower(chs[1])), 0, 32)
	if err != nil {
		return -1, -1, -1,
			fmt.Errorf("convert %q to int failed at line %d: %s\n", chs[1], i, err)
	}
	unicode := int(ret)
	utf8 := unicode2utf8(unicode)

	return gb2312, unicode, utf8, nil
}
