package scripting

import "fmt"

// StdLib provides safe Lua standard library functions for plugins
type StdLib struct {
	allowed map[string]bool
}

// NewStdLib creates a new restricted standard library
func NewStdLib() *StdLib {
	lib := &StdLib{
		allowed: make(map[string]bool),
	}

	// Only allow safe modules
	lib.allowed["table"] = true
	lib.allowed["string"] = true
	lib.allowed["math"] = true
	lib.allowed["utf8"] = true

	return lib
}

// LoadTableLibrary loads safe table functions
func (lib *StdLib) LoadTableLibrary() error {
	if !lib.allowed["table"] {
		return fmt.Errorf("table library is not allowed")
	}

	// Available functions from table library:
	// - insert(table, value)
	// - remove(table, [pos])
	// - concat(table, [sep])
	// - sort(table, [comp])
	// - unpack(table) / table.unpack(table)

	// Real implementation would register these with Lua VM
	return nil
}

// LoadStringLibrary loads safe string functions
func (lib *StdLib) LoadStringLibrary() error {
	if !lib.allowed["string"] {
		return fmt.Errorf("string library is not allowed")
	}

	// Available functions from string library:
	// - sub(s, i [, j])
	// - len(s)
	// - rep(s, n)
	// - upper(s)
	// - lower(s)
	// - find(s, pattern [, init [, plain]])
	// - match(s, pattern [, init])
	// - gsub(s, pattern, repl [, n])
	// - format(formatstring, ...)
	// - reverse(s)
	// - byte(s [, i [, j]])
	// - char(...)

	// Real implementation would register these with Lua VM
	return nil
}

// LoadMathLibrary loads safe math functions
func (lib *StdLib) LoadMathLibrary() error {
	if !lib.allowed["math"] {
		return fmt.Errorf("math library is not allowed")
	}

	// Available functions from math library:
	// - abs(x)
	// - floor(x)
	// - ceil(x)
	// - min(x, ...)
	// - max(x, ...)
	// - fmod(x, y)
	// - modf(x)
	// - pow(x, y)
	// - sqrt(x)
	// - exp(x)
	// - log(x [, base])
	// - sin(x), cos(x), tan(x)
	// - asin(x), acos(x), atan(x), atan2(y, x)
	// - sinh(x), cosh(x), tanh(x)
	// - deg(x), rad(x)
	// - random([m [, n]])
	// - randomseed([x])
	// - pi (constant)

	// Real implementation would register these with Lua VM
	return nil
}

// LoadUtf8Library loads safe UTF-8 functions
func (lib *StdLib) LoadUtf8Library() error {
	if !lib.allowed["utf8"] {
		return fmt.Errorf("utf8 library is not allowed")
	}

	// Available functions from utf8 library:
	// - len(s)
	// - codes(s)
	// - codepoint(s, i, j)
	// - offset(s, n, i)
	// - char(...)

	// Real implementation would register these with Lua VM
	return nil
}

// BlockModule blocks a module from being loaded
func (lib *StdLib) BlockModule(name string) {
	lib.allowed[name] = false
}

// AllowModule allows a module to be loaded
func (lib *StdLib) AllowModule(name string) {
	if name == "table" || name == "string" || name == "math" || name == "utf8" {
		lib.allowed[name] = true
	}
}

// IsModuleAllowed checks if a module is allowed
func (lib *StdLib) IsModuleAllowed(name string) bool {
	return lib.allowed[name]
}

// LoadAll loads all safe standard library modules
func (lib *StdLib) LoadAll() error {
	var errors []error

	if err := lib.LoadTableLibrary(); err != nil {
		errors = append(errors, err)
	}

	if err := lib.LoadStringLibrary(); err != nil {
		errors = append(errors, err)
	}

	if err := lib.LoadMathLibrary(); err != nil {
		errors = append(errors, err)
	}

	if err := lib.LoadUtf8Library(); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors loading standard libraries: %v", errors)
	}

	return nil
}

// DisallowedModules returns a list of disallowed modules
func DisallowedModules() []string {
	return []string{
		"io",        // File I/O
		"os",        // Operating system
		"debug",     // Debugging
		"package",   // Package loading
		"load",      // Dynamic code loading
		"loadfile",  // File loading
		"require",   // Module requiring
		"dofile",    // File execution
		"coroutine", // Coroutines
	}
}

// ForbiddenGlobals returns a list of forbidden global functions
func ForbiddenGlobals() []string {
	return []string{
		"loadstring",     // Execute string as code
		"load",           // Load code
		"dofile",         // Execute file
		"loadfile",       // Load file
		"require",        // Require module
		"module",         // Module declaration
		"setmetatable",   // Metatables
		"getmetatable",   // Metatables
		"rawget",         // Raw table access
		"rawset",         // Raw table access
		"rawlen",         // Raw table length
		"debug",          // Debug library
		"collectgarbage", // Memory management
	}
}

// Array functions (placeholders - actual implementation would use Lua tables)

// Example built-in scripts

const (
	// MaxStatsScript maximizes all character stats
	MaxStatsScript = `
-- Maximize all character stats
function maxAllStats()
    for charID = 0, 15 do
        setCharacterLevel(charID, 99)
        setCharacterHP(charID, 9999)
        setCharacterMP(charID, 999)
        setCharacterStat(charID, "vigor", 128)
        setCharacterStat(charID, "speed", 128)
        setCharacterStat(charID, "stamina", 128)
        setCharacterStat(charID, "magic", 128)
        log("Maxed stats for character " .. charID)
    end
    print("All characters maximized!")
end

maxAllStats()
`

	// GiveAllItemsScript gives 99 of all items
	GiveAllItemsScript = `
-- Give 99 of all items
function giveAllItems()
    for itemID = 0, 255 do
        setItemCount(itemID, 99)
    end
    print("Gave 99 of all items!")
end

giveAllItems()
`

	// LearnAllMagicScript learns all magic for all characters
	LearnAllMagicScript = `
-- Learn all magic for all characters
function learnAllMagic()
    for charID = 0, 15 do
        for spellID = 0, 53 do
            learnMagic(charID, spellID)
        end
        log("Learned all magic for character " .. charID)
    end
    print("All characters know all magic!")
end

learnAllMagic()
`

	// BalancedPartyScript creates a balanced 4-character party
	BalancedPartyScript = `
-- Create balanced party
function createBalancedParty()
    -- Terra, Edgar, Celes, Sabin
    setPartyMembers({0, 1, 6, 5})
    
    -- Set levels to 50
    for _, charID in ipairs({0, 1, 6, 5}) do
        setCharacterLevel(charID, 50)
    end
    
    print("Created balanced party!")
end

createBalancedParty()
`
)

// BuiltInScripts returns a map of built-in scripts
func BuiltInScripts() map[string]string {
	return map[string]string{
		"Max All Stats":   MaxStatsScript,
		"Give All Items":  GiveAllItemsScript,
		"Learn All Magic": LearnAllMagicScript,
		"Balanced Party":  BalancedPartyScript,
	}
}
