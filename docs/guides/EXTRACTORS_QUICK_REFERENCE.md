# Type-Safe Extractors - Quick Reference

## Overview
The `type_safe_extractors.go` module provides type-safe alternatives to unsafe type assertions. Use these instead of `value.(type)` patterns.

---

## Quick Examples

### Extract String from interface{}
```go
// ✅ NEW - Safe
str, err := ExtractString(value)
if err != nil {
    return fmt.Errorf("invalid string: %w", err)
}
// str is guaranteed to be a valid string

// ❌ OLD - Unsafe (panics if type wrong)
str := value.(string)
```

### Extract Integer Array from interface{}
```go
// ✅ NEW - Safe
numbers, err := ExtractInt64Array(value)
if err != nil {
    return fmt.Errorf("invalid number array: %w", err)
}
for i, num := range numbers {
    // Process validated numbers
}

// ❌ OLD - Unsafe (panics if type wrong)
arr := value.([]interface{})
for _, item := range arr {
    num := item.(int64)  // Panics!
}
```

### Extract Nested Target Data
```go
// ✅ NEW - Safe (most common use case)
characterStrings, err := SafeGetFromTarget(userDataMap, OwnedCharacterList)
if err != nil {
    return fmt.Errorf("failed to extract characters: %w", err)
}

for j, charJSON := range characterStrings {
    char := jo.NewOrderedMap()
    if err := char.UnmarshalJSON([]byte(charJSON)); err != nil {
        return fmt.Errorf("failed to unmarshal character %d: %w", j, err)
    }
    // Process character
}

// ❌ OLD - Unsafe (multiple panics possible)
chars := userDataMap[OwnedCharacterList]
charList := chars.(string)
parsed := jo.NewOrderedMap()
parsed.UnmarshalJSON([]byte(charList))
// No error handling, type not checked!
```

### Extract Single Integer
```go
// ✅ NEW - Safe
maxHP, err := ExtractInt64(value)
if err != nil {
    return fmt.Errorf("invalid HP value: %w", err)
}

// ❌ OLD - Unsafe
jn := value.(json.Number)
maxHP, _ := jn.Int64()  // Ignores error!
```

---

## Function Reference

### Array Extractors

#### `ExtractStringArray(value interface{}) ([]string, error)`
Safely convert `[]interface{}` to `[]string`.
- **Returns**: Array of strings, error if any element is not string or value is not array
- **Use case**: Character names, inventory item codes, skill names
```go
names, err := ExtractStringArray(userDataMap["characterNames"])
```

#### `ExtractInt64Array(value interface{}) ([]int64, error)`
Safely convert `[]interface{}` to `[]int64`.
- **Returns**: Array of int64 values, error if any element is not json.Number
- **Use case**: Experience totals, money amounts, coordinates
```go
coords, err := ExtractInt64Array(mapPositions)
```

#### `ExtractFloat64Array(value interface{}) ([]float64, error)`
Safely convert `[]interface{}` to `[]float64`.
- **Returns**: Array of float64 values
- **Use case**: Less common; use Int64Array for game data
```go
stats, err := ExtractFloat64Array(floatArray)
```

---

### Scalar Extractors

#### `ExtractString(value interface{}) (string, error)`
Safely extract string from interface{}.
- **Returns**: String value or error if not string
- **Use case**: Single string extraction with error handling
```go
savePath, err := ExtractString(configMap["savePath"])
```

#### `ExtractInt64(value interface{}) (int64, error)`
Safely extract int64 from json.Number.
- **Returns**: int64 value or error if not valid number
- **Use case**: Single number extraction from JSON
```go
maxEspers, err := ExtractInt64(settingsMap["maxEspers"])
```

---

### High-Level Helpers

#### `SafeGetFromTarget(om *jo.OrderedMap, key string) ([]string, error)`
Extract "target" field from nested JSON structure (most common pattern).
- **Returns**: Array of JSON strings from target field
- **Use case**: Loading character/equipment/inventory data
- **Pattern**: 
  ```
  userData: { 
    key: "{\"target\": [\"json1\", \"json2\"]}"
  }
  ```
```go
characters, err := SafeGetFromTarget(userDataMap, OwnedCharacterList)
// Returns array of character JSON strings ready to unmarshal
```

#### `SafeGetFromTargetRaw(om *jo.OrderedMap, key string) ([]interface{}, error)`
Same as SafeGetFromTarget but returns raw interface{} array.
- **Returns**: Raw interface array from target field
- **Use case**: When you need to process mixed types
```go
items, err := SafeGetFromTargetRaw(inventoryMap, OwnedInventoryList)
```

#### `ValidateArray(value interface{}, fieldName string) (int, error)`
Validate that value is an array and return its length.
- **Returns**: Array length or error if not array
- **Use case**: Bounds checking before processing
```go
count, err := ValidateArray(data["party"], "party members")
if err != nil {
    return err  // Better error message with field name
}
```

#### `UnmarshalOrderedMapFromString(json string) (*jo.OrderedMap, error)`
Safely unmarshal JSON string to OrderedMap.
- **Returns**: OrderedMap or error if invalid JSON
- **Use case**: Parsing JSON strings extracted from parent object
```go
char := jo.NewOrderedMap()
charData, err := UnmarshalOrderedMapFromString(charJSONString)
// More reliable than direct json.Unmarshal
```

#### `SafeUnmarshalJSON(data []byte) (*jo.OrderedMap, error)`
Safely unmarshal JSON bytes to OrderedMap.
- **Returns**: OrderedMap or error
- **Use case**: Direct JSON bytes parsing
```go
om, err := SafeUnmarshalJSON([]byte(jsonString))
```

---

## Common Patterns

### Pattern 1: Extract and Process Array
```go
items, err := ExtractStringArray(data["items"])
if err != nil {
    return fmt.Errorf("failed to extract items: %w", err)
}

for i, item := range items {
    om := jo.NewOrderedMap()
    if err := om.UnmarshalJSON([]byte(item)); err != nil {
        return fmt.Errorf("failed to parse item %d: %w", i, err)
    }
    // Process item
}
```

### Pattern 2: Extract and Validate Numbers
```go
numbers, err := ExtractInt64Array(numberArray)
if err != nil {
    return fmt.Errorf("invalid number array: %w", err)
}

for i, num := range numbers {
    if num < 0 || num > 255 {
        return fmt.Errorf("number at index %d out of range: %d", i, num)
    }
    // Use validated number
}
```

### Pattern 3: Nested Target Extraction
```go
// Most common: extract "target" from map
chars, err := SafeGetFromTarget(userData, "characters")
if err != nil {
    return fmt.Errorf("failed to load characters: %w", err)
}

for j, charJSON := range chars {
    character := jo.NewOrderedMap()
    if err := character.UnmarshalJSON([]byte(charJSON)); err != nil {
        return fmt.Errorf("failed to unmarshal character %d: %w", j, err)
    }
    
    // Now safely extract fields from character
    name, err := ExtractString(character.Get("name"))
    if err != nil {
        return fmt.Errorf("character %d missing name: %w", j, err)
    }
    // Use name
}
```

---

## Error Handling

All extractors return descriptive errors:

```go
value, err := ExtractString(data["name"])
if err != nil {
    // Error contains:
    // - What failed (e.g., "not a string")
    // - What was expected (e.g., "expected string")
    // - Actual type received
    return fmt.Errorf("invalid name: %w", err)
}
```

Add context with indices for array processing:
```go
for i, item := range items {
    val, err := ExtractInt64(item)
    if err != nil {
        // Include index in error message
        return fmt.Errorf("invalid value at index %d: %w", i, err)
    }
}
```

---

## Migration Checklist

When refactoring code to use type-safe extractors:

- [ ] Replace `value.(string)` with `ExtractString(value)`
- [ ] Replace `value.([]interface{})` with `ExtractStringArray(value)` or `ExtractInt64Array(value)`
- [ ] Replace `data.(json.Number)` with `ExtractInt64(data)`
- [ ] Add error checking for all extractions
- [ ] Include index context in error messages for loops
- [ ] Use `SafeGetFromTarget()` for common "target" pattern
- [ ] Validate array bounds with `ValidateArray()` when indexing

---

## Performance Notes

- **Negligible overhead**: Each extractor is a single type check + conversion
- **Benchmark included**: `BenchmarkExtract*` functions in test file
- **Memory**: No additional allocations beyond the destination array
- **Expected cost**: <1% performance difference vs unsafe assertions

---

## When to Use Each Function

| Situation | Function | Example |
|-----------|----------|---------|
| Single string value | `ExtractString()` | Character name |
| Single number from JSON | `ExtractInt64()` | Max HP value |
| Array of strings | `ExtractStringArray()` | Character names |
| Array of numbers | `ExtractInt64Array()` | Experience values |
| Array of mixed types | `SafeGetFromTargetRaw()` | Varied inventory items |
| Nested JSON string | `SafeGetFromTarget()` | Character data object |
| Validate array exists | `ValidateArray()` | Check bounds before indexing |
| Parse JSON string | `UnmarshalOrderedMapFromString()` | Extracted JSON data |
| Parse JSON bytes | `SafeUnmarshalJSON()` | Raw file data |

---

## Deprecation Notes

The following functions are still available but should use extractors instead:

- `getIntFromSlice()` → Use `ExtractInt64Array()` + manual bounds check
- `getJsonInts()` → Use `ExtractInt64Array()` 
- Manual `value.(type)` patterns → Use specific `Extract*()` function

---

## Questions?

Refer to `TYPE_SAFE_REFACTORING.md` for:
- Detailed before/after examples
- Full migration guide
- Architecture overview
- Benefits and design decisions

Or check `type_safe_extractors_test.go` for:
- Complete test examples
- Edge case handling
- Benchmark patterns
