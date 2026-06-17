package mem

import "strings"

type Kind struct {
	Name   string
	Desc   string
	Prefix string
	Dir    string
}

var kindList = []Kind{
	{Name: "bug", Desc: "バグ、不具合、失敗", Prefix: "BUG", Dir: "bugs"},
	{Name: "prob", Desc: "根本問題、構造的課題", Prefix: "PROB", Dir: "problems"},
	{Name: "adr", Desc: "設計判断", Prefix: "ADR", Dir: "decisions"},
	{Name: "chg", Desc: "実際に行った変更", Prefix: "CHG", Dir: "changes"},
	{Name: "rev", Desc: "巻き戻し、撤回", Prefix: "REV", Dir: "reversions"},
	{Name: "sol", Desc: "再利用可能な解決策", Prefix: "SOL", Dir: "solutions"},
	{Name: "sup", Desc: "過去判断や実装の無効化", Prefix: "SUP", Dir: "supersessions"},
}

var kinds = buildKindMap()
var kindByPrefix = buildPrefixMap()

func buildKindMap() map[string]Kind {
	m := make(map[string]Kind, len(kindList))
	for _, kind := range kindList {
		m[kind.Name] = kind
	}
	return m
}

func buildPrefixMap() map[string]Kind {
	m := make(map[string]Kind, len(kindList))
	for _, kind := range kindList {
		m[kind.Prefix] = kind
	}
	return m
}

func ParseKind(name string) (Kind, bool) {
	kind, ok := kinds[name]
	return kind, ok
}

func KindByPrefix(prefix string) (Kind, bool) {
	kind, ok := kindByPrefix[strings.ToUpper(prefix)]
	return kind, ok
}

func Kinds() []Kind {
	out := make([]Kind, len(kindList))
	copy(out, kindList)
	return out
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

func Statuses() []string {
	return []string{
		"OPEN",
		"INVESTIGATING",
		"UNRESOLVED",
		"PROPOSED",
		"ACCEPTED",
		"FIXED",
		"VERIFIED",
		"SHIPPED",
		"CLOSED",
		"SUPERSEDED",
		"ARCHIVED",
		"ACTIVE",
	}
}

func FutureValues() []string {
	return []string{"ignore", "watch", "reusable"}
}
