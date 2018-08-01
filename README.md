# go-soxr
Go bindings for the libsoxr audio  down and up sample  library (go包装的soxr音频升降采样库)

#### Usage:
Create an SoxR instance
```go
  
  var spec soxr.IoSpec
  
  //输入采样点精度
  //Input sampling point width
  spec.Itype = soxr.SOXR_INT16_I
  //输出采样点精度
  //Output sampling point width
  spec.Otype = soxr.SOXR_INT16_I
  spec.Scale = 1

  //Create an SoxR instance
  s, err := soxr.Create(float64(*inSample), float64(*outSample), 1, spec)
```
For processing functions, n is the number of bytes exported, and err is an error.
```go
  n, err = s.Process(inBuf[:n], outBuf)
  if err != nil {
    fmt.Printf("soxr process fail:%s\n", err)
  }
```
Destroy the SoxR instance
```go
defer s.Close()
```

#### examples:
```
env GOPATH=`pwd` go build github.com/guonaihong/go-soxr/examples/resample
```

Audio downsampling
```
./resample -in 21.16k.pcm -in_sample 16000 -out out_pcm_file_8k -out_sample 8000
```

Audio upsampling
```
./resample -in 21.8k.pcm -in_sample 8000 -out out_pcm_file_16k -out_sample 16000
```

This raw audio file may be played using ffplay:
```bash
ffplay -f s16le -ar 16000 ./out_pcm_file_16k
ffplay -f s16le -ar 8000 ./out_pcm_file_8k
```
