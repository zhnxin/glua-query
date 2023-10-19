package query

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func registerQeuryType(module *lua.LTable, L *lua.LState) {
	mt := L.NewTypeMetatable(luaQeuryTypeName)
	L.SetGlobal(luaQeuryTypeName, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), queryFunction))
}

func ApiNewWithEncoding(L *lua.LState) int {
	bodyString := L.CheckString(1)
	body := []byte(bodyString)
	var selection string
	var encoding string
	var doc *goquery.Document
	var err error
	if L.GetTop() > 1 {
		selection = L.CheckString(2)
	}
	if L.GetTop() > 2 {
		encoding = L.CheckString(3)
	}
	L.Pop(L.GetTop())
	switch strings.ToLower(encoding) {
	case "gb18030":
		doc, err = goquery.NewDocumentFromReader(transform.NewReader(bytes.NewReader(body),
			simplifiedchinese.GB18030.NewDecoder()))
	case "gbk":
		doc, err = goquery.NewDocumentFromReader(transform.NewReader(bytes.NewReader(body), simplifiedchinese.GBK.NewDecoder()))
	case "", "utf-8", "utf8":
		doc, err = goquery.NewDocumentFromReader(bytes.NewReader(body))
	default:
		err = fmt.Errorf("encoding only support gb18030, gbk, utf-8")
	}
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	if selection == "" {
		L.Push(newLUserData(&querySelection{element: doc}, L))
	} else {
		L.Push(newLUserData(&querySelection{element: doc.Find(selection)}, L))
	}
	L.Push(lua.LNil)
	return 2
}
func ApiNew(L *lua.LState) int {
	body := L.CheckString(1)
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(body))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(newLUserData(&querySelection{element: doc}, L))
	return 1
}
func checkQuerySelection(L *lua.LState) *querySelection {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*querySelection); ok {
		return v
	}
	L.ArgError(1, "query.selection expected")
	return nil
}

var queryFunction = map[string]lua.LGFunction{
	"text": func(L *lua.LState) int {
		return checkQuerySelection(L).text(L)
	},
	"attr": func(L *lua.LState) int {
		return checkQuerySelection(L).attr(L)
	},
	"is_empty": func(L *lua.LState) int {
		return checkQuerySelection(L).isEmpty(L)
	},
	"is_not_empty": func(L *lua.LState) int {
		return checkQuerySelection(L).isNotEmpty(L)
	},
	"find": func(L *lua.LState) int {
		return checkQuerySelection(L).find(L)
	},
	"find_all": func(L *lua.LState) int {
		return checkQuerySelection(L).find(L)
	},
	"find_first": func(L *lua.LState) int {
		return checkQuerySelection(L).findFirst(L)
	},
	"html": func(L *lua.LState) int {
		return checkQuerySelection(L).html(L)
	},
}
