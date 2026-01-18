# Contributing to the Plugin Registry

Thank you for your interest in contributing plugins to the FF6 Save Editor! This guide will walk you through the process of creating and submitting a plugin.

---

## Table of Contents

1. [Before You Start](#before-you-start)
2. [Plugin Development](#plugin-development)
3. [Testing Your Plugin](#testing-your-plugin)
4. [Submission Process](#submission-process)
5. [Review Process](#review-process)
6. [Plugin Standards](#plugin-standards)
7. [Common Issues](#common-issues)

---

## Before You Start

### Prerequisites

- **Editor Version:** FF6 Save Editor 3.4.0 or higher
- **Development Tools:** Any text editor (VS Code, Sublime, Notepad++)
- **Testing:** Access to test save files
- **GitHub Account:** Required for submitting plugins

### Required Reading

Please familiarize yourself with:
- [Plugin Specification](PLUGIN_SPECIFICATION.md) - Plugin format requirements
- [Plugin API Reference](https://github.com/ff6-save-editor/editor/blob/main/PHASE_4B_API_REFERENCE.md) - Complete API documentation
- [Plugin Development Guide](https://github.com/ff6-save-editor/editor/blob/main/PHASE_4B_PLUGIN_GUIDE.md) - Best practices and examples

---

## Plugin Development

### 1. Choose a Plugin Idea

Good plugin ideas:
- ✅ Solve a specific problem or need
- ✅ Provide clear value to users
- ✅ Have a well-defined scope
- ✅ Complement existing editor features

Avoid:
- ❌ Duplicating existing plugins
- ❌ Plugins that are too broad or vague
- ❌ Plugins that could corrupt save data
- ❌ Plugins requiring external dependencies

### 2. Plan Your Plugin

Before coding, document:
- **Purpose:** What does your plugin do?
- **Features:** List of specific features
- **Permissions:** What permissions are needed?
- **UI/UX:** How will users interact with it?
- **Testing:** How will you test it?

### 3. Create Plugin Structure

```
plugins/your-plugin-name/
├── plugin.lua              # Main plugin code (required)
├── metadata.json           # Plugin metadata (required)
├── README.md               # Plugin documentation (required)
├── CHANGELOG.md            # Version history (required)
├── LICENSE                 # Plugin license (optional, defaults to MIT)
├── screenshot.png          # Plugin screenshot (recommended)
└── checksum.sha256         # SHA256 checksum (auto-generated)
```

### 4. Write Plugin Code

**plugin.lua Template:**

```lua
-- ============================================================================
-- Plugin Metadata (Required)
-- ============================================================================
-- @id: your-plugin-name
-- @name: Your Plugin Name
-- @version: 1.0.0
-- @author: Your Name
-- @description: Brief description of what your plugin does
-- @permissions: read_save, ui_display

-- ============================================================================
-- Plugin Configuration (Optional)
-- ============================================================================
local config = {
    -- Add any configuration options here
    debug = false,
    maxResults = 10
}

-- ============================================================================
-- Helper Functions
-- ============================================================================
local function formatHP(current, max)
    return string.format("%d/%d (%.1f%%)", current, max, (current / max) * 100)
end

-- ============================================================================
-- Main Entry Point (Required)
-- ============================================================================
function main(api)
    -- Check permissions
    if not api:HasPermission("read_save") then
        api:Log("error", "Missing required permission: read_save")
        return false
    end
    
    -- Your plugin logic here
    local char = api:GetCharacter(0)
    if not char then
        api:ShowDialog("Error", "No character found")
        return false
    end
    
    -- Display results
    local message = string.format(
        "Character: %s\nLevel: %d\nHP: %s",
        char.Name,
        char.Level,
        formatHP(char.CurrentHP, char.MaxHP)
    )
    
    api:ShowDialog("Character Info", message)
    
    -- Return true on success
    return true
end

-- ============================================================================
-- Error Handling
-- ============================================================================
-- Wrap main in pcall for error handling
local success, result = pcall(main, api)
if not success then
    api:Log("error", "Plugin error: " .. tostring(result))
    api:ShowDialog("Error", "Plugin encountered an error. Check logs for details.")
    return false
end

return result
```

### 5. Create metadata.json

```json
{
  "id": "your-plugin-name",
  "name": "Your Plugin Name",
  "version": "1.0.0",
  "author": "Your Name",
  "contact": "your.email@example.com",
  "description": "Brief description of your plugin (max 200 characters)",
  "longDescription": "Detailed description of your plugin's features, use cases, and capabilities. Can be multiple paragraphs.",
  "category": "utility",
  "tags": ["stats", "character", "analysis"],
  "permissions": ["read_save", "ui_display"],
  "minEditorVersion": "3.4.0",
  "homepage": "https://github.com/yourusername/plugin-name",
  "repository": "https://github.com/yourusername/plugin-name",
  "license": "MIT",
  "dependencies": [],
  "screenshots": ["screenshot.png"],
  "documentation": "README.md",
  "changelog": "CHANGELOG.md"
}
```

**Field Descriptions:**

- **id:** Unique identifier (lowercase, hyphens only)
- **name:** Human-readable name
- **version:** Semantic version (MAJOR.MINOR.PATCH)
- **author:** Your name or GitHub username
- **contact:** Email or GitHub profile URL
- **description:** Brief description (shown in marketplace)
- **longDescription:** Detailed description (shown in plugin details)
- **category:** One of: utility, automation, analysis, enhancement, debug
- **tags:** Array of relevant keywords (max 5)
- **permissions:** Array of required permissions
- **minEditorVersion:** Minimum required editor version
- **homepage:** Plugin homepage or repository URL
- **repository:** Source code repository URL
- **license:** License identifier (MIT, GPL-3.0, etc.)
- **dependencies:** Array of plugin IDs this plugin depends on
- **screenshots:** Array of screenshot filenames
- **documentation:** Path to README.md
- **changelog:** Path to CHANGELOG.md

### 6. Write Documentation

**README.md Template:**

```markdown
# Your Plugin Name

Brief description of your plugin.

## Features

- Feature 1: Description
- Feature 2: Description
- Feature 3: Description

## Installation

1. Open FF6 Save Editor
2. Go to Tools → Marketplace
3. Search for "Your Plugin Name"
4. Click "Install"

## Usage

1. Load a save file
2. Go to Tools → Plugin Manager
3. Select "Your Plugin Name"
4. Click "Run"

## Screenshots

![Screenshot 1](screenshot.png)

## Permissions

This plugin requires the following permissions:
- **read_save:** To read character data
- **ui_display:** To show results in a dialog

## Configuration

This plugin supports the following configuration options:
- `maxResults`: Maximum number of results to display (default: 10)

## Compatibility

- **Editor Version:** 3.4.0 or higher
- **Game:** Final Fantasy VI Pixel Remastered

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## License

MIT License - See [LICENSE](LICENSE) for details.

## Support

For issues or questions:
- GitHub Issues: [your-repo-url]
- Discord: @yourusername
```

**CHANGELOG.md Template:**

```markdown
# Changelog

All notable changes to this plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-16

### Added
- Initial release
- Feature 1
- Feature 2
- Feature 3

### Changed
- N/A (initial release)

### Deprecated
- N/A (initial release)

### Removed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- N/A (initial release)
```

---

## Testing Your Plugin

### Local Testing

1. **Install Plugin Locally:**
   ```
   Copy plugin.lua to: %USERPROFILE%\.ff6editor\plugins\your-plugin-name\plugin.lua
   ```

2. **Test in Editor:**
   - Open FF6 Save Editor
   - Load a test save file
   - Go to Tools → Plugin Manager
   - Find your plugin in the list
   - Click "Run"
   - Verify expected behavior

3. **Test Error Handling:**
   - Test with invalid data
   - Test with missing characters
   - Test with edge cases
   - Verify error messages are clear

4. **Test Permissions:**
   - Verify plugin only uses declared permissions
   - Test with minimal permissions
   - Ensure no permission bypass

### Validation

Before submitting, run the validation script:

```bash
python scripts/validate-plugin.py plugins/your-plugin-name
```

This checks:
- ✅ Directory structure
- ✅ Required files present
- ✅ metadata.json schema compliance
- ✅ Lua syntax correctness
- ✅ Checksum generation
- ✅ Documentation completeness

---

## Submission Process

### 1. Fork Repository

```bash
git clone https://github.com/yourusername/plugin-registry.git
cd plugin-registry
```

### 2. Create Plugin Branch

```bash
git checkout -b add-your-plugin-name
```

### 3. Add Your Plugin

```bash
# Create plugin directory
mkdir -p plugins/your-plugin-name

# Copy your plugin files
cp /path/to/your/plugin.lua plugins/your-plugin-name/
cp /path/to/your/metadata.json plugins/your-plugin-name/
cp /path/to/your/README.md plugins/your-plugin-name/
cp /path/to/your/CHANGELOG.md plugins/your-plugin-name/
cp /path/to/your/screenshot.png plugins/your-plugin-name/ (optional)
```

### 4. Generate Checksum

```bash
python scripts/update-checksums.py plugins/your-plugin-name
```

### 5. Validate Plugin

```bash
python scripts/validate-plugin.py plugins/your-plugin-name
```

Fix any errors reported by the validator.

### 6. Commit Changes

```bash
git add plugins/your-plugin-name
git commit -m "Add Your Plugin Name plugin

- Brief description of features
- Version 1.0.0
- Category: utility/automation/analysis/enhancement
"
```

### 7. Push to GitHub

```bash
git push origin add-your-plugin-name
```

### 8. Create Pull Request

1. Go to your fork on GitHub
2. Click "Pull requests" → "New pull request"
3. Select your branch: `add-your-plugin-name`
4. Fill out the pull request template:

```markdown
## Plugin Submission

**Plugin Name:** Your Plugin Name
**Plugin ID:** your-plugin-name
**Version:** 1.0.0
**Category:** utility
**Author:** Your Name

### Description

Brief description of what your plugin does and why it's useful.

### Features

- Feature 1
- Feature 2
- Feature 3

### Testing

Describe how you tested your plugin:
- ✅ Tested with multiple save files
- ✅ Tested error handling
- ✅ Tested all features
- ✅ Ran validation script
- ✅ Verified permissions

### Checklist

- [x] Plugin follows [PLUGIN_SPECIFICATION.md](PLUGIN_SPECIFICATION.md)
- [x] All required files included
- [x] Documentation complete
- [x] Validation script passes
- [x] Tested locally
- [x] Screenshots included (if applicable)
- [x] License specified

### Additional Notes

Any additional information for reviewers.
```

5. Submit the pull request

---

## Review Process

### Automated Checks

When you submit a pull request, automated CI checks will run:

1. **Structure Validation:** Verify all required files are present
2. **Metadata Validation:** Check metadata.json against schema
3. **Syntax Check:** Verify Lua syntax is correct
4. **Checksum Verification:** Ensure checksum is correct
5. **Security Scan:** Check for forbidden patterns or API calls
6. **Documentation Check:** Verify README and CHANGELOG exist

### Manual Review

After automated checks pass, maintainers will review:

1. **Code Quality:**
   - Readable and well-commented
   - Follows best practices
   - Proper error handling
   - Efficient implementation

2. **Functionality:**
   - Works as described
   - No bugs or crashes
   - User-friendly interface
   - Clear error messages

3. **Security:**
   - No malicious code
   - Appropriate permission usage
   - Safe data modifications
   - No forbidden API usage

4. **Documentation:**
   - Clear and complete
   - Accurate descriptions
   - Usage examples included
   - Screenshots helpful

### Review Timeline

- **Automated checks:** Immediate (minutes)
- **Initial review:** 1-3 business days
- **Feedback/revisions:** Varies based on complexity
- **Final approval:** 1-2 business days after all checks pass

### Possible Outcomes

1. **Approved:** Plugin is merged and published
2. **Changes Requested:** Reviewer feedback provided, revisions needed
3. **Rejected:** Plugin doesn't meet standards or duplicates existing functionality

---

## Plugin Standards

### Code Quality

- ✅ Clear variable and function names
- ✅ Comprehensive comments explaining logic
- ✅ Consistent indentation (2 or 4 spaces)
- ✅ Error handling for all operations
- ✅ Input validation
- ✅ No hardcoded values (use configuration)

### Performance

- ✅ Execution time <5 seconds for typical operations
- ✅ Memory usage <10MB
- ✅ No infinite loops
- ✅ Efficient algorithms
- ✅ Avoid unnecessary API calls

### Security

- ✅ Only use declared permissions
- ✅ Validate all user input
- ✅ Sanitize data before display
- ✅ No hardcoded credentials
- ✅ No attempts to bypass sandbox

### Documentation

- ✅ README.md with clear usage instructions
- ✅ CHANGELOG.md with version history
- ✅ Inline comments explaining complex logic
- ✅ Examples of usage
- ✅ Known limitations documented

### Testing

- ✅ Tested with multiple save files
- ✅ Edge cases tested
- ✅ Error conditions tested
- ✅ Permissions tested
- ✅ No false positives/negatives

---

## Common Issues

### Validation Failures

**Issue:** "Metadata validation failed: missing required field 'category'"
**Solution:** Add the category field to metadata.json

**Issue:** "Lua syntax error on line 42"
**Solution:** Check Lua syntax at line 42, likely missing `end` or quote

**Issue:** "Checksum mismatch"
**Solution:** Regenerate checksum: `python scripts/update-checksums.py plugins/your-plugin`

### Plugin Not Loading

**Issue:** Plugin doesn't appear in Plugin Manager
**Solution:** Check plugin filename is exactly `plugin.lua` (case-sensitive)

**Issue:** "Permission denied" error
**Solution:** Add required permission to metadata.json permissions array

### API Errors

**Issue:** `api:GetCharacter(0)` returns nil
**Solution:** Check if save file has characters loaded

**Issue:** "attempt to call method 'ShowDialog' (a nil value)"
**Solution:** Ensure API object is passed to main function correctly

---

## Getting Help

### Resources

- **Discord:** Join our community Discord for real-time help
- **GitHub Discussions:** Ask questions in GitHub Discussions
- **Examples:** Review [official example plugins](https://github.com/ff6-save-editor/editor/blob/main/PHASE_4B_PLUGIN_EXAMPLES.md)

### Contact

- **Email:** plugins@ff6editor.dev
- **Discord:** discord.gg/ff6editor
- **GitHub Issues:** github.com/ff6-save-editor/plugin-registry/issues

---

## Code of Conduct

By contributing, you agree to:
- Be respectful and constructive
- Help maintain a welcoming community
- Follow all plugin guidelines and standards
- Not submit malicious code
- Respond to review feedback in a timely manner

---

## License

By submitting a plugin, you agree that:
- Your plugin code is your own or properly attributed
- You grant users the right to use your plugin under the specified license
- The registry maintainers can redistribute your plugin
- You will maintain your plugin and address critical issues

Default license for plugins without explicit LICENSE file: MIT

---

Thank you for contributing to the FF6 Save Editor plugin ecosystem!
