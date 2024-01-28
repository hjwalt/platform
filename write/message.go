package write

import (
	"bufio"
	"fmt"

	"github.com/hjwalt/platform/model"
)

func WriteProtobufMessage(w *bufio.Writer, f *model.ProtobufMessage) {
	w.WriteString(fmt.Sprintf("message %s { \n", f.GetName()))
	for _, fi := range f.GetFields() {
		WriteProtobufMessageField(w, fi)
	}
	w.WriteString("} \n")
}

func WriteProtobufMessageField(w *bufio.Writer, f *model.ProtobufMessageField) {
	switch cf := f.GetField().(type) {
	case *model.ProtobufMessageField_BasicField:
		WriteProtobufMessageBasicField(w, cf.BasicField)
	case *model.ProtobufMessageField_MapField:
		WriteProtobufMessageMapField(w, cf.MapField)
	case *model.ProtobufMessageField_RepeatedField:
		WriteProtobufMessageRepeatedField(w, cf.RepeatedField)
	case *model.ProtobufMessageField_OneofField:
		WriteProtobufMessageOneofField(w, cf.OneofField)
	}
}

func WriteProtobufMessageBasicField(w *bufio.Writer, f *model.ProtobufMessageBasicField) {
	w.WriteString(fmt.Sprintf("%s %s = %d ;\n", f.GetTypeName(), f.GetName(), f.GetIndex()))
}

func WriteProtobufMessageMapField(w *bufio.Writer, f *model.ProtobufMessageMapField) {
	w.WriteString(fmt.Sprintf("map< %s , %s > %s = %d ;\n", f.GetKeyTypeName(), f.GetValueTypeName(), f.GetName(), f.GetIndex()))
}

func WriteProtobufMessageRepeatedField(w *bufio.Writer, f *model.ProtobufMessageRepeatedField) {
	w.WriteString(fmt.Sprintf("repeated %s %s = %d ;\n", f.GetTypeName(), f.GetName(), f.GetIndex()))
}

func WriteProtobufMessageOneofField(w *bufio.Writer, f *model.ProtobufMessageOneofField) {
	w.WriteString(fmt.Sprintf("oneof %s { \n", f.GetName()))
	for _, fi := range f.GetFields() {
		WriteProtobufMessageField(w, fi)
	}
	w.WriteString("} \n")
}
