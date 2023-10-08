package query

import (
	"log"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestUsageCode(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	Preload(L)
	args := L.CreateTable(1, 1)
	args.Append(lua.LString(string([]byte{60, 117, 108, 62, 60, 108, 105, 32, 104, 114, 101, 102, 61, 34, 104, 116, 116, 112, 58, 47, 47, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 34, 62, 210, 188, 60, 47, 108, 105, 62, 60, 108, 105, 62, 183, 161, 60, 47, 108, 105, 62, 60, 108, 105, 62, 200, 254, 60, 47, 108, 105, 62})))
	L.SetGlobal("args", args)
	err := L.DoString(`
local query = require("query")
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
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestFindFirst(t *testing.T) {
	state := lua.NewState()
	defer state.Close()
	state.PreloadModule("query", Loader)
	source := `
local query = require("query")
local select,err = query.new('<ul><li>1</li><li>2</li><li>3</li>','','ass')
if err == nil then
	local e = select:find_first('li')
	if e then
		print(222)
		print(e:text())
	else
		print(333)
	end
else
	print(err)
end
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}
func TestFindFirstM2(t *testing.T) {
	state := lua.NewState()
	defer state.Close()
	state.PreloadModule("query", Loader)
	source := `
local query = require("query")
local select = query.new('<ul><li>1</li><li>2</li><li>3</li>','li')
if err then
print(111)
	print(err)
else
	local e = select:find_first()
	if e then
		print(222)
		print(e:text())
	else
		print(333)
	end
end
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}
func TestFindFirstEmpty(t *testing.T) {
	state := lua.NewState()
	defer state.Close()
	state.PreloadModule("query", Loader)
	source := `
local query = require("query")
local select = query.new('<ul><li>1</li><li>2</li><li>3</li>','h1')
if err then
print(111)
	print(err)
else
	local e = select:find_first()
	if e ~= nil then
		print(e:text())
	else
		print('')
	end
end
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}

func TestFindAll(t *testing.T) {
	state := lua.NewState()
	defer state.Close()
	state.PreloadModule("query", Loader)
	source := `
	local query = require("query")
	local select = query.new('<ul><li href="http://github.com">1</li><li>2</li><li>3</li>','li')
	if err then
		print(111)
		print(err)
	else
		for _,e in ipairs(select:find()) do
			print(e:attr('href'))
			print(e:text())
		end
	end
	`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}

func TestEncoding(t *testing.T) {
	body := []byte{60, 117, 108, 62, 60, 108, 105, 32, 104, 114, 101, 102, 61, 34, 104, 116, 116, 112, 58, 47, 47, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 34, 62, 210, 188, 60, 47, 108, 105, 62, 60, 108, 105, 62, 183, 161, 60, 47, 108, 105, 62, 60, 108, 105, 62, 200, 254, 60, 47, 108, 105, 62}
	L := lua.NewState()
	defer L.Close()
	args := L.CreateTable(1, 1)
	L.PreloadModule("query", Loader)
	args.Append(lua.LString(string(body)))
	L.SetGlobal("args", args)
	source := `
	local query = require("query")
	local select = query.new(args[1],'li','gb18030')
	if err then
		print(111)
		print(err)
	else
		for _,e in ipairs(select:find()) do
			print(e:attr('href'))
			print(e:text())
		end
	end
	`
	if err := L.DoString(source); err != nil {
		log.Fatal(err.Error())
	}

}
