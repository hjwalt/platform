package harness

import (
	"context"
	"errors"
	"testing"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/hjwalt/platform/agent"
	skill_tool "github.com/hjwalt/platform/agent/tool/skill"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/type/either"
	"github.com/hjwalt/platform/type/optional"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

// mockToolContainer is a behaviour-configurable agent.ToolContainer.
type mockToolContainer struct {
	exists map[string]bool
	auto   map[string]bool

	execResult optional.Optional[string]
	execErr    error
	execCalls  int
}

func newMockToolContainer() *mockToolContainer {
	return &mockToolContainer{
		exists:     make(map[string]bool),
		auto:       make(map[string]bool),
		execResult: optional.Empty[string](),
	}
}

func (m *mockToolContainer) AddSync(agent.SyncToolWrapper)   {}
func (m *mockToolContainer) AddAsync(agent.AsyncToolWrapper) {}

func (m *mockToolContainer) Execute(ctx context.Context, in agent.Message, call agent.ToolCall) (optional.Optional[string], error) {
	m.execCalls++
	return m.execResult, m.execErr
}

func (m *mockToolContainer) DescribeRequest(call agent.ToolCall) (string, error) {
	return "", nil
}

func (m *mockToolContainer) Exists(call agent.ToolCall) bool {
	return m.exists[call.Name]
}

func (m *mockToolContainer) Auto(call agent.ToolCall) bool {
	return m.auto[call.Name]
}

func (m *mockToolContainer) OpenAiParamsFiltered(allowed []string) []openai.ChatCompletionToolUnionParam {
	return nil
}

func (m *mockToolContainer) DeepSeekParams(allowed []string) []deepseek.Tool {
	return nil
}

// mockSkillContainer returns a fixed assistant message.
type mockSkillContainer struct {
	assistant agent.Message
}

func (m *mockSkillContainer) Add(agent.Instruction) {}
func (m *mockSkillContainer) Get(string) (agent.Instruction, bool) {
	return agent.Instruction{}, false
}
func (m *mockSkillContainer) Assistant(ctx string) agent.Message {
	return m.assistant
}

// mockLanguageModel is a behaviour-configurable agent.LanguageModel.
type mockLanguageModel struct {
	chatResult []agent.Message
	chatErr    error

	calls       int
	gotMessages []agent.Message
	gotTools    []string
}

func (m *mockLanguageModel) Start() error { return nil }
func (m *mockLanguageModel) Stop()        {}

func (m *mockLanguageModel) Chat(ctx context.Context, messages []agent.Message, allowedTools []string) ([]agent.Message, error) {
	m.calls++
	m.gotMessages = append([]agent.Message(nil), messages...)
	m.gotTools = append([]string(nil), allowedTools...)
	return m.chatResult, m.chatErr
}

// ---------------------------------------------------------------------------
// Test harness
// ---------------------------------------------------------------------------

type flowFixture struct {
	flow       *Flow
	tools      *mockToolContainer
	skills     *mockSkillContainer
	model      *mockLanguageModel
	assistant  agent.Message
	background context.Context
}

func newFlowFixture(t *testing.T) *flowFixture {
	t.Helper()

	assistant := agent.NewMessage("ctx-1", agent.MessageType_System, "available skills: none", "", agent.ToolCall{})

	fixture := &flowFixture{
		tools:      newMockToolContainer(),
		skills:     &mockSkillContainer{assistant: assistant},
		model:      &mockLanguageModel{},
		assistant:  assistant,
		background: context.Background(),
	}
	fixture.flow = &Flow{
		Tools:  fixture.tools,
		Skills: fixture.skills,
		Model:  fixture.model,
	}
	return fixture
}

func newUserMessage(ctx string, text string) agent.Message {
	return agent.NewMessage(ctx, agent.MessageType_User, text, "", agent.ToolCall{})
}

func assertLeft(t *testing.T, result either.Either[ExecutionState, agent.Message]) ExecutionState {
	t.Helper()
	require.True(t, result.IsLeft())
	return result.Left()
}

// ---------------------------------------------------------------------------
// FlowMetadata
// ---------------------------------------------------------------------------

func TestFlowMetadataKey(t *testing.T) {
	assert := assert.New(t)

	r := &FlowMetadata{}

	key, err := r.Key(context.Background(), agent.NewMessage("ctx-42", agent.MessageType_User, "hi", "", agent.ToolCall{}))
	assert.NoError(err)
	assert.Equal("ctx-42", key)
}

func TestFlowMetadataKeyDefaultsWhenContextEmpty(t *testing.T) {
	assert := assert.New(t)

	r := &FlowMetadata{}
	msg := agent.NewMessage("", agent.MessageType_User, "hi", "", agent.ToolCall{})

	key, err := r.Key(context.Background(), msg)
	assert.NoError(err)
	assert.Equal("DEFAULT", key)
}

func TestFlowMetadataResultMetadata(t *testing.T) {
	assert := assert.New(t)

	r := &FlowMetadata{}
	pref := flow.Metadata{Id: "prev", Group: "group-1", Attempt: 3, Sequence: 7, Source: "other"}
	value := agent.SingleResult(agent.NewMessage("ctx", agent.MessageType_Agent, "m", "", agent.ToolCall{}))

	out := r.ResultMetadata(context.Background(), pref, value)

	assert.Equal(value.Id, out.Id)
	assert.Equal("group-1", out.Group)
	assert.Equal(int32(0), out.Attempt)
	assert.Equal(int64(8), out.Sequence)
	assert.Equal("AGENT_HARNESS", out.Source)
}

func TestFlowMetadataMessageMetadata(t *testing.T) {
	assert := assert.New(t)

	r := &FlowMetadata{}
	pref := flow.Metadata{Id: "prev", Group: "g", Attempt: 2, Sequence: 5, Source: "other"}
	value := agent.NewMessage("ctx-9", agent.MessageType_User, "m", "", agent.ToolCall{})

	out := r.MessageMetadata(context.Background(), pref, value)

	assert.Equal(value.Id, out.Id)
	assert.Equal("ctx-9", out.Group)
	assert.Equal(int32(0), out.Attempt)
	assert.Equal(int64(6), out.Sequence)
	assert.Equal("AGENT_HARNESS", out.Source)
}

// ---------------------------------------------------------------------------
// Flow.Update - MessageType_Start
// ---------------------------------------------------------------------------

func TestUpdateStartAgent(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	parent := agent.AgentContext{
		ParentContext: "parent-1",
		SystemMessage: "you are a helpful agent",
		AllowedTools:  []string{"web_search"},
	}
	in := agent.Start("ctx-1", "do the thing", agent.ToolCall{Id: "t0", Name: "start-tool"}, parent)

	result := fixture.flow.Update(fixture.background, in, ExecutionState{})
	st := assertLeft(t, result)

	// context copied from incoming message
	assert.Equal("ctx-1", st.Context)
	// parent context and tool call stored
	assert.Equal(parent, st.Parent)
	assert.Equal(in.Tool, st.ParentToolCall)
	// incoming message appended first, then a System message from parent
	assert.Len(st.Messages, 2)
	assert.Equal(in, st.Messages[0])
	assert.Equal(agent.MessageType_System, st.Messages[1].Type)
	assert.Equal(parent.SystemMessage, st.Messages[1].Message)
	assert.Equal("ctx-1", st.Messages[1].Context)
	// next contains a single USER message carrying the original request
	require.Len(t, st.Next.Messages, 1)
	assert.Equal(agent.MessageType_User, st.Next.Messages[0].Type)
	assert.Equal(in.Message, st.Next.Messages[0].Message)
	assert.Equal(agent.ToolCall{}, st.Next.Messages[0].Tool)
}

func TestUpdateStartAgentWithEmptySystemMessage(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	in := agent.Start("ctx-1", "msg", agent.ToolCall{}, agent.AgentContext{})

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	require.Len(t, st.Messages, 2)
	assert.Equal(agent.MessageType_System, st.Messages[1].Type)
	assert.Empty(st.Messages[1].Message)
}

func TestUpdateKeepsExistingContext(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	in := agent.NewMessage("incoming-ctx", agent.MessageType_Agent, "m", "", agent.ToolCall{})
	pre := ExecutionState{Context: "pre-set-ctx"}

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, pre))

	assert.Equal("pre-set-ctx", st.Context)
}

// ---------------------------------------------------------------------------
// Flow.Update - MessageType_Agent (mergeMessage)
// ---------------------------------------------------------------------------

func TestUpdateMergeMessageWithoutParentProducesEmptyNext(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	in := agent.NewMessage("ctx-1", agent.MessageType_Agent, "agent reply", "", agent.ToolCall{})

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Empty(st.Next.Messages, "no parent context means nothing is merged upstream")
	require.Len(t, st.Messages, 1)
	assert.Equal(in, st.Messages[0])
}

func TestUpdateMergeMessageWithParentForwardsToolResult(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	parentTool := agent.ToolCall{Id: "tc-7", Name: "subagent"}
	parent := agent.AgentContext{ParentContext: "parent-ctx"}
	st := ExecutionState{Parent: parent, ParentToolCall: parentTool}

	in := agent.NewMessage("ctx-1", agent.MessageType_Agent, "child result", "reasoning", agent.ToolCall{})

	out := assertLeft(t, fixture.flow.Update(fixture.background, in, st))

	require.Len(t, out.Next.Messages, 1)
	next := out.Next.Messages[0]
	assert.Equal(agent.MessageType_ToolResult, next.Type)
	assert.Equal("parent-ctx", next.Context)
	assert.Equal("child result", next.Message)
	assert.Equal("reasoning", next.ReasoningContent)
	assert.Equal(parentTool, next.Tool)
}

// ---------------------------------------------------------------------------
// Flow.Update - MessageType_ToolRequest
// ---------------------------------------------------------------------------

func TestUpdateToolRequestAutoToolProceedsToExecute(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.tools.exists["web_search"] = true
	fixture.tools.auto["web_search"] = true

	tool := agent.ToolCall{Id: "tr-1", Name: "web_search", Arguments: `{"term":"go"}`}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolRequest, "search for go", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Equal(ToolState_Requested, st.ToolStates[tool.Id])
	require.Len(t, st.Next.Messages, 1)
	assert.Equal(agent.MessageType_ToolExecute, st.Next.Messages[0].Type)
	assert.Equal("execution approved to search for go", st.Next.Messages[0].Message)
	assert.Equal(tool, st.Next.Messages[0].Tool)
}

func TestUpdateToolRequestToolExistsButNotAutoWaits(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.tools.exists["web_search"] = true
	fixture.tools.auto["web_search"] = false

	tool := agent.ToolCall{Id: "tr-2", Name: "web_search"}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolRequest, "search", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Equal(ToolState_Requested, st.ToolStates[tool.Id])
	assert.Empty(st.Next.Messages, "non-auto tools do not emit an approval message")
}

func TestUpdateToolRequestUnknownToolErrors(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	tool := agent.ToolCall{Id: "tr-3", Name: "ghost_tool"}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolRequest, "run it", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Equal(ToolState_Failed, st.ToolStates[tool.Id])
	require.Len(t, st.Next.Messages, 1)
	next := st.Next.Messages[0]
	assert.Equal(agent.MessageType_Error, next.Type)
	assert.Equal("tool ghost_tool does not exist", next.Message)
	assert.Equal(tool, next.Tool)
}

// ---------------------------------------------------------------------------
// Flow.Update - MessageType_ToolExecute
// ---------------------------------------------------------------------------

func TestUpdateToolExecuteAlreadyExecutedErrors(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	tool := agent.ToolCall{Id: "te-1", Name: "web_search"}
	st := ExecutionState{ToolStates: map[string]ToolState{tool.Id: ToolState_Executed}}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "exec", "", tool)

	out := assertLeft(t, fixture.flow.Update(fixture.background, in, st))

	require.Len(t, out.Next.Messages, 1)
	next := out.Next.Messages[0]
	assert.Equal(agent.MessageType_Error, next.Type)
	assert.Equal("tool web_search already executed", next.Message)
	assert.Equal(tool, next.Tool)
	assert.Equal(0, fixture.tools.execCalls, "Execute must not be invoked for an executed tool")
}

func TestUpdateToolExecuteRejectedToolErrors(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	tool := agent.ToolCall{Id: "te-2", Name: "web_search"}
	st := ExecutionState{ToolStates: map[string]ToolState{tool.Id: ToolState_Rejected}}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "exec", "", tool)

	out := assertLeft(t, fixture.flow.Update(fixture.background, in, st))

	require.Len(t, out.Next.Messages, 1)
	assert.Equal("tool web_search already executed", out.Next.Messages[0].Message)
	assert.Equal(0, fixture.tools.execCalls)
}

func TestUpdateToolExecuteSkillToolParseFailure(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	// reset any tool execution response so the test fails loudly if reached
	fixture.tools.execErr = errors.New("execute should not be reached")

	tool := agent.ToolCall{Id: "te-3", Name: skill_tool.Name, Arguments: `{not json`}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "load skill", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	_, parseErr := skill_tool.Parse(tool.Arguments)
	require.Error(t, parseErr)

	assert.Equal(ToolState_Failed, st.ToolStates[tool.Id])
	require.Len(t, st.Next.Messages, 1)
	assert.Equal(agent.MessageType_Error, st.Next.Messages[0].Type)
	assert.Equal(parseErr.Error(), st.Next.Messages[0].Message)
	assert.Equal(0, fixture.tools.execCalls)
}

func TestUpdateToolExecuteSkillToolAlreadyLoaded(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	tool := agent.ToolCall{Id: "te-4", Name: skill_tool.Name, Arguments: `{"name":"MySkill"}`}
	st := ExecutionState{LoadedSkills: map[string]bool{"myskill": true}}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "load skill", "", tool)

	out := assertLeft(t, fixture.flow.Update(fixture.background, in, st))

	assert.Equal(ToolState_Executed, out.ToolStates[tool.Id])
	require.Len(t, out.Next.Messages, 1)
	next := out.Next.Messages[0]
	assert.Equal(agent.MessageType_ToolResult, next.Type)
	assert.Equal("skill MySkill is already loaded", next.Message)
	assert.Equal(0, fixture.tools.execCalls)
}

func TestUpdateToolExecuteSkillToolNotLoadedExecutes(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.tools.execResult = optional.Of("loaded skill `myskill`")

	tool := agent.ToolCall{Id: "te-5", Name: skill_tool.Name, Arguments: `{"name":"myskill"}`}
	st := ExecutionState{}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "load skill", "", tool)

	out := assertLeft(t, fixture.flow.Update(fixture.background, in, st))

	assert.Equal(1, fixture.tools.execCalls)
	assert.Equal(ToolState_Executed, out.ToolStates[tool.Id])
	require.Len(t, out.Next.Messages, 1)
	assert.Equal(agent.MessageType_ToolResult, out.Next.Messages[0].Type)
	assert.Equal("loaded skill `myskill`", out.Next.Messages[0].Message)
}

func TestUpdateToolExecuteSuccessWithResult(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.tools.execResult = optional.Of("the result output")

	tool := agent.ToolCall{Id: "te-6", Name: "web_search", Arguments: `{"term":"go"}`}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "exec", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Equal(1, fixture.tools.execCalls)
	assert.Equal(ToolState_Executed, st.ToolStates[tool.Id])
	require.Len(t, st.Next.Messages, 1)
	next := st.Next.Messages[0]
	assert.Equal(agent.MessageType_ToolResult, next.Type)
	assert.Equal("the result output", next.Message)
	assert.Equal(tool, next.Tool)
}

func TestUpdateToolExecuteToolError(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.tools.execResult = optional.Empty[string]()
	fixture.tools.execErr = errors.New("boom")

	tool := agent.ToolCall{Id: "te-7", Name: "web_search"}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "exec", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Equal(ToolState_Failed, st.ToolStates[tool.Id])
	require.Len(t, st.Next.Messages, 1)
	assert.Equal(agent.MessageType_Error, st.Next.Messages[0].Type)
	assert.Equal("boom", st.Next.Messages[0].Message)
}

func TestUpdateToolExecuteEmptyResultNoErrorProducesEmptyNext(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.tools.execResult = optional.Empty[string]()
	fixture.tools.execErr = nil

	tool := agent.ToolCall{Id: "te-8", Name: "async_tool"}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "exec", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Equal(ToolState_Executed, st.ToolStates[tool.Id])
	assert.Empty(st.Next.Messages)
}

// ---------------------------------------------------------------------------
// Flow.Update - MessageType_ToolResult
// ---------------------------------------------------------------------------

func TestUpdateToolResultForwardsToModel(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	modelReply := agent.NewMessage("ctx-1", agent.MessageType_Agent, "model says done", "", agent.ToolCall{})
	fixture.model.chatResult = []agent.Message{modelReply}

	tool := agent.ToolCall{Id: "tr-9", Name: "web_search"}
	in := agent.NewMessage("ctx-1", agent.MessageType_ToolResult, "tool outcome", "", tool)

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	assert.Equal(ToolState_Executed, st.ToolStates[tool.Id])
	require.Len(t, st.Next.Messages, 1)
	assert.Equal(modelReply, st.Next.Messages[0])
	// the executed tool result message is included in the model conversation
	require.Equal(t, 1, fixture.model.calls)
	require.Len(t, fixture.model.gotMessages, 1)
	assert.Equal(in, fixture.model.gotMessages[0])
}

// ---------------------------------------------------------------------------
// Flow.Update - default / unknown type
// ---------------------------------------------------------------------------

func TestUpdateUnknownMessageTypeReturnsStateUnchanged(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	pre := agent.NewMessage("ctx-1", agent.MessageType_Agent, "old", "", agent.ToolCall{})
	st := ExecutionState{
		Context:  "ctx-1",
		Messages: []agent.Message{pre},
		Next:     agent.SingleResult(pre),
	}

	in := agent.NewMessage("ctx-1", agent.MessageType_System, "unhandled", "", agent.ToolCall{})

	out := assertLeft(t, fixture.flow.Update(fixture.background, in, st))

	// message appended, next reset to empty, nothing else produced
	require.Len(t, out.Messages, 2)
	assert.Equal(pre, out.Messages[0])
	assert.Equal(in, out.Messages[1])
	assert.Empty(out.Next.Messages)
	assert.Equal(0, fixture.model.calls)
}

func TestUpdateResetsNextOnEveryCall(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	first := agent.NewMessage("ctx-1", agent.MessageType_Agent, "first", "", agent.ToolCall{})
	second := agent.NewMessage("ctx-1", agent.MessageType_Agent, "second", "", agent.ToolCall{})

	st := assertLeft(t, fixture.flow.Update(fixture.background, first, ExecutionState{}))
	// seed a Next to prove the next Update clears it
	st.Next = agent.SingleResult(first)
	st = assertLeft(t, fixture.flow.Update(fixture.background, second, st))

	assert.Empty(st.Next.Messages, "Next is reset to an empty result each update")
	require.Len(t, st.Messages, 2)
	assert.Equal(first, st.Messages[0])
	assert.Equal(second, st.Messages[1])
}

// ---------------------------------------------------------------------------
// Flow.Update - MessageType_User (modelExecute)
// ---------------------------------------------------------------------------

func TestUpdateUserMessageCallsModelAndStoresResult(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	reply := agent.NewMessage("ctx-1", agent.MessageType_Agent, "hello there", "", agent.ToolCall{})
	fixture.model.chatResult = []agent.Message{reply}

	in := newUserMessage("ctx-1", "hi")
	st := ExecutionState{Parent: agent.AgentContext{AllowedTools: []string{"web_search"}}}

	out := assertLeft(t, fixture.flow.Update(fixture.background, in, st))

	require.Len(t, out.Next.Messages, 1)
	assert.Equal(reply, out.Next.Messages[0])
	assert.Equal(1, fixture.model.calls)
	assert.Equal([]string{"web_search"}, fixture.model.gotTools)
}

func TestUpdateUserMessagePrependsAssistantWhenFirstMessageIsDefault(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.model.chatResult = []agent.Message{agent.NewMessage("ctx-1", agent.MessageType_Agent, "r", "", agent.ToolCall{})}

	in := newUserMessage("ctx-1", "hi")

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	require.Len(t, st.Next.Messages, 1)
	require.Len(t, fixture.model.gotMessages, 2)
	assert.Equal(fixture.assistant, fixture.model.gotMessages[0],
		"skills assistant message is prepended before the first user message")
	assert.Equal(in, fixture.model.gotMessages[1])
}

func TestUpdateUserMessagePrependsAssistantAfterFirstSystemMessage(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.model.chatResult = []agent.Message{agent.NewMessage("ctx-1", agent.MessageType_Agent, "r", "", agent.ToolCall{})}

	systemMsg := agent.NewMessage("ctx-1", agent.MessageType_System, "system prompt", "", agent.ToolCall{})
	in := newUserMessage("ctx-1", "hi")
	pre := ExecutionState{Messages: []agent.Message{systemMsg}}

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, pre))

	require.Len(t, st.Next.Messages, 1)
	require.Len(t, fixture.model.gotMessages, 3)
	assert.Equal(systemMsg, fixture.model.gotMessages[0])
	assert.Equal(fixture.assistant, fixture.model.gotMessages[1],
		"skills assistant message is prepended right after the leading system message")
	assert.Equal(in, fixture.model.gotMessages[2])
}

func TestUpdateUserMessageOnlyPassesExecutedToolMessages(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.model.chatResult = []agent.Message{agent.NewMessage("ctx-1", agent.MessageType_Agent, "r", "", agent.ToolCall{})}

	executed := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "executed", "", agent.ToolCall{Id: "t1"})
	requested := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "requested", "", agent.ToolCall{Id: "t2"})
	failed := agent.NewMessage("ctx-1", agent.MessageType_ToolResult, "failed", "", agent.ToolCall{Id: "t3"})
	in := newUserMessage("ctx-1", "hi")

	pre := ExecutionState{
		Messages: []agent.Message{executed, requested, failed},
		ToolStates: map[string]ToolState{
			"t1": ToolState_Executed,
			"t2": ToolState_Requested,
			"t3": ToolState_Failed,
		},
	}

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, pre))

	require.Len(t, st.Next.Messages, 1)
	require.Len(t, fixture.model.gotMessages, 2)
	assert.Equal(executed, fixture.model.gotMessages[0])
	assert.Equal(in, fixture.model.gotMessages[1])
}

func TestUpdateUserMessageChatErrorProducesErrorResult(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	fixture.model.chatErr = errors.New("upstream failure")

	in := newUserMessage("ctx-1", "hi")

	st := assertLeft(t, fixture.flow.Update(fixture.background, in, ExecutionState{}))

	require.Len(t, st.Next.Messages, 1)
	next := st.Next.Messages[0]
	assert.Equal(agent.MessageType_Error, next.Type)
	assert.Equal("upstream failure", next.Message)
	assert.Equal(1, fixture.model.calls)
}

// ---------------------------------------------------------------------------
// Flow.Next / Flow.Explode
// ---------------------------------------------------------------------------

func TestFlowNextReturnsCurrentResult(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	result := agent.SingleResult(agent.NewMessage("ctx-1", agent.MessageType_Agent, "m", "", agent.ToolCall{}))
	st := ExecutionState{Next: result}

	optionalResult, optionalMessage := fixture.flow.Next(fixture.background, agent.Message{}, st)

	assert.True(optionalResult.IsPresent())
	assert.Equal(result, optionalResult.Get())
	assert.False(optionalMessage.IsPresent())
}

func TestFlowExplodeReturnsMessages(t *testing.T) {
	assert := assert.New(t)
	fixture := newFlowFixture(t)

	messages := []agent.Message{
		agent.NewMessage("ctx-1", agent.MessageType_Agent, "one", "", agent.ToolCall{}),
		agent.NewMessage("ctx-1", agent.MessageType_Agent, "two", "", agent.ToolCall{}),
	}
	result := agent.NewResult(messages)

	optionalMessages, optionalMessage := fixture.flow.Explode(fixture.background, result)

	assert.True(optionalMessages.IsPresent())
	assert.Equal(messages, optionalMessages.Get())
	assert.False(optionalMessage.IsPresent())
}
