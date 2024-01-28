package write

import (
	"bufio"
	"fmt"

	"github.com/hjwalt/platform/model"
)

func WriteProtobufSchema(w *bufio.Writer, f *model.ProtobufSchema) {
	w.WriteString("syntax = \"proto3\";\n")
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf("package %s ; \n", f.GetPackage()))
	w.WriteString(fmt.Sprintf("option go_package = \"%s\" ; \n", f.GetGoPackage()))
	w.WriteString("\n")
	for _, fi := range f.GetTypes() {
		WriteProtobufType(w, fi)
		w.WriteString("\n")
	}
}

func WriteProtobufType(w *bufio.Writer, f *model.ProtobufType) {
	switch cf := f.GetType().(type) {
	case *model.ProtobufType_Message:
		WriteProtobufMessage(w, cf.Message)
	}
}
