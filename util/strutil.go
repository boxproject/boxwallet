package util

import (
	"bytes"
	"strings"
)

type StrUtil struct {
}

func (*StrUtil) StrAppend(strs ...string) string {
	buff := &bytes.Buffer{}
	for _, v := range strs {
		buff.WriteString(v)
	}
	return buff.String()
}

//切分数值类型为整数和小数
//strCheck：是否正则验证，正则验证耗费性能，100%确认的数传false
func (*StrUtil) SplitStrToNum(str string, strCheck bool) (intstr, decstr string, err error) {
	return splitStrToNum(str, strCheck)
}

func splitStrToNum(str string, strCheck bool) (intstr, decstr string, err error) {
	if strCheck {
		err = canPraseBigFloat(str)
		if err != nil {
			return
		}
	}
	arr := strings.Split(str, ".")
	count := len(arr)
	intstr = arr[0]
	if count == 2 {
		decstr = arr[1]
	}
	return
}

//将数值类型移位
//strCheck：是否正则验证，正则验证耗费性能，100%确认的数传false
func (*StrUtil) MoveDecimalPosition(str string, gap int, strCheck bool) (newstr string, err error) {
	return moveDecimalPosition(str, gap, strCheck)
}

// gap:  代表数值 num * 1e gap
func moveDecimalPosition(str string, gap int, strCheck bool) (newstr string, err error) {
	if gap == 0 {
		newstr = str
		return
	}
	if strCheck {
		err = canPraseBigFloat(str)
		if err != nil {
			return
		}
	}
	//这个方法会用的比较多，内部变量合并，会稍微快100ns/op
	tmp := struct {
		curIdx, strLen, spliIdx, overflow int
		strNoPoint                        string
		buffer                            *bytes.Buffer
	}{
		strings.Index(str, "."),
		len(str),
		0,
		0,
		strings.Replace(str, ".", "", 1),
		new(bytes.Buffer),
	}
	//左移右移可能超过len上限
	if gap > 0 {
		//<-
		if gap >= tmp.curIdx {
			goto HEADADD0
		}
		tmp.spliIdx = tmp.curIdx - gap
		goto INSERTPOINT
	} else {
		//right
		tmp.overflow = tmp.curIdx - gap + 1 - tmp.strLen
		if tmp.overflow > 0 {
			goto TAILADD0
		} else {
			tmp.spliIdx = tmp.curIdx - gap
			goto INSERTPOINT
		}
		return
	}
HEADADD0: //头补0
	tmp.buffer.WriteString("0")
	tmp.buffer.WriteString(".")
	for i := 0; i < gap-tmp.curIdx; i++ {
		tmp.buffer.WriteString("0")
	}
	tmp.buffer.WriteString(tmp.strNoPoint)
	newstr = tmp.buffer.String()
	return
TAILADD0: //尾补0
	tmp.buffer.WriteString(tmp.strNoPoint)
	for i := 0; i < tmp.overflow; i++ {
		tmp.buffer.WriteString("0")
	}
	newstr = tmp.buffer.String()
	return
INSERTPOINT: //Insert string
	tmp.buffer.WriteString(tmp.strNoPoint[:tmp.spliIdx]) //String on the left
	if tmp.curIdx-gap != tmp.strLen-1 {
		tmp.buffer.WriteString(".")
		tmp.buffer.WriteString(tmp.strNoPoint[tmp.spliIdx:]) //String on the right
	}
	newstr = tmp.buffer.String()
	return
}
