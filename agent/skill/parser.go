package agent_skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjwalt/platform/type/optional"
	"github.com/stretchr/testify/assert/yaml"
)

// FindSkillMd finds the SKILL.md file in a skill directory.
// It prefers SKILL.md (uppercase) but accepts skill.md (lowercase).
// Returns the path to the file and true if found, empty string and false otherwise.
func FindSkillMd(skillDir string) (string, bool) {
	for _, name := range []string{"SKILL.md", "skill.md", "AGENTS.md", "agents.md"} {
		path := filepath.Join(skillDir, name)
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}
	return "", false
}

// ParseFrontmatter parses YAML frontmatter from SKILL.md content.
//
// The content must start with "---" followed by YAML, then closed with "---".
// Returns the parsed metadata as a map, the markdown body (content after
// the closing "---"), and any parse error.
//
// The metadata["metadata"] field, if present, is normalized to map[string]string.
func ParseFrontmatter(content string) (map[string]interface{}, string, error) {
	if !strings.HasPrefix(content, "---") {
		return nil, "", NewParseError("SKILL.md must start with YAML frontmatter (---)")
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, "", NewParseError("SKILL.md frontmatter not properly closed with ---")
	}

	frontmatterStr := parts[1]
	body := strings.TrimSpace(parts[2])

	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatterStr), &metadata); err != nil {
		return nil, "", NewParseErrorf("Invalid YAML in frontmatter: %v", err)
	}

	if metadata == nil {
		return nil, "", NewParseError("SKILL.md frontmatter must be a YAML mapping")
	}

	// Normalize metadata["metadata"] to map[string]string
	if meta, ok := metadata["metadata"]; ok {
		if metaMap, ok := meta.(map[string]interface{}); ok {
			normalized := make(map[string]string)
			for k, v := range metaMap {
				normalized[k] = fmt.Sprintf("%v", v)
			}
			metadata["metadata"] = normalized
		}
	}

	return metadata, body, nil
}

// ReadProperties reads skill properties from SKILL.md frontmatter.
//
// It locates the SKILL.md file, parses the YAML frontmatter, and returns
// the extracted properties. This function performs basic validation
// (required fields must exist and be non-empty strings) but does NOT
// perform full validation. Use [Validate] for complete validation.
//
// Returns a ParseError if the file is missing or has invalid YAML,
// or a ValidationError if required fields are missing or invalid.
func ReadProperties(skillDir string) (*Skill, error) {
	skillMd, found := FindSkillMd(skillDir)
	if !found {
		return nil, NewParseErrorf("SKILL.md not found in %s", skillDir)
	}

	content, err := os.ReadFile(skillMd)
	if err != nil {
		return nil, NewParseErrorf("Failed to read SKILL.md: %v", err)
	}

	metadata, body, err := ParseFrontmatter(string(content))
	if err != nil {
		return nil, err
	}

	// Check required fields
	nameVal, hasName := metadata["name"]
	if !hasName {
		return nil, NewValidationError("Missing required field in frontmatter: name")
	}

	descVal, hasDesc := metadata["description"]
	if !hasDesc {
		return nil, NewValidationError("Missing required field in frontmatter: description")
	}

	name, ok := nameVal.(string)
	if !ok || strings.TrimSpace(name) == "" {
		return nil, NewValidationError("Field 'name' must be a non-empty string")
	}

	description, ok := descVal.(string)
	if !ok || strings.TrimSpace(description) == "" {
		return nil, NewValidationError("Field 'description' must be a non-empty string")
	}

	props := &Skill{
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Body:        body,
	}

	// Optional fields
	if license, ok := metadata["license"].(string); ok {
		props.License = optional.Of(license)
	}

	if compatibility, ok := metadata["compatibility"].(string); ok {
		props.Compatibility = optional.Of(compatibility)
	}

	if allowedTools, ok := metadata["allowed-tools"].([]string); ok {
		props.AllowedTools = allowedTools
	}

	if meta, ok := metadata["metadata"].(map[string]string); ok {
		props.Metadata = meta
	}

	return props, nil
}
