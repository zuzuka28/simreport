package bot

type docDetails struct {
	ParentID string   `yaml:"parent_id"`
	Name     string   `yaml:"name"`
	Version  int      `yaml:"version"`
	GroupID  []string `yaml:"group_id"`
}

const docDetailsExample = `
parent_id: 123
name: "Document"
version: 1
group_id: ["123", "456"]
`
