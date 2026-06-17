package mem

type Kind struct {
	Name   string
	Prefix string
	Dir    string
}

var kinds = map[string]Kind{
	"bug":  {Name: "bug", Prefix: "BUG", Dir: "bugs"},
	"prob": {Name: "prob", Prefix: "PROB", Dir: "problems"},
	"adr":  {Name: "adr", Prefix: "ADR", Dir: "decisions"},
	"chg":  {Name: "chg", Prefix: "CHG", Dir: "changes"},
	"rev":  {Name: "rev", Prefix: "REV", Dir: "reversions"},
	"sol":  {Name: "sol", Prefix: "SOL", Dir: "solutions"},
	"sup":  {Name: "sup", Prefix: "SUP", Dir: "supersessions"},
}

func ParseKind(name string) (Kind, bool) {
	kind, ok := kinds[name]
	return kind, ok
}

func Dirs() []string {
	return []string{
		"memory",
		"memory/bugs",
		"memory/problems",
		"memory/decisions",
		"memory/changes",
		"memory/reversions",
		"memory/solutions",
		"memory/supersessions",
		"memory/generated",
	}
}
