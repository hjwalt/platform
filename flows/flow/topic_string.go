package flow

import "github.com/hjwalt/platform/format"

func StringTopic(topic string) Topic[string, string] {
	return GenericTopic(topic, format.String(), format.String())
}
