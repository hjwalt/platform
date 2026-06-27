package harness_container

import (
	"testing"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/type/optional"
	"github.com/stretchr/testify/assert"
)

func TestNewSkillContainerReturnsNonNil(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()
	assert.NotNil(c)
}

func TestAddAndGetRoundTripsAllFields(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()

	in := agent.Instruction{
		Name:          "code-review",
		Description:   "Review code for bugs and improvements",
		License:       optional.Of("MIT"),
		Compatibility: optional.Of(">=1.0.0"),
		AllowedTools:  []string{"linux_shell", "web_fetch"},
		Metadata:      map[string]string{"author": "team-platform"},
		Body:          "## Code Review Playbook\n\n1. Check correctness\n2. Check performance\n",
	}

	c.Add(in)

	out, ok := c.Get("code-review")
	assert.True(ok)
	assert.Equal(in.Name, out.Name)
	assert.Equal(in.Description, out.Description)
	assert.Equal(in.License, out.License)
	assert.Equal(in.Compatibility, out.Compatibility)
	assert.Equal(in.AllowedTools, out.AllowedTools)
	assert.Equal(in.Metadata, out.Metadata)
	assert.Equal(in.Body, out.Body)
}

func TestGetMissingSkillReturnsFalse(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()

	out, ok := c.Get("nonexistent")
	assert.False(ok)
	assert.Equal(agent.Instruction{}, out)
}

func TestGetEmptyNameReturnsFalse(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()

	c.Add(agent.Instruction{Name: "code-review", Description: "desc", Body: "body"})

	out, ok := c.Get("")
	assert.False(ok)
	assert.Equal(agent.Instruction{}, out)
}

func TestAddMultipleSkillsThenRetrieveAll(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()

	skills := []agent.Instruction{
		{Name: "code-review", Description: "Review code", Body: "## Review\n"},
		{Name: "deep-research", Description: "Research topics", Body: "## Research\n"},
		{Name: "simplify", Description: "Simplify code", Body: "## Simplify\n"},
	}

	for _, s := range skills {
		c.Add(s)
	}

	for _, s := range skills {
		out, ok := c.Get(s.Name)
		assert.True(ok, "skill %q should be found", s.Name)
		assert.Equal(s.Description, out.Description)
		assert.Equal(s.Body, out.Body)
	}
}

func TestAddOverwritesExistingSkill(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()

	original := agent.Instruction{
		Name:        "code-review",
		Description: "Original description",
		Body:        "## Original body\n",
	}
	updated := agent.Instruction{
		Name:        "code-review",
		Description: "Updated description",
		Body:        "## Updated body\n",
	}

	c.Add(original)
	c.Add(updated)

	out, ok := c.Get("code-review")
	assert.True(ok)
	assert.Equal("Updated description", out.Description)
	assert.Equal("## Updated body\n", out.Body)
}

func TestAddWithOptionalFieldsEmpty(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()

	in := agent.Instruction{
		Name:        "minimal-skill",
		Description: "A minimal skill",
		Body:        "## Minimal\n",
	}

	c.Add(in)

	out, ok := c.Get("minimal-skill")
	assert.True(ok)
	assert.False(out.License.IsPresent())
	assert.False(out.Compatibility.IsPresent())
	assert.Nil(out.AllowedTools)
	assert.Nil(out.Metadata)
}

func TestAddWithNilMetadataAndAllowedTools(t *testing.T) {
	assert := assert.New(t)
	c := NewSkillContainer()

	in := agent.Instruction{
		Name:         "nil-slices-skill",
		Description:  "Skill with nil slices",
		Body:         "## Nil slices\n",
		AllowedTools: nil,
		Metadata:     nil,
	}

	c.Add(in)

	out, ok := c.Get("nil-slices-skill")
	assert.True(ok)
	assert.Nil(out.AllowedTools)
	assert.Nil(out.Metadata)
}
