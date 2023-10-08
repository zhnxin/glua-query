package query

import (
	"github.com/PuerkitoBio/goquery"
	lua "github.com/yuin/gopher-lua"
)

var (
	luaQeuryTypeName = "query"
)

type IQuerySelection interface {
	Length() int
	Slice(start, end int) *goquery.Selection
	Filter(selector string) *goquery.Selection
	Find(selector string) *goquery.Selection
	First() *goquery.Selection
	Attr(name string) (string, bool)
	Text() string
	Html() (string, error)
}

type querySelection struct {
	element IQuerySelection
}

func newLUserData(q *querySelection, L *lua.LState) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = q
	L.SetMetatable(ud, L.GetTypeMetatable(luaQeuryTypeName))
	return ud
}

func newLUserDataWithSelection(element *goquery.Selection, L *lua.LState) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &querySelection{element: element}
	L.SetMetatable(ud, L.GetTypeMetatable(luaQeuryTypeName))
	return ud
}

func (q *querySelection) text(L *lua.LState) int {
	L.Push(lua.LString(q.element.Text()))
	return 1
}

func (q *querySelection) attr(L *lua.LState) int {
	attrName := L.CheckString(2)
	if attrName == "" {
		L.Push(lua.LString(""))
		return 1
	}
	if attr, ok := q.element.Attr(attrName); ok {
		L.Push(lua.LString(attr))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func (q *querySelection) html(L *lua.LState) int {
	html, err := q.element.Html()
	L.Push(lua.LString(html))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}
func (q *querySelection) isEmpty(L *lua.LState) int {
	L.Push(lua.LBool(q.element.Length() == 0))
	return 1
}
func (q *querySelection) isNotEmpty(L *lua.LState) int {
	L.Push(lua.LBool(q.element.Length() > 0))
	return 1
}
func (q *querySelection) findFirst(L *lua.LState) int {
	var selector string
	if L.GetTop() > 1 {
		selector = L.CheckString(2)
	}
	var element *goquery.Selection
	if selector == "" {
		element = q.element.First()
	} else {
		element = q.element.Find(selector).First()
	}
	if element.Length() > 0 {
		L.Push(newLUserDataWithSelection(element, L))
		return 1
	}
	return 1

}
func (q *querySelection) find(L *lua.LState) int {
	var selector string
	if L.GetTop() > 1 {
		selector = L.CheckString(2)
	}
	if selector != "" {
		qs := q.element.Find(selector)
		if qs.Length() == 0 {
			L.Push(lua.LNil)
			return 1
		}
		table := L.CreateTable(qs.Length(), 0)
		for i := 0; i < qs.Length(); i++ {
			table.Append(newLUserDataWithSelection(qs.Slice(i, i+1), L))
		}
		L.Push(table)
	} else {
		table := L.CreateTable(q.element.Length(), 0)
		for i := 0; i < q.element.Length(); i++ {
			table.Append(newLUserDataWithSelection(q.element.Slice(i, i+1), L))
		}
		L.Push(table)
	}
	return 1
}
