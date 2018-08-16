package main

import (
	"flag"
	"fmt"
	"github.com/guonaihong/go-soxr"
	"os"
)

func bzero(buf []byte) {
	for k, _ := range buf {
		buf[k] = 0
	}
}

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

	fmt.Printf("in sample(%d) out sample(%d)\n", *inSample, *outSample)
	s, err := soxr.Create(float64(*inSample), float64(*outSample), 1, spec)
	if err != nil {
		fmt.Printf("soxr create fail:%s\n", err)
		return
	}
	defer s.Close()

	if *inSample%2 != 0 {
		if (*inSample)--; *inSample < 0 {
			panic("in_sample < 0")
		}
	}

	inBuf := make([]byte, *inSample)
	outBuf := make([]byte, *outSample+8000)
	for {

		n, err := inFd.Read(inBuf)
		if err != nil {
			break
		}

		bzero(outBuf)
		n, err = s.Process(inBuf[:n], outBuf)
		if err != nil {
			fmt.Printf("soxr process fail:%s\n", err)
		}

		fmt.Printf("n = %d\n", n)
		outFd.Write(outBuf[:n])
	}
}
