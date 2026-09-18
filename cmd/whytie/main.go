package main

import (
	"fmt"
	"os"

	"github.com/JumoKookbob/whytie/internal/formatter"
	"github.com/JumoKookbob/whytie/internal/reasoning"
	"github.com/JumoKookbob/whytie/internal/scanner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("WhyTie")
		return
	}

	switch os.Args[1] {
	case "scan":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: whytie scan <path>")
			os.Exit(1)
		}

		//? scan 결과를 바로 저장할 것인가?
		//+ 먼저 화면에 출력 가능한 구조를 완성한다
		//< persistence 전에 scanner와 reasoning pipeline을 검증하기 위해
		comments, err := scanner.ScanDir(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			os.Exit(1)
		}

		blocks := reasoning.Group(comments)

		for i, block := range blocks {
			if i > 0 {
				fmt.Println()
			}

			formatter.WriteBlock(os.Stdout, block)
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
