package gotag

import (
	"errors"
	"fmt"
	"regexp"

	mxtag_pb "github.com/CPtung/gotag/protobuf"
	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
)

func DecodeTopic(topic string) (string, string, error) {
	re := regexp.MustCompile("/(\\+|[\\w-]+)")
	matches := re.FindAllStringSubmatch(topic, -1)
	if matches != nil && len(matches) > 3 {
		return matches[1][1], matches[3][1], nil
	}
	return "", "", errors.New("invalid topic format")
}

func getDecodeValue(v *mxtag_pb.Value, t int32) *Value {
	value := &Value{}
	switch t {
	case TAG_VALUE_TYPE_BOOLEAN:
		value.i = v.GetIntValue()
		break
	case TAG_VALUE_TYPE_INT, TAG_VALUE_TYPE_INT8,
		TAG_VALUE_TYPE_INT16, TAG_VALUE_TYPE_INT32:
		value.i = v.GetIntValue()
		break
	case TAG_VALUE_TYPE_UINT, TAG_VALUE_TYPE_UINT8,
		TAG_VALUE_TYPE_UINT16, TAG_VALUE_TYPE_UINT32:
		value.u = v.GetUintValue()
		break
	case TAG_VALUE_TYPE_FLOAT:
		value.f = v.GetFloatValue()
		break
	case TAG_VALUE_TYPE_DOUBLE:
		value.d = v.GetDoubleValue()
		break
	case TAG_VALUE_TYPE_STRING:
		value.s = v.GetStrValue()
		break
	case TAG_VALUE_TYPE_BYTEARRAY:
		value.b = v.GetBytesValue()
		break
	default:
		logger.Debugf("decode invalid value type: %v", t)
		break
	}
	return value
}

func DecodePayload(payload []byte, tag *Tag) error {
	data := &mxtag_pb.Tag{}
	if err := proto.Unmarshal(payload, data); err != nil {
		return errors.New(fmt.Sprintf("unmarshaling error: %v", err))
	}
	tag.SourceName = data.GetEquipment()
	tag.TagName = data.GetTag()
	tag.Val = getDecodeValue(data.GetValue(), data.GetValueType())
	tag.ValType = data.GetValueType()
	tag.Ts = data.GetAtMs()
	tag.Unit = data.GetUnit()
	return nil
}
