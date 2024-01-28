package construct

import "github.com/hjwalt/platform/model"

func ProtobufMessage(
	name string,
	fields []*model.ProtobufMessageField,
) *model.ProtobufType {
	return &model.ProtobufType{
		Type: &model.ProtobufType_Message{
			Message: &model.ProtobufMessage{
				Name:   name,
				Fields: fields,
			},
		},
	}
}

func ProtobufMessageBasicField(
	typeName string,
	name string,
	i int64,
) *model.ProtobufMessageField {
	return &model.ProtobufMessageField{
		Field: &model.ProtobufMessageField_BasicField{
			BasicField: &model.ProtobufMessageBasicField{
				Index:    i,
				Name:     name,
				TypeName: typeName,
			},
		},
	}
}

func ProtobufMessageRepeatedField(
	typeName string,
	name string,
	i int64,
) *model.ProtobufMessageField {
	return &model.ProtobufMessageField{
		Field: &model.ProtobufMessageField_RepeatedField{
			RepeatedField: &model.ProtobufMessageRepeatedField{
				Index:    i,
				Name:     name,
				TypeName: typeName,
			},
		},
	}
}

func ProtobufMessageMapField(
	keyTypeName string,
	valueTypeName string,
	name string,
	i int64,
) *model.ProtobufMessageField {
	return &model.ProtobufMessageField{
		Field: &model.ProtobufMessageField_MapField{
			MapField: &model.ProtobufMessageMapField{
				Index:         i,
				Name:          name,
				KeyTypeName:   keyTypeName,
				ValueTypeName: valueTypeName,
			},
		},
	}
}

func ProtobufMessageOneofField(
	name string,
	fields []*model.ProtobufMessageField,
) *model.ProtobufMessageField {
	return &model.ProtobufMessageField{
		Field: &model.ProtobufMessageField_OneofField{
			OneofField: &model.ProtobufMessageOneofField{
				Name:   name,
				Fields: fields,
			},
		},
	}
}
