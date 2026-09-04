package harness

import (
	"testing"

	"github.com/hjwalt/platform/agent"
	"github.com/stretchr/testify/assert"
)

func TestSetContextReturnsUpdatedCopy(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{}
	out := st.SetContext("ctx-1")

	// original is untouched (value receiver semantics)
	assert.Empty(st.Context)
	assert.Equal("ctx-1", out.Context)
}

func TestSetAgentContext(t *testing.T) {
	assert := assert.New(t)

	parent := agent.AgentContext{
		ParentContext: "parent-1",
		SystemMessage: "sys",
		AllowedTools:  []string{"web_search"},
	}
	tool := agent.ToolCall{Id: "t1", Name: "web_search", Arguments: "{}"}

	st := ExecutionState{}
	out := st.SetAgentContext(parent, tool)

	assert.Empty(st.Parent.ParentContext)
	assert.Empty(st.ParentToolCall.Name)
	assert.Equal(parent, out.Parent)
	assert.Equal(tool, out.ParentToolCall)
}

func TestSetNext(t *testing.T) {
	assert := assert.New(t)

	result := agent.SingleResult(agent.NewMessage("c", agent.MessageType_Agent, "m", "", agent.ToolCall{}))

	st := ExecutionState{}
	out := st.SetNext(result)

	assert.Nil(st.Next.Messages, "original untouched")
	assert.Equal(result, out.Next)
}

func TestUpdateToolStateInitialisesNilMap(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{}
	out := st.UpdateToolState("web_search", ToolState_Requested)

	assert.Nil(st.ToolStates, "original untouched")
	assert.NotNil(out.ToolStates)
	assert.Equal(ToolState_Requested, out.ToolStates["web_search"])
}

func TestUpdateToolStateAddsAndOverwrites(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{ToolStates: map[string]ToolState{"a": ToolState_Requested}}
	out := st.UpdateToolState("b", ToolState_Executed)
	out = out.UpdateToolState("a", ToolState_Rejected)

	assert.Equal(ToolState_Rejected, out.ToolStates["a"])
	assert.Equal(ToolState_Executed, out.ToolStates["b"])
}

func TestAppendMessageInitialisesNilSlice(t *testing.T) {
	assert := assert.New(t)

	msg := agent.NewMessage("c", agent.MessageType_User, "hi", "", agent.ToolCall{})

	st := ExecutionState{}
	out := st.AppendMessage(msg)

	assert.Nil(st.Messages, "original untouched")
	assert.NotNil(out.Messages)
	assert.Len(out.Messages, 1)
	assert.Equal(msg, out.Messages[0])
}

func TestAppendMessageGrowsSlice(t *testing.T) {
	assert := assert.New(t)

	m1 := agent.NewMessage("c", agent.MessageType_User, "one", "", agent.ToolCall{})
	m2 := agent.NewMessage("c", agent.MessageType_Agent, "two", "", agent.ToolCall{})

	out := ExecutionState{}.AppendMessage(m1).AppendMessage(m2)

	assert.Len(out.Messages, 2)
	assert.Equal(m1, out.Messages[0])
	assert.Equal(m2, out.Messages[1])
}

func TestSkillLoadedNilMapReturnsFalse(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{}
	assert.False(st.SkillLoaded("anything"))
	assert.False(st.SkillLoaded(""))
}

func TestSkillLoadedNotPresentReturnsFalse(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{LoadedSkills: map[string]bool{"a": true}}
	assert.False(st.SkillLoaded("b"))
}

func TestAppendSkillLoadedMarksLoaded(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{}
	out := st.AppendSkillLoaded("my-skill")

	assert.Nil(st.LoadedSkills, "original untouched")
	assert.NotNil(out.LoadedSkills)
	assert.True(out.SkillLoaded("my-skill"))
}

func TestSkillLoadedIsCaseInsensitive(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{}.AppendSkillLoaded("MySkill")

	assert.True(st.SkillLoaded("myskill"), "lookup lowercases the loaded name")
	assert.True(st.SkillLoaded("MYSKILL"))
	assert.True(st.SkillLoaded("MySkill"))
}

func TestSkillLoadedStoredLowercased(t *testing.T) {
	assert := assert.New(t)

	st := ExecutionState{}
	out := st.AppendSkillLoaded("UpperCaseSkill")

	_, present := out.LoadedSkills["uppercaseskill"]
	assert.True(present, "stored key should be lowercased")
	_, raw := out.LoadedSkills["UpperCaseSkill"]
	assert.False(raw)
}
