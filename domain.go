package main

import model "github.com/hjwalt/platform/model"

func Parse(s *model.ProtobufSchema) map[string]*model.ProtobufType {
	typeMap := map[string]*model.ProtobufType{}
	for _, t := range s.GetTypes() {
		switch ti := t.Type.(type) {
		case *model.ProtobufType_Message:
			typeMap[ti.Message.GetName()] = t
		case *model.ProtobufType_Enum:
			typeMap[ti.Enum.GetName()] = t
		}
	}
	return typeMap
}

func Flatten(s *model.ProtobufSchema, typeMap map[string]*model.ProtobufType) {
	s.Types = []*model.ProtobufType{}
	for _, t := range typeMap {
		s.Types = append(s.Types, t)
	}
}
