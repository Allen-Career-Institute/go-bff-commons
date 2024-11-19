package pageds

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (*Helper) ConvertStructToRawMessage(s *structpb.Struct) (json.RawMessage, error) {
	if s == nil {
		return nil, nil
	}
	jsonBytes, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (*Helper) InterfaceToStructPb(input interface{}) (*structpb.Struct, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var dataStruct structpb.Struct
	if err := protojson.Unmarshal(jsonData, &dataStruct); err != nil {
		return nil, err
	}

	return &dataStruct, nil
}
