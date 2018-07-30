package main

import (
	"flag"
	"fmt"
	"github.com/guonaihong/go-soxr"
	"os"
)

func main() {
	in := flag.String("in", "", "(must)Input pcm file")
	out := flag.String("out", "", "(must)Output pcm file")
	inSample := flag.Int("in_sample", 8000, "(must)Input pcm sample rate")
	outSample := flag.Int("out_sample", 16000, "(must)Output pcm sample rate")
	flag.Parse()

	if *in == "" || *out == "" {
		flag.Usage()
		return
	}

	inFd, err := os.Open(*in)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	defer inFd.Close()

	outFd, err := os.Create(*out)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer outFd.Close()

	var spec soxr.IoSpec
	spec.Itype = soxr.SOXR_INT16_I
	spec.Otype = soxr.SOXR_INT16_I
	spec.Scale = 1

	s, err := soxr.Create(float64(*inSample), float64(*outSample), 1, spec)
	if err != nil {
		fmt.Printf("soxr create fail:%s\n", err)
		return
	}
	defer s.Close()

	inBuf := make([]byte, 8000)
	outBuf := make([]byte, 16000)
	for {
		n, err := inFd.Read(inBuf)
		if err != nil {
			break
		}
		n, err = s.Process(inBuf[:n], outBuf)
		if err != nil {
			fmt.Printf("soxr process fail:%s\n", err)
		}

		fmt.Printf("n = %d\n", n)
		outFd.Write(outBuf[:n])
	}
}
