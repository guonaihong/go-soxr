package soxr

/*
#include <soxr.h>
#cgo LDFLAGS: -l soxr -lm -lgomp
*/
import "C"

import (
	"errors"
	"unsafe"
)

type DataType int

const (
	SOXR_FLOAT32 = DataType(C.SOXR_FLOAT32)
	SOXR_FLOAT64 = DataType(C.SOXR_FLOAT64)
	SOXR_INT32   = DataType(C.SOXR_INT32)
	SOXR_INT16   = DataType(C.SOXR_INT16)
	SOXR_SPLIT   = DataType(C.SOXR_SPLIT)

	SOXR_FLOAT32_I = DataType(C.SOXR_FLOAT32_I)
	SOXR_FLOAT64_I = DataType(C.SOXR_FLOAT64_I)
	SOXR_INT32_I   = DataType(C.SOXR_INT32_I)
	SOXR_INT16_I   = DataType(C.SOXR_INT16_I)

	SOXR_FLOAT32_S = DataType(C.SOXR_FLOAT32_S)
	SOXR_FLOAT64_S = DataType(C.SOXR_FLOAT64_S)
	SOXR_INT32_S   = DataType(C.SOXR_INT32_S)
	SOXR_INT16_S   = DataType(C.SOXR_INT16_S)
)

type IoSpec struct {
	Itype DataType
	Otype DataType
	Scale float64
	Flags uint
}

type Soxr struct {
	soxr  C.soxr_t
	close bool
	in    uint8
	out   uint8
}

func getuint(d DataType) uint8 {
	switch {
	case d == SOXR_INT16, d == SOXR_INT16_I, d == SOXR_INT16_S:
		return 2
	case d == SOXR_FLOAT32, d == SOXR_FLOAT32_I, d == SOXR_FLOAT32_S, d == SOXR_INT32, d == SOXR_INT32_I, d == SOXR_INT32_S:
		return 4
	case d == SOXR_FLOAT64, d == SOXR_FLOAT64_I, d == SOXR_FLOAT64_S:
		return 8
	}

	return 1
}

func Create(inputRate, outputRate float64, numChannels uint32, spec IoSpec) (*Soxr, error) {

	var cspec C.soxr_io_spec_t

	var s Soxr
	cspec.itype = C.soxr_datatype_t(spec.Itype)
	cspec.otype = C.soxr_datatype_t(spec.Otype)
	cspec.scale = C.double(spec.Scale)
	cspec.flags = C.ulong(spec.Flags)

	var e C.soxr_error_t
	s.soxr = C.soxr_create(C.double(inputRate), C.double(outputRate), C.uint(numChannels), &e, &cspec, nil, nil)
	if s.soxr == nil {
		return nil, errors.New(C.GoString(e))
	}

	s.in = getuint(spec.Itype)
	s.out = getuint(spec.Otype)
	return &s, nil
}

func (s *Soxr) Process(in []byte, out []byte) (int, error) {

	odone := C.size_t(0)
	rv := C.soxr_process(s.soxr, C.soxr_in_t(unsafe.Pointer(&in[0])), C.size_t(len(in))/C.size_t(s.in),
		nil, C.soxr_out_t(unsafe.Pointer(&out[0])), C.size_t(len(out))/C.size_t(s.out), &odone)
	if rv != nil {
		return 0, errors.New(C.GoString((*C.char)(rv)))
	}

	return int(odone) * int(s.out), nil
}

func (s *Soxr) Close() {
	if !s.close {
		s.close = true
		C.soxr_delete(s.soxr)
	}
}
