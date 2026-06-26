package skill_tool

import (
	"context"
	"strings"
	"testing"

	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
	tool_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/stretchr/testify/assert"
)

// testRegistry returns a skill registry with two skills for testing.
func testRegistry() map[string]agent_skill.Skill {
	return map[string]agent_skill.Skill{
		"code-review": {
			Name:        "code-review",
			Description: "Review code for bugs and improvements",
			Body:        "## Code Review Playbook\n\n1. Check for correctness\n2. Check for performance\n",
			AllowedTools: []string{"linux_shell", "web_fetch"},
		},
		"deep-research": {
			Name:        "deep-research",
			Description: "Perform deep research on a topic",
			Body:        "## Deep Research Playbook\n\n1. Search the web\n2. Fetch sources\n3. Synthesize findings\n",
			AllowedTools: []string{"web_search", "web_fetch"},
		},
	}
}

// --- FR-SKILL-002, FR-SKILL-003: Successful skill lookup ---

func TestSkillLookupReturnsFullContent(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	skill := Create(testRegistry())

	resp, err := skill.Apply(ctx, Request{Name: "code-review"})
	assert.NoError(err)
	assert.True(resp.Found)
	assert.Equal("code-review", resp.Name)
	assert.Equal("Review code for bugs and improvements", resp.Description)
	assert.Contains(resp.Body, "Code Review Playbook")
	assert.Equal([]string{"linux_shell", "web_fetch"}, resp.AllowedTools)
}

func TestSkillLookupCaseInsensitive(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	skill := Create(testRegistry())

	for _, name := range []string{"CODE-REVIEW", "Code-Review", "code-review", "  code-review  "} {
		resp, err := skill.Apply(ctx, Request{Name: name})
		assert.NoError(err, "failed for name: %q", name)
		assert.True(resp.Found, "should be found for name: %q", name)
		assert.Equal("code-review", resp.Name)
	}
}

func TestSkillLookupMultipleSkills(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	skill := Create(testRegistry())

	// Look up first skill
	resp1, err1 := skill.Apply(ctx, Request{Name: "code-review"})
	assert.NoError(err1)
	assert.True(resp1.Found)
	assert.Contains(resp1.Body, "Code Review Playbook")

	// Look up second skill
	resp2, err2 := skill.Apply(ctx, Request{Name: "deep-research"})
	assert.NoError(err2)
	assert.True(resp2.Found)
	assert.Contains(resp2.Body, "Deep Research Playbook")
}

// --- FR-SKILL-004: Missing skill handling ---

func TestSkillNotFound(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	skill := Create(testRegistry())

	resp, err := skill.Apply(ctx, Request{Name: "nonexistent-skill"})
	assert.NoError(err)
	assert.False(resp.Found)
	assert.Equal("nonexistent-skill", resp.Name) // requested name preserved for error reporting
	assert.Empty(resp.Body)
}

func TestSkillEmptyName(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	skill := Create(testRegistry())

	resp, err := skill.Apply(ctx, Request{Name: ""})
	assert.NoError(err)
	assert.False(resp.Found)

	resp2, err2 := skill.Apply(ctx, Request{Name: "   "})
	assert.NoError(err2)
	assert.False(resp2.Found)
}

// --- FR-SKILL-005: Auto policy ---

func TestSkillAutoPolicy(t *testing.T) {
	assert := assert.New(t)

	skill := Create(testRegistry())
	assert.True(skill.Auto())
}

// --- FR-SKILL-006: Metadata and schemas ---

func TestSkillMetadata(t *testing.T) {
	assert := assert.New(t)

	skill := Create(testRegistry())

	assert.Equal("skill", skill.Name())
	assert.NotEmpty(skill.Description())
	assert.NotNil(skill.RequestSchema())
	assert.NotNil(skill.ResultSchema())

	// RequestFormat and ResultFormat should not be nil
	assert.NotNil(skill.RequestFormat())
	assert.NotNil(skill.ResultFormat())
}

func TestSkillDescribeRequest(t *testing.T) {
	assert := assert.New(t)

	skill := Create(testRegistry())
	desc := skill.DescribeRequest(Request{Name: "code-review"})
	assert.Contains(desc, "code-review")
	assert.Contains(desc, "loading skill")
}

func TestSkillDescribeResultFound(t *testing.T) {
	assert := assert.New(t)

	skill := Create(testRegistry())
	desc := skill.DescribeResult(Response{
		Name:        "code-review",
		Description: "Review code for bugs",
		Body:        "## Playbook\n\nStep 1",
		Found:       true,
	})

	assert.Contains(desc, "loaded skill")
	assert.Contains(desc, "code-review")
	assert.Contains(desc, "Review code for bugs")
	assert.Contains(desc, "## Playbook")
}

func TestSkillDescribeResultNotFound(t *testing.T) {
	assert := assert.New(t)

	skill := Create(testRegistry())
	desc := skill.DescribeResult(Response{
		Name:  "missing-skill",
		Found: false,
	})

	assert.Contains(desc, "missing-skill")
	assert.Contains(desc, "not found")
}

// --- FR-SKILL-007: Empty registry handling ---

func TestSkillEmptyRegistry(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	skill := Create(make(map[string]agent_skill.Skill))

	resp, err := skill.Apply(ctx, Request{Name: "any-skill"})
	assert.NoError(err)
	assert.False(resp.Found)
}

func TestSkillNilRegistry(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// Should initialize empty map internally, not panic
	skill := Create(nil)

	resp, err := skill.Apply(ctx, Request{Name: "any-skill"})
	assert.NoError(err)
	assert.False(resp.Found)
}

// --- FR-SKILL-008: Registration ---

func TestAddToContainerRegistersSkillTool(t *testing.T) {
	assert := assert.New(t)
	container := tool_container.New()

	AddToContainer(container, testRegistry())

	assert.True(container.Exists(agent.ToolCall{Name: "skill"}))
}

// --- Integration: DescribeResult contains full body for context injection ---

func TestDescribeResultContainsFullBody(t *testing.T) {
	assert := assert.New(t)

	skill := Create(testRegistry())
	desc := skill.DescribeResult(Response{
		Name:        "code-review",
		Description: "Review code for bugs",
		Body:        "## Step 1\nDo X\n\n## Step 2\nDo Y\n",
		Found:       true,
	})

	// The body must be fully included so the LLM gets the complete playbook
	assert.True(strings.Contains(desc, "Do X"), "DescribeResult must include full body content")
	assert.True(strings.Contains(desc, "Do Y"), "DescribeResult must include full body content")
}
