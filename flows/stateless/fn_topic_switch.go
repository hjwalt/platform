package stateless

import (
	"context"
	"errors"

	"github.com/hjwalt/platform/commons/logger"
	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/flows/flow"
	"go.uber.org/zap"
)

// constructor
func NewTopicSwitch(configurations ...runtime.Configuration[*TopicSwitch]) BatchFunction {
	fn := &TopicSwitch{
		functions: make(map[string]BatchFunction),
	}
	for _, configuration := range configurations {
		fn = configuration(fn)
	}
	return fn.Apply
}

// configuration
func WithTopicSwitchFunction(topic string, f BatchFunction) runtime.Configuration[*TopicSwitch] {
	return func(sts *TopicSwitch) *TopicSwitch {
		sts.functions[topic] = f
		return sts
	}
}

// implementation
type TopicSwitch struct {
	functions map[string]BatchFunction
}

func (r *TopicSwitch) Apply(c context.Context, m []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {

	messageMultiMap := structure.NewMultiMap[string, flow.Message[structure.Bytes, structure.Bytes]]()
	for _, mi := range m {
		messageMultiMap.Add(mi.Topic, mi)
	}

	resultMessages := []flow.Message[structure.Bytes, structure.Bytes]{}
	for k, v := range messageMultiMap.GetAll() {
		logger.Info("stateless switch", zap.String("topic", k), zap.Int("messages", len(v)))

		fn, fnExists := r.functions[k]
		if !fnExists {
			return make([]flow.Message[[]byte, []byte], 0), errors.Join(errors.New(k), ErrSwitchMissingTopic)
		}

		currGroupMessages, currGroupHandlerErr := fn(c, v)
		if currGroupHandlerErr != nil {
			return make([]flow.Message[[]byte, []byte], 0), currGroupHandlerErr
		}
		resultMessages = append(resultMessages, currGroupMessages...)
	}

	return resultMessages, nil
}

var (
	ErrSwitchMissingTopic = errors.New("stateless switch missing topic")
)
