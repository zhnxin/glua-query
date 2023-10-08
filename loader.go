package query

import lua "github.com/yuin/gopher-lua"

// Preload adds query to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local query = require("query")
func Preload(L *lua.LState) {
	L.PreloadModule("query", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": ApiNewWithEncoding,
	})
	registerQeuryType(mod, L)
	L.Push(mod)
	return 1
}
