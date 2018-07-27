package qrcode

import (
	"github.com/qpliu/qrencode-go/qrencode"
	"fmt"
	"github.com/mattn/go-colorable"
	"os"
	"github.com/fumiyas/qrc/lib"
	"github.com/fumiyas/go-tty"
)


type cmdOptions struct {
	Help    bool `short:"h" long:"help" description:"show this help message"`
	Inverse bool `short:"i" long:"invert" description:"invert color"`
}


func pErr(format string, a ...interface{}) {
	fmt.Fprint(os.Stdout, os.Args[0], ": ")
	fmt.Fprintf(os.Stdout, format, a...)
}

func Makeqr(text string) {
	ret := 0
	defer func() { os.Exit(ret) }()

	opts := &cmdOptions{}

	grid, err := qrencode.Encode(text, qrencode.ECLevelL)
	if err != nil {
		pErr("encode failed: %v\n", err)
		ret = 1
		return
	}

	da1, err := tty.GetDeviceAttributes1(os.Stdout)
	if err == nil && da1[tty.DA1_SIXEL] {
		qrc.PrintSixel(os.Stdout, grid, opts.Inverse)
	} else {
		stdout := colorable.NewColorableStdout()
		qrc.PrintAA(stdout, grid, opts.Inverse)
	}
}
