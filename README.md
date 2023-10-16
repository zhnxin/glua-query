# glua-query

glua-query provide an easy way to use [goquery](https://github.com/PuerkitoBio/goquery) - a little like that j-thing, only in Go.

## Usage

```go
package main

import (
	"github.com/zhnxin/glua-query"
	"github.com/yuin/gopher-lua"
)

func main(){
    L := lua.NewState()
	defer L.Close()
    query.Preload(L)
	args := L.CreateTable(1, 1)
	args.Append(lua.LString(string([]byte{60, 117, 108, 62, 60, 108, 105, 32, 104, 114, 101, 102, 61, 34, 104, 116, 116, 112, 58, 47, 47, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 34, 62, 210, 188, 60, 47, 108, 105, 62, 60, 108, 105, 62, 183, 161, 60, 47, 108, 105, 62, 60, 108, 105, 62, 200, 254, 60, 47, 108, 105, 62})))
	L.SetGlobal("args", args)
    err := L.DoString(`
local query = require("goquery")
local node,err = query.new('<ul>    <li><a href="http://github.com/zhnxin/glua-query">1</a></li>    <li>2</li>    <li>3</li></ul>')
if err == nil then
    local element = node:find_first('li')
    if element:is_not_empty() then
        print('href:'..element:find_first('a'):attr('href'))
        print('text:'..element:text())
    else
        print('nothing match')
    end
else
    print(err)
end

print('example for decode document with gb18030')

node,err = query.new(args[1],'li','gb18030')
if err == nil then
    local element = node:find_first()
    if element:is_not_empty() then
        print('href:'..element:attr('href'))
        print('text:'..element:text())
    else
        print('nothing match')
    end
else
    print(err)
end
`)
    if err != nil{
        panic(err)
    }
}
```

## API

- `query.new(document [, selector, encoding])`
- `query:attr(attr_name)`
- `query:text()`
- `query:html()`
- `query:is_empty()`
- `query:is_not_empty()`
- `query:find([selector])`
- `query:find_all([selector])`-- alias name of `query:find([selector])`
- `query:find_first()`
