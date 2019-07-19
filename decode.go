package gotag

import (
    "fmt"
    "regexp"
    "errors"
    "strings"

    "github.com/golang/protobuf/proto"
    "github.com/CPtung/gotag/protobuf"
    logger "github.com/sirupsen/logrus"
)

func Split(str string, symbol string) (string, string, error) {
    s := strings.Split(str, symbol)
    if len(s) < 2 {
        return str, "", errors.New("no matched")
    }
    return s[0], s[1], nil
}

func DecodeTopic(topic string) (string, string) {
    re, _ := regexp.Compile("/equs/(.*)/tags(.*)")
    str := re.ReplaceAllString(topic, "$1$2")
    source, tag, err := Split(str, "/")
    if err != nil {
        return "", ""
    }
    return source, tag
}

func getDecodeValue(v *mxtag_pb.Value, t int32) *Value{
    value := &Value{}
    switch (t) {
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
    tag.sourceName = data.GetEquipment()
    tag.tagName	= data.GetTag()
    tag.val = getDecodeValue(data.GetValue(), data.GetValueType())
    tag.valType = data.GetValueType()
    tag.ts = data.GetAtMs()
    tag.unit = data.GetUnit()
    return nil
}
