package xdb

import (
	"testing"

	"github.com/gogf/gf/frame/g"
)

func TestDefault(t *testing.T) {
	g.Dump(g.Cfg().Get("database"))
	t.Log(g.Cfg().Get("database.user.0.link"))
}

func TestGetTable(t *testing.T) {
	GetTable()
}
