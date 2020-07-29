# PSlog

A golang log wrapper with prefix.

## Features

1. Now is based on `logrus`.
2. Log content with prefix string.
3. Create sub logger with grouped prefix
4. Error with prefix string.

## Usage

```go

package main

import (
	"github.com/menduo/pslog"
	"github.com/sirupsen/logrus"
	"strings"
)

func main() {

	l1 := pslog.NewPsLogger("l1")
	l1sub1 := l1.Sub("sub1", pslog.WithErrorFormat("我的天哪出错了: %s"))

	l2 := pslog.NewPsLogger("l2")

	l2sub1 := l2.Sub("sub1")

	l3 := pslog.NewPsLogger("l3", pslog.WithLogger(logrus.New()))

	l1.Warn("hello")
	l2.Warn("hello")
	l3.Warn("hello")
	l1sub1.Warn("hello")
	l2sub1.Warn("hello")

	// output:
	/*
		WARN[0000] [l1]hello
		WARN[0000] [l2]hello
		WARN[0000] [l3]hello
		WARN[0000] [l1.sub1]hello
		WARN[0000] [l2.sub2]hello
	*/
	err1 := l1.NewErrWithMsgs("this is an error")
	err1sub1 := l1sub1.NewErrWithFormat("i am what i am %s", "what i am...")

	l2.Infoln("err1==nil", err1 == nil)         // false
	l2.Infoln("err1.Error():   ", err1.Error()) // -> [l1]:  this is an error

	l1.Infoln("err1sub1==nil", err1sub1 == nil)  // false
	l3.Infoln("err1sub1 ===:", err1sub1.Error()) // 	l2.Infoln("err1.Error():   ", err1.Error()) // -> [l1]:  this is an error

	l4 := pslog.NewPsLogger("l4", pslog.WithPSGenerOption(func(pslist []string) string {
		return "这是什么前缀啊...." + strings.Join(pslist, "---")
	}))

	l4.Infoln("hello") // 这是什么前缀啊....l4 hello
}


```

## License

MIT License

## Contact

shimenduo@gmail.com

