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
