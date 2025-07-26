package task

import "github.com/hjwalt/platform/commons/format"

func StringChannel(channel string) Channel[string] {
	return GenericChannel(channel, format.String())
}
