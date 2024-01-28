package domain

import "github.com/hjwalt/platform/model"

func InUse(typeName string, s *model.ProtobufSchema) bool {
	for _, t := range s.Types {
		if InUseProtobufType(typeName, t) {
			return true
		}
	}
	return false
}

func InUseProtobufType(typeName string, t *model.ProtobufType) bool {
	switch ti := t.Type.(type) {
	case *model.ProtobufType_Message:
		return InUseProtobufMessage(typeName, ti.Message)
	case *model.ProtobufType_Enum:
		return InUseProtobufEnum(typeName, ti.Enum)
	}
	return false
}

func InUseProtobufMessage(typeName string, s *model.ProtobufMessage) bool {
	if typeName == s.GetName() {
		return false
	}

	for _, f := range s.Fields {
		if InUseProtobufMessageField(typeName, f) {
			return true
		}
	}
	return false
}

func InUseProtobufMessageField(typeName string, s *model.ProtobufMessageField) bool {
	switch cf := s.GetField().(type) {
	case *model.ProtobufMessageField_BasicField:
		return InUseProtobufMessageBasicField(typeName, cf.BasicField)
	case *model.ProtobufMessageField_MapField:
		return InUseProtobufMessageMapField(typeName, cf.MapField)
	case *model.ProtobufMessageField_RepeatedField:
		return InUseProtobufMessageRepeatedField(typeName, cf.RepeatedField)
	case *model.ProtobufMessageField_OneofField:
		return InUseProtobufMessageOneofField(typeName, cf.OneofField)
	}
	return false
}

func InUseProtobufMessageBasicField(typeName string, s *model.ProtobufMessageBasicField) bool {
	return typeName == s.GetTypeName()
}

func InUseProtobufMessageMapField(typeName string, s *model.ProtobufMessageMapField) bool {
	return typeName == s.GetKeyTypeName() || typeName == s.GetValueTypeName()
}

func InUseProtobufMessageRepeatedField(typeName string, s *model.ProtobufMessageRepeatedField) bool {
	return typeName == s.GetTypeName()
}

func InUseProtobufMessageOneofField(typeName string, s *model.ProtobufMessageOneofField) bool {
	for _, f := range s.Fields {
		if InUseProtobufMessageField(typeName, f) {
			return true
		}
	}
	return false
}

func InUseProtobufEnum(typeName string, s *model.ProtobufEnum) bool {
	return typeName == s.GetName()
}
