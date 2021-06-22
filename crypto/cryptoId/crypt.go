package cryptoId

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/lights-T/lib-go/util/uid"
	"github.com/modern-go/reflect2"
)

// 结构体上的json标签必须带有cv_id属性才会自动转换id
func MarshalJSON(v interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ext := &IdExtension{}
	json.RegisterExtension(ext)
	return json.Marshal(v)
}

// 结构体上的json标签必须带有cv_id属性才会自动转换id
func UnmarshalJSON(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ext := &IdExtension{}
	json.RegisterExtension(ext)
	return json.Unmarshal(data, v)
}

type idCodec struct {
}

func (codec *idCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	i := *((*int64)(ptr))
	newId := uid.UID(i).Encode()
	stream.WriteInt64(newId)
}

func (codec *idCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if !iter.ReadNil() {
		newId := int64(uid.UID(iter.ReadInt64()).Decode())
		*((*int64)(ptr)) = newId
	}
}
func (codec *idCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int64)(ptr)) == 0
}

type idSliceEncoder struct {
	sliceType *reflect2.UnsafeSliceType
}

func encoderIdOfSlice(typ *reflect2.UnsafeSliceType) jsoniter.ValEncoder {
	return &idSliceEncoder{sliceType: typ}
}

func (encoder *idSliceEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	if encoder.sliceType.UnsafeIsNil(ptr) {
		stream.WriteEmptyArray()
		return
	}
	length := encoder.sliceType.UnsafeLengthOf(ptr)
	if length == 0 {
		stream.WriteEmptyArray()
		return
	}
	stream.WriteArrayStart()

	for i := 0; i < length; i++ {
		n := encoder.sliceType.UnsafeGetIndex(ptr, i)
		newId := uid.UID(*((*int64)(n))).Encode()
		stream.WriteInt64(newId)
		if i != length-1 {
			stream.WriteMore()
		}
	}
	stream.WriteArrayEnd()
	if stream.Error != nil && stream.Error != io.EOF {
		stream.Error = fmt.Errorf("%v: %s", encoder.sliceType, stream.Error.Error())
	}
}

func (encoder *idSliceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.sliceType.UnsafeLengthOf(ptr) == 0
}

func decoderOfSlice(typ reflect2.Type) jsoniter.ValDecoder {
	sliceType := typ.(*reflect2.UnsafeSliceType)
	//decoder := decoderOfType(ctx.append("[sliceElem]"), sliceType.Elem())
	return &idSliceDecoder{sliceType: sliceType}
}

type idSliceDecoder struct {
	sliceType *reflect2.UnsafeSliceType
	//elemDecoder jsoniter.ValDecoder
}

func (decoder *idSliceDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.doDecode(ptr, iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.sliceType, iter.Error.Error())
	}
}

func (decoder *idSliceDecoder) doDecode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	sliceType := decoder.sliceType
	b := iter.SkipAndReturnBytes()

	if b[0] == 'n' {
		sliceType.UnsafeSetNil(ptr)
		return
	}
	if b[0] != '[' {
		iter.ReportError("decode slice", "expect [ or n, but found "+string([]byte{b[0]}))
		return
	}
	empty := true
	val := ""
	length := 0
	for i := 1; i < len(b); i++ {
		if b[i] == ' ' {
			continue
		}
		if empty && b[i] != ' ' && b[i] != ']' {
			empty = false
		}
		if len(val) != 0 && (b[i] == ',' || b[i] == ']') {
			i, _ := strconv.Atoi(val)
			newId := int64(uid.UID(i).Decode())
			idx := length
			length += 1
			sliceType.UnsafeGrow(ptr, length)
			elemPtr := sliceType.UnsafeGetIndex(ptr, idx)
			*((*int64)(elemPtr)) = newId
			val = ""
			continue
		}
		val += string(b[i])
	}

	if empty && b[len(b)-1] == ']' {
		sliceType.UnsafeSet(ptr, sliceType.UnsafeMakeSlice(0, 0))
		return
	}
}

type IdExtension struct {
}

// UpdateStructDescriptor No-op
func (extension *IdExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for i, item := range structDescriptor.Fields {
		tag := item.Field.Tag().Get("json")
		tagArr := strings.Split(tag, ",")
		for _, k := range tagArr {
			if k == "cv_id" {
				t := item.Field.Type().Kind().String()
				switch t {
				case "int64":
					structDescriptor.Fields[i].Encoder = new(idCodec)
					structDescriptor.Fields[i].Decoder = new(idCodec)
				case "slice":
					if item.Field.Type().(*reflect2.UnsafeSliceType).Type1().Elem().Name() == "int64" {
						sliceType := item.Field.Type().(*reflect2.UnsafeSliceType)
						structDescriptor.Fields[i].Encoder = encoderIdOfSlice(sliceType)
						structDescriptor.Fields[i].Decoder = decoderOfSlice(sliceType)
					}
				}
				break
			}

		}
	}
}

// CreateMapKeyDecoder No-op
func (extension *IdExtension) CreateMapKeyDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	return nil
}

// CreateMapKeyEncoder No-op
func (extension *IdExtension) CreateMapKeyEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	return nil
}

// CreateDecoder No-op
func (extension *IdExtension) CreateDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	return nil
}

// CreateEncoder No-op
func (extension *IdExtension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	return nil
}

// DecorateDecoder No-op
func (extension *IdExtension) DecorateDecoder(typ reflect2.Type, decoder jsoniter.ValDecoder) jsoniter.ValDecoder {
	return decoder
}

// DecorateEncoder No-op
func (extension *IdExtension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	return encoder
}
