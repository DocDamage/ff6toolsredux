package scripting

import (
	"context"
	"ffvi_editor/io/pr"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// LuaResult is a generic map result from Lua tables.
type LuaResult map[string]interface{}

// RunSnippet executes a Lua snippet with sandboxed VM and returns a LuaResult if a table is returned.
func RunSnippet(ctx context.Context, code string) (LuaResult, error) {
	if code == "" {
		return nil, fmt.Errorf("code is empty")
	}
	// Context timeout safeguard
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()

	openSafeLibs(L)

	// Restrict package.path to local plugins directory
	pluginPaths := []string{
		"./?.lua",
		"./?/init.lua",
		"./plugins/?.lua",
		"./plugins/?/init.lua",
		"./plugins/?/v1_0_core.lua",
	}
	_ = L.DoString(fmt.Sprintf(`package.path = package.path .. ';%s'`, strings.Join(pluginPaths, ";")))
	disableUnsafeGlobals(L)

	done := make(chan error, 1)
	go func() {
		done <- L.DoString(code)
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, err
		}
		// If a value was returned, convert when table
		if L.GetTop() >= 1 {
			val := L.Get(-1)
			if tbl, ok := val.(*lua.LTable); ok {
				return tableToMap(tbl), nil
			}
		}
		return nil, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("lua execution timeout: %w", ctx.Err())
	}
}

// RunSnippetWithSave executes a Lua snippet with save data bindings.
func RunSnippetWithSave(ctx context.Context, code string, save *pr.PR) (LuaResult, error) {
	if code == "" {
		return nil, fmt.Errorf("code is empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()

	openSafeLibs(L)

	// Restrict package.path to local plugins directory
	pluginPaths := []string{
		"./?.lua",
		"./?/init.lua",
		"./plugins/?.lua",
		"./plugins/?/init.lua",
		"./plugins/?/v1_0_core.lua",
	}
	_ = L.DoString(fmt.Sprintf(`package.path = package.path .. ';%s'`, strings.Join(pluginPaths, ";")))
	disableUnsafeGlobals(L)

	// Register save bindings if save provided
	if save != nil {
		registerSaveBindings(L, save)
	}

	done := make(chan error, 1)
	go func() {
		done <- L.DoString(code)
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, err
		}
		if L.GetTop() >= 1 {
			val := L.Get(-1)
			if tbl, ok := val.(*lua.LTable); ok {
				return tableToMap(tbl), nil
			}
		}
		return nil, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("lua execution timeout: %w", ctx.Err())
	}
}

func tableToMap(tbl *lua.LTable) LuaResult {
	result := LuaResult{}
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		k := key.String()
		switch v := value.(type) {
		case *lua.LTable:
			result[k] = tableToMap(v)
		case lua.LBool:
			result[k] = bool(v)
		case lua.LNumber:
			result[k] = float64(v)
		default:
			result[k] = v.String()
		}
	})
	return result
}

// BuildEncounterScript builds Lua code for encounter tuning.
func BuildEncounterScript(zone string, rate, elite string) string {
	return fmt.Sprintf(`local pack = require('plugins.combat-depth-pack.v1_0_core')
local res = pack.EncounterTuner.configureEncounterRates({zone = '%s', encounter_rate = %s, elite_chance = %s})
return res`, escape(zone), rate, elite)
}

// BuildBossScript builds Lua code for boss remix.
func BuildBossScript(affixes string) string {
	return fmt.Sprintf(`local pack = require('plugins.combat-depth-pack.v1_0_core')
local affix_list = {}
for affix in string.gmatch('%s', '([^,]+)') do table.insert(affix_list, affix) end
local res = pack.BossRemix.generateRemix({name = 'Boss', hp = 50000, attack = 1200}, {affixes = affix_list})
return res`, escape(affixes))
}

// BuildCompanionScript builds Lua code for companion profile definition and evaluation.
func BuildCompanionScript(profile, risk string) string {
	return fmt.Sprintf(`local pack = require('plugins.combat-depth-pack.v1_0_core')
pack.CompanionDirector.defineProfile('%s', {risk_tolerance = '%s', priorities = {'damage','survival'}})
local res = pack.CompanionDirector.recommendAction({hp_status='stable', threat_level='medium', primary_target='boss'}, '%s')
return res`, escape(profile), escape(risk), escape(profile))
}

// BuildSmokeScript builds Lua code to run smoke tests.
func BuildSmokeScript() string {
	return "local t = require('plugins.combat_depth_pack_smoke_tests'); return t.run_all()"
}

func escape(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

// openSafeLibs loads a restricted set of Lua standard libraries.
func openSafeLibs(L *lua.LState) {
	lua.OpenBase(L)
	lua.OpenTable(L)
	lua.OpenString(L)
	lua.OpenMath(L)
	lua.OpenPackage(L)
	// Note: OpenUtf8 not available in gopher-lua v1.1.0
}

// disableUnsafeGlobals removes dangerous functions and modules from the global environment.
func disableUnsafeGlobals(L *lua.LState) {
	unsafeGlobals := []string{"dofile", "loadfile", "load", "loadstring", "collectgarbage", "module"}
	for _, name := range unsafeGlobals {
		L.SetGlobal(name, lua.LNil)
	}

	// Remove IO/OS/debug tables if present (not opened by openSafeLibs, but defensive).
	L.SetGlobal("io", lua.LNil)
	L.SetGlobal("os", lua.LNil)
	L.SetGlobal("debug", lua.LNil)

	// Trim package fields that could escape the sandbox.
	if pkg, ok := L.GetGlobal("package").(*lua.LTable); ok {
		pkg.RawSetString("loadlib", lua.LNil)
		pkg.RawSetString("cpath", lua.LString(""))
	}
}

// PluginPath returns the plugins directory (currently relative).
func PluginPath() string {
	return filepath.Join(".", "plugins")
}

// registerSaveBindings registers Go functions for save data manipulation in Lua.
func registerSaveBindings(L *lua.LState, save *pr.PR) {
	// Create save table
	saveTable := L.NewTable()

	// Character access
	L.SetField(saveTable, "getCharacterCount", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(len(save.Characters)))
		return 1
	}))

	L.SetField(saveTable, "getCharacterName", L.NewFunction(func(L *lua.LState) int {
		idx := int(L.CheckNumber(1))
		if idx < 0 || idx >= len(save.Characters) {
			L.Push(lua.LNil)
			L.Push(lua.LString("invalid character index"))
			return 2
		}
		if save.Characters[idx] == nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("character not initialized"))
			return 2
		}
		nameVal := save.Characters[idx].Get("name")
		if nameVal != nil {
			if nameStr, ok := nameVal.(string); ok {
				L.Push(lua.LString(nameStr))
				return 1
			}
		}
		L.Push(lua.LString("Unknown"))
		return 1
	}))

	L.SetField(saveTable, "setCharacterLevel", L.NewFunction(func(L *lua.LState) int {
		idx := int(L.CheckNumber(1))
		level := int(L.CheckNumber(2))
		if idx < 0 || idx >= len(save.Characters) || save.Characters[idx] == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("invalid character index"))
			return 2
		}
		save.Characters[idx].Set("level", float64(level))
		L.Push(lua.LBool(true))
		return 1
	}))

	L.SetField(saveTable, "setCharacterHP", L.NewFunction(func(L *lua.LState) int {
		idx := int(L.CheckNumber(1))
		hp := int(L.CheckNumber(2))
		if idx < 0 || idx >= len(save.Characters) || save.Characters[idx] == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("invalid character index"))
			return 2
		}
		save.Characters[idx].Set("current_hp", float64(hp))
		L.Push(lua.LBool(true))
		return 1
	}))

	L.SetField(saveTable, "setCharacterMP", L.NewFunction(func(L *lua.LState) int {
		idx := int(L.CheckNumber(1))
		mp := int(L.CheckNumber(2))
		if idx < 0 || idx >= len(save.Characters) || save.Characters[idx] == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("invalid character index"))
			return 2
		}
		save.Characters[idx].Set("current_mp", float64(mp))
		L.Push(lua.LBool(true))
		return 1
	}))

	// Inventory/party placeholder stubs
	L.SetField(saveTable, "getGil", L.NewFunction(func(L *lua.LState) int {
		val := save.UserData.Get("gil")
		if val != nil {
			if gil, ok := val.(float64); ok {
				L.Push(lua.LNumber(gil))
				return 1
			}
		}
		L.Push(lua.LNumber(0))
		return 1
	}))

	L.SetField(saveTable, "setGil", L.NewFunction(func(L *lua.LState) int {
		gil := int(L.CheckNumber(1))
		save.UserData.Set("gil", float64(gil))
		L.Push(lua.LBool(true))
		return 1
	}))

	L.SetField(saveTable, "log", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		fmt.Printf("[LUA] %s\n", msg)
		return 0
	}))

	// Register global save table
	L.SetGlobal("save", saveTable)
}
