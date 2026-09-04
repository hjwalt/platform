package agent_skill_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hjwalt/platform/agent/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// FindSkillMd
// ---------------------------------------------------------------------------

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestFindSkillMdPrefersUppercase(t *testing.T) {
	dir := t.TempDir()
	upper := writeFile(t, dir, "SKILL.md", "# x")
	writeFile(t, dir, "skill.md", "# lower")

	path, found := agent_skill.FindSkillMd(dir)

	require.True(t, found)
	assert.Equal(t, upper, path)
}

func TestFindSkillMdAcceptsLowercase(t *testing.T) {
	dir := t.TempDir()
	lower := writeFile(t, dir, "skill.md", "# x")

	path, found := agent_skill.FindSkillMd(dir)

	require.True(t, found)
	assert.Equal(t, lower, path)
}

func TestFindSkillMdAcceptsAgentsMdVariants(t *testing.T) {
	t.Run("AGENTS.md", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "AGENTS.md", "# x")

		foundPath, found := agent_skill.FindSkillMd(dir)

		require.True(t, found)
		assert.Equal(t, path, foundPath)
	})

	t.Run("agents.md", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "agents.md", "# x")

		foundPath, found := agent_skill.FindSkillMd(dir)

		require.True(t, found)
		assert.Equal(t, path, foundPath)
	})
}

func TestFindSkillMdReturnsFalseWhenNothingMatches(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "README.md", "hi")

	path, found := agent_skill.FindSkillMd(dir)

	assert.False(t, found)
	assert.Empty(t, path)
}

func TestFindSkillMdReturnsTrueForDirectoriesNamedSkillMd(t *testing.T) {
	// os.Stat succeeds for directories too, so a directory called SKILL.md is
	// "found" even though it is not a readable markdown file.
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "SKILL.md"), 0o700))

	path, found := agent_skill.FindSkillMd(dir)

	require.True(t, found)
	assert.Equal(t, filepath.Join(dir, "SKILL.md"), path)
}

// ---------------------------------------------------------------------------
// ParseFrontmatter
// ---------------------------------------------------------------------------

func TestParseFrontmatterParsesValidFrontmatter(t *testing.T) {
	assert := assert.New(t)

	content := "---\n" +
		"name: test-skill\n" +
		"description: a test skill\n" +
		"count: 3\n" +
		"enabled: true\n" +
		"---\n" +
		"\n" +
		"# Body\n" +
		"\n" +
		"Body line two.\n"

	metadata, body, err := agent_skill.ParseFrontmatter(content)

	assert.NoError(err)
	require.NotNil(t, metadata)
	assert.Equal("test-skill", metadata["name"])
	assert.Equal("a test skill", metadata["description"])
	assert.Equal("# Body\n\nBody line two.", body)
}

func TestParseFrontmatterTrimsBodyWhitespace(t *testing.T) {
	assert := assert.New(t)

	content := "---\nname: x\n---\n\n   \nkeep me\n  \n"

	_, body, err := agent_skill.ParseFrontmatter(content)

	assert.NoError(err)
	assert.Equal("keep me", body)
}

func TestParseFrontmatterRequiresOpeningDelimiter(t *testing.T) {
	assert := assert.New(t)

	_, _, err := agent_skill.ParseFrontmatter("name: x\n---\nbody")

	require.Error(t, err)
	var parseErr agent_skill.SkillError
	assert.True(errors.As(err, &parseErr))
	assert.Equal("SKILL.md must start with YAML frontmatter (---)", err.Error())
}

func TestParseFrontmatterRequiresClosingDelimiter(t *testing.T) {
	assert := assert.New(t)

	_, _, err := agent_skill.ParseFrontmatter("---\nname: x\n")

	require.Error(t, err)
	assert.Equal("SKILL.md frontmatter not properly closed with ---", err.Error())
}

func TestParseFrontmatterRejectsInvalidYaml(t *testing.T) {
	assert := assert.New(t)

	_, _, err := agent_skill.ParseFrontmatter("---\nname: [unclosed\n---\nbody")

	require.Error(t, err)
	var parseErr agent_skill.SkillError
	assert.True(errors.As(err, &parseErr))
	assert.Contains(err.Error(), "Invalid YAML in frontmatter")
}

func TestParseFrontmatterRejectsNullDocument(t *testing.T) {
	assert := assert.New(t)

	_, _, err := agent_skill.ParseFrontmatter("---\nnull\n---\nbody")

	require.Error(t, err)
	assert.Equal("SKILL.md frontmatter must be a YAML mapping", err.Error())
}

func TestParseFrontmatterRejectsEmptyDocument(t *testing.T) {
	assert := assert.New(t)

	_, _, err := agent_skill.ParseFrontmatter("---\n---\nbody")

	require.Error(t, err)
	assert.Equal("SKILL.md frontmatter must be a YAML mapping", err.Error())
}

func TestParseFrontmatterRejectsScalarDocument(t *testing.T) {
	assert := assert.New(t)

	_, _, err := agent_skill.ParseFrontmatter("---\njust a string\n---\nbody")

	require.Error(t, err)
	assert.Contains(err.Error(), "Invalid YAML in frontmatter")
}

func TestParseFrontmatterRejectsSequenceDocument(t *testing.T) {
	assert := assert.New(t)

	_, _, err := agent_skill.ParseFrontmatter("---\n- a\n- b\n---\nbody")

	require.Error(t, err)
	assert.Contains(err.Error(), "Invalid YAML in frontmatter")
}

func TestParseFrontmatterNormalizesMetadataSubmap(t *testing.T) {
	assert := assert.New(t)

	content := "---\n" +
		"name: x\n" +
		"metadata:\n" +
		"  owner: alice\n" +
		"  retention_days: 30\n" +
		"  verified: true\n" +
		"---\nbody"

	metadata, _, err := agent_skill.ParseFrontmatter(content)

	assert.NoError(err)
	normalized, ok := metadata["metadata"].(map[string]string)
	require.True(t, ok, "metadata submap must be normalised to map[string]string")
	assert.Equal(map[string]string{
		"owner":          "alice",
		"retention_days": "30",
		"verified":       "true",
	}, normalized)
}

func TestParseFrontmatterLeavesNonMapMetadataUntouched(t *testing.T) {
	assert := assert.New(t)

	content := "---\nname: x\nmetadata: plain-string\n---\nbody"

	metadata, _, err := agent_skill.ParseFrontmatter(content)

	assert.NoError(err)
	assert.Equal("plain-string", metadata["metadata"])
}

// ---------------------------------------------------------------------------
// ReadProperties
// ---------------------------------------------------------------------------

const validSkillMd = "---\n" +
	"name: my-skill\n" +
	"description: does things\n" +
	"license: MIT\n" +
	"compatibility: linux\n" +
	"allowed-tools:\n" +
	"  - web_search\n" +
	"  - web_fetch\n" +
	"metadata:\n" +
	"  owner: alice\n" +
	"  retention_days: 30\n" +
	"---\n" +
	"\n" +
	"# My Skill\n" +
	"\n" +
	"Body content.\n"

func TestReadPropertiesExtractsValidSkill(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "SKILL.md", validSkillMd)

	props, err := agent_skill.ReadProperties(dir)

	assert.NoError(err)
	assert.Equal("my-skill", props.Name)
	assert.Equal("does things", props.Description)
	assert.Equal("# My Skill\n\nBody content.", props.Body)
	require.True(t, props.License.IsPresent())
	assert.Equal("MIT", props.License.Get())
	require.True(t, props.Compatibility.IsPresent())
	assert.Equal("linux", props.Compatibility.Get())
	assert.Equal(map[string]string{"owner": "alice", "retention_days": "30"}, props.Metadata)

	// current behaviour: yaml.v3 decodes the sequence into []interface{} so the
	// []string type assertion in ReadProperties never matches; AllowedTools
	// stays empty even though the field is present in the file.
	assert.Empty(props.AllowedTools)
}

func TestReadPropertiesAcceptsLowercaseFilename(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "skill.md", validSkillMd)

	props, err := agent_skill.ReadProperties(dir)

	assert.NoError(err)
	assert.Equal("my-skill", props.Name)
}

func TestReadPropertiesMissingSkillMdReturnsParseError(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "README.md", "hi")

	props, err := agent_skill.ReadProperties(dir)

	assert.Error(err)
	assert.Empty(props)
	var parseErr agent_skill.SkillError
	assert.True(errors.As(err, &parseErr))
	assert.Equal("SKILL.md not found in "+dir, err.Error())
}

func TestReadPropertiesDirectoryNamedSkillMdFailsToRead(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "SKILL.md"), 0o700))

	_, err := agent_skill.ReadProperties(dir)

	require.Error(t, err)
	assert.Contains(err.Error(), "Failed to read SKILL.md")
}

func TestReadPropertiesMissingNameReturnsValidationError(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "SKILL.md", "---\ndescription: only desc\n---\nbody")

	_, err := agent_skill.ReadProperties(dir)

	require.Error(t, err)
	var validationErr *agent_skill.ValidationError
	assert.True(errors.As(err, &validationErr))
	assert.Equal("Missing required field in frontmatter: name", err.Error())
}

func TestReadPropertiesMissingDescriptionReturnsValidationError(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "SKILL.md", "---\nname: skill-x\n---\nbody")

	_, err := agent_skill.ReadProperties(dir)

	require.Error(t, err)
	var validationErr *agent_skill.ValidationError
	assert.True(errors.As(err, &validationErr))
	assert.Equal("Missing required field in frontmatter: description", err.Error())
}

func TestReadPropertiesNonStringNameReturnsValidationError(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "SKILL.md", "---\nname: 123\ndescription: desc\n---\nbody")

	_, err := agent_skill.ReadProperties(dir)

	require.Error(t, err)
	var validationErr *agent_skill.ValidationError
	assert.True(errors.As(err, &validationErr))
	assert.Equal("Field 'name' must be a non-empty string", err.Error())
}

func TestReadPropertiesEmptyNameReturnsValidationError(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "SKILL.md", "---\nname: ''\ndescription: desc\n---\nbody")

	_, err := agent_skill.ReadProperties(dir)

	require.Error(t, err)
	assert.Equal("Field 'name' must be a non-empty string", err.Error())
}

func TestReadPropertiesInvalidYamlReturnsParseError(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "SKILL.md", "---\nname: [unclosed\n---\nbody")

	_, err := agent_skill.ReadProperties(dir)

	require.Error(t, err)
	var parseErr agent_skill.SkillError
	assert.True(errors.As(err, &parseErr))
	var validationErr *agent_skill.ValidationError
	assert.False(errors.As(err, &validationErr))
}
