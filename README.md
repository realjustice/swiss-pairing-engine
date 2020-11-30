# swiss-pairing-engine

[![ren_forbid compliant](https://img.shields.io/badge/swissPairingEngine-realjustice-green.svg)](https://github.com/realjustice/renju_forbid)

瑞士制赛事编排引擎，适用于围棋类赛事的对阵生成

本仓库包含以下内容：

1. 瑞士制赛事编排，生成对阵
2. 对局结果的录入

## 内容列表

- [安装](#安装)
- [快速开始](#快速开始)
  - [编排并生成对阵](#编排并生成对阵)
  - [录入结果](#录入结果)
- [示例](#示例)
- [相关仓库](#相关仓库)
- [维护者](#维护者)
- [如何贡献](#如何贡献)
- [使用许可](#使用许可)

## 安装

本项目使用 [go](https://gomirrors.org/)语言开发。请确保你本地安装了它们。

```sh
$ go get -u github.com/realjustice/swiss-pairing-engine
```

## 快速开始

```sh
$ cd example
```

编排并生成对阵

```sh
$ cat example/pair/main.go
```

```go
package main

import (
   "flag"
   "fmt"
   "github.com/realjustice/swiss-pairing-engine/src/gotha"
   "io"
   "io/ioutil"
   "log"
   "os"
)

var (
   round  = flag.Int("round", 1, "The round number")
   system = flag.String("system", "SWISS", "the pair system")
)

func main() {
   flag.Parse()

   // Step1 init the pair engine
   g := gotha.NewGotha()

   // Step2 import the data resource
   // you can import from the xml file
   filePath := `../demo.xml` //  your file path
   importFromXMLFile(filePath, g)

   // or from bytes
   // importFromBytes(filePath, g)

   // Step3 chose the pair system
   // Currently only supports Swiss-made arrangements
   g.SelectSystem(*system)
   t := g.GetTournament()

   // Step4 choose the players （via keyString）
   t.SetSelectedPlayers([]string{"ARNAUDANCELIN","WURZINGERRALF"})

   // Step5 pair
   t.Pair(*round)
   for _, game := range t.SortGameByTableNumber() {
      fmt.Printf("white : %s  <> black : %s\n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName)
   }

   //  Step6 you will get a io.reader
   rd, err := g.IO.FlushGameToXML(t.SortGameByTableNumber())
   if err != nil {
      log.Fatal(err)
   }

   // overwrite your xml file (Optional operation)
   file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
   if err != nil {
      log.Fatal(err)
   }
   defer func() { _ = file.Close() }()
   _, err = io.Copy(file, rd)
   if err != nil {
      log.Fatal(err)
   }
}

func importFromXMLFile(filePath string, g *gotha.Gotha) {
   if err := g.ImportFromXMLFile(filePath); err != nil {
      log.Fatal(err)
   }
}

func importFromBytes(filePath string, g *gotha.Gotha) {
   tempXML, err := os.Open(filePath)
   if err != nil {
      log.Fatal(err)
   }

   defer func() { _ = tempXML.Close() }()
   fd, err := ioutil.ReadAll(tempXML)
   if err != nil {
      log.Fatal(err)
   }
   if err = g.ImportFromBytes(fd); err != nil {
      log.Fatal(err)
   }
}

```

录入结果

```sh
$ cat example/record_result/main.go
```

```go
package main

import (
	"flag"
	"github.com/realjustice/swiss-pairing-engine/src/gotha"
	"io"
	"log"
	"os"
)

var (
	round       = flag.Int("round", 0, "The round number")
	tableNumber = flag.Int("table", 0, "The table number")
	result      = flag.String("result", "RESULT_WHITEWINS", "The Result")
)

func main() {
	flag.Parse()
	// Step1 init the pair engine
	g := gotha.NewGotha()

	// Step2 import the data resource
	filePath := `../demo.xml` //  your file path
	if err := g.ImportFromXMLFile(filePath); err != nil {
		log.Fatal(err)
	}
	t := g.GetTournament()

	// Step2 set the result via roundNumber and tableNumber
	// 	RESULT_WHITEWINS 白胜
	//	RESULT_BLACKWINS 黑胜
	//	RESULT_EQUAL 和棋
	//	RESULT_BOTHWIN 双胜
	//	RESULT_BOTHLOSE 双负
	t.SetGameResult(*round, *tableNumber, *result)

	//  Step3 you will get a io.reader
	rd, err := g.IO.FlushGameToXML(t.SortGameByTableNumber())
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	_, err = io.Copy(file, rd)
	if err != nil {
		log.Fatal(err)
	}
}
```

## 示例

更多示例，请参考 [example-readmes](example-readmes/)。

## 相关仓库

- [maximum_weight_matching](https://github.com/realjustice/maximum_weight_matching) —最大权重匹配算法。

## 维护者

[@realjustice](https://github.com/realjustice)。

## 如何贡献

[提一个 Issue](https://github.com/RichardLitt/standard-readme/issues/new) 或者提交一个 Pull Request，或者发送邮件至 [z_s_c_p@163.com](z_s_c_p@163.com)

## 使用许可

[MIT](LICENSE) © realjustice