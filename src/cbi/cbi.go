package cbi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type Pack struct{}

const (
	packName    = "Pack"
	packDefault = 8
)

type cbiError error

// Marshal returns the byte stream cbi encoded of v
func Marshal(v any) ([]byte, error) {
	c := codecState{bytes.NewBuffer(nil)}
	err := c.encode(v)
	if err != nil {
		return nil, err
	}

	buf := append([]byte(nil), c.Bytes()...)
	return buf, nil
}

func Unmarshal(data []byte, v any) error {
	c := codecState{bytes.NewBuffer(data)}
	return c.decode(v)
}

type codecState struct {
	*bytes.Buffer
}

func (c *codecState) encode(v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if ce, ok := r.(cbiError); ok {
				err = ce
			} else {
				panic(r)
			}
		}
	}()

	rv := reflect.ValueOf(v)
	newTypeCodec(rv.Type()).encode(c, rv, nil)
	return nil
}

func (c *codecState) decode(v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if ce, ok := r.(cbiError); ok {
				err = ce
			} else {
				panic(r)
			}
		}
	}()

	rv := reflect.ValueOf(v)
	newTypeCodec(rv.Type()).decode(c, rv, nil)
	if c.Len() != 0 {
		return errors.New("cbi: data overflow")
	}
	return nil
}

type codec interface {
	encode(*codecState, reflect.Value, *codecOpts)
	decode(*codecState, reflect.Value, *codecOpts)
}

type codecOpts struct {
	size      int
	bigEndian bool
}

func newTypeCodec(t reflect.Type) codec {
	switch t.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return newBasicCodec()
	case reflect.String:
		return newStringCodec()
	case reflect.Struct:
		return newStructCodec(t)
	case reflect.Slice:
		return newSliceCodec(t)
	case reflect.Array:
		return newArrayCodec(t)
	case reflect.Pointer:
		return newPtrCodec(t)
	}
	err := errors.New("cbi: unsupported type: " + t.Name())
	panic(cbiError(err))
}

type basicCodec struct{}

func (basicCodec) encode(c *codecState, v reflect.Value, opts *codecOpts) {
	t := v.Type()
	size := int(t.Size())
	bigEndian := false
	if opts != nil {
		size = opts.size
		bigEndian = opts.bigEndian
	}

	// get value
	var val uint64
	switch t.Kind() {
	case reflect.Bool:
		if v.Bool() {
			val = 1
		} else {
			val = 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = uint64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val = v.Uint()
	}

	// encode value to byte array
	addr := unsafe.Pointer(&val)
	buf := make([]byte, 8)
	switch size {
	case 1:
		buf = buf[:1]
		buf[0] = *(*uint8)(addr)
	case 2:
		buf = buf[:2]
		value := *(*uint16)(addr)
		if bigEndian {
			binary.BigEndian.PutUint16(buf, value)
		} else {
			binary.LittleEndian.PutUint16(buf, value)
		}
	case 4:
		buf = buf[:4]
		value := *(*uint32)(addr)
		if bigEndian {
			binary.BigEndian.PutUint32(buf, value)
		} else {
			binary.LittleEndian.PutUint32(buf, value)
		}
	case 8:
		buf = buf[:8]
		value := *(*uint64)(addr)
		if bigEndian {
			binary.BigEndian.PutUint64(buf, value)
		} else {
			binary.LittleEndian.PutUint64(buf, value)
		}
	}
	c.Write(buf)
}

func (basicCodec) decode(c *codecState, v reflect.Value, opts *codecOpts) {
	t := v.Type()
	size := int(t.Size())
	bigEndian := false
	if opts != nil {
		size = opts.size
		bigEndian = opts.bigEndian
	}

	// decode byte array to value
	var val uint64
	buf := make([]byte, 8)
	switch size {
	case 1:
		c.Read(buf[:1])
		val = uint64(buf[0])
	case 2:
		c.Read(buf[:2])
		var v uint16
		if bigEndian {
			v = binary.BigEndian.Uint16(buf)
		} else {
			v = binary.LittleEndian.Uint16(buf)
		}
		val = uint64(v)
	case 4:
		c.Read(buf[:4])
		var value uint32
		if bigEndian {
			value = binary.BigEndian.Uint32(buf)
		} else {
			value = binary.LittleEndian.Uint32(buf)
		}
		val = uint64(value)
	case 8:
		c.Read(buf[:8])
		var value uint64
		if bigEndian {
			value = binary.BigEndian.Uint64(buf)
		} else {
			value = binary.LittleEndian.Uint64(buf)
		}
		val = value
	}

	// set value
	switch t.Kind() {
	case reflect.Bool:
		if val == 0 {
			v.SetBool(false)
		} else {
			v.SetBool(true)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(val))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v.SetUint(val)
	}
}

func newBasicCodec() codec {
	return basicCodec{}
}

type stringCodec struct{}

func (stringCodec) encode(c *codecState, v reflect.Value, opts *codecOpts) {
	s := v.String()
	c.WriteString(strconv.Itoa(len(s)))
	c.WriteString(s)
}

func (stringCodec) decode(c *codecState, v reflect.Value, opts *codecOpts) {
}

func newStringCodec() codec {
	return stringCodec{}
}

type structCodec struct {
	si structInfo
}

func (sc structCodec) encode(c *codecState, v reflect.Value, opts *codecOpts) {
	for i := range sc.si.fields {
		f := &sc.si.fields[i]
		fv := v.Field(f.index)
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				panic(errors.New("cbi: unsupported value: " + v.Type().Name()))
			}
			fv = fv.Elem()
		}
		// padding before encoding filed value
		sc.padding(c, f.align)
		f.codec.encode(c, fv, &codecOpts{f.size, f.bigEndian})
	}
	// padding struct tailer
	sc.padding(c, sc.si.align)
}

func (sc structCodec) decode(c *codecState, v reflect.Value, opts *codecOpts) {
	for i := range sc.si.fields {
		f := &sc.si.fields[i]
		fv := v.Field(f.index)
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				panic(errors.New("cbi: unsupported value: " + v.Type().Name()))
			}
			fv = fv.Elem()
		}
		// unpadding before decoding filed value
		sc.unpadding(c, f.align)
		f.codec.decode(c, fv, &codecOpts{f.size, f.bigEndian})
	}
	// unpadding struct tailer
	sc.unpadding(c, sc.si.align)
}

func (sc structCodec) padding(c *codecState, align int) {
	offset := c.Len()
	size := 0
	for offset%align != 0 {
		offset++
		size++
	}
	if size == 0 {
		return
	}
	zeros := make([]byte, size)
	c.Write(zeros)
}

func (sc structCodec) unpadding(c *codecState, align int) {
	offset := c.Len()
	size := 0
	for offset%align != 0 {
		offset--
		size++
	}
	if size == 0 {
		return
	}
	zeros := make([]byte, size)
	c.Read(zeros[:])
}

func newStructCodec(t reflect.Type) codec {
	se := structCodec{si: typeStruct(t)}
	return se
}

type sliceCodec struct {
	arrayEnc codec
}

func (sc sliceCodec) encode(e *codecState, v reflect.Value, opts *codecOpts) {
	if v.IsNil() {
		return
	}
	sc.arrayEnc.encode(e, v, opts)
}

func (sc sliceCodec) decode(e *codecState, v reflect.Value, opts *codecOpts) {
}

func newSliceCodec(t reflect.Type) codec {
	return sliceCodec{newArrayCodec(t)}
}

type arrayCodec struct {
	elemCodec codec
}

func (ac arrayCodec) encode(e *codecState, v reflect.Value, opts *codecOpts) {
	n := v.Len()
	e.WriteString(strconv.Itoa(n))
	e.WriteByte('[')
	for i := 0; i < n; i++ {
		if i > 0 {
			e.WriteByte(',')
		}
		ac.elemCodec.encode(e, v.Index(i), opts)
	}
	e.WriteByte(']')
}

func (ac arrayCodec) decode(e *codecState, v reflect.Value, opts *codecOpts) {
}

func newArrayCodec(t reflect.Type) codec {
	return arrayCodec{newTypeCodec(t.Elem())}
}

type ptrCodec struct {
	elemCodec codec
}

func (pc ptrCodec) encode(e *codecState, v reflect.Value, opts *codecOpts) {
	if v.IsNil() {
		return
	}
	pc.elemCodec.encode(e, v.Elem(), opts)
}

func (pc ptrCodec) decode(e *codecState, v reflect.Value, opts *codecOpts) {
	pc.elemCodec.decode(e, v.Elem(), opts)
}

func newPtrCodec(t reflect.Type) codec {
	return ptrCodec{newTypeCodec(t.Elem())}
}

type structInfo struct {
	align  int
	fields []field
}

type field struct {
	index     int
	align     int
	size      int
	bigEndian bool
	codec     codec
}

func typeStruct(t reflect.Type) structInfo {
	pack := packDefault
	var fields []field
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Anonymous {

		} else if !sf.IsExported() {
			// Private field
			continue
		}
		tag := sf.Tag.Get("cbi")
		if tag == "-" {
			continue
		}
		size, opts := parseTag(tag)
		if sf.Name == packName {
			pack = size
			continue
		}
		ft := sf.Type
		if ft.Name() == "" && ft.Kind() == reflect.Pointer {
			// Follow pointer.
			ft = ft.Elem()
		}
		if size == 0 && ft.Kind() == reflect.Struct {
			si := typeStruct(ft)
			size = si.align
		}
		field := field{
			index:     i,
			size:      size,
			bigEndian: opts.contains("bigEndian"),
			codec:     newTypeCodec(ft),
		}
		fields = append(fields, field)
	}

	// set field align = min(pack, field.size)
	maxFieldSize := 1
	for i := range fields {
		f := &fields[i]
		f.align = pack
		if f.size < pack {
			f.align = f.size
		}
		if f.size > maxFieldSize {
			maxFieldSize = f.size
		}
	}

	// set struct align = min(pack, maxFieldSize)
	align := pack
	if maxFieldSize < pack {
		align = maxFieldSize
	}

	return structInfo{align: align, fields: fields}
}

var table = map[string]int{
	"1":         1,
	"2":         2,
	"4":         4,
	"8":         8,
	"bool":      1,
	"char":      1,
	"uchar":     1,
	"short":     2,
	"ushort":    2,
	"int":       4,
	"uint":      4,
	"float":     4,
	"double":    8,
	"longlong":  8,
	"ulonglong": 8,
}

func parseTag(tag string) (int, tagOptions) {
	c, opt, _ := strings.Cut(tag, ",")
	size, ok := table[c]
	if !ok {
		size = table["int"]
	}
	return size, tagOptions(opt)
}

type tagOptions string

func (o tagOptions) contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}
