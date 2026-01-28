#!/usr/bin/env python3
"""
Plugin Validation Script for FF6 Save Editor Plugin Registry

This script validates plugin submissions against the plugin specification.
It checks:
- Directory structure
- Required files
- metadata.json schema compliance
- Lua syntax
- Checksum generation
- Documentation completeness
- Security patterns

Usage:
    python validate-plugin.py <plugin-directory>
    
Example:
    python validate-plugin.py plugins/stats-display
"""

import os
import sys
import json
import hashlib
import re
from pathlib import Path
from typing import Dict, List, Tuple, Optional

# ============================================================================
# Configuration
# ============================================================================

REQUIRED_FILES = [
    "plugin.lua",
    "metadata.json",
    "README.md",
    "CHANGELOG.md"
]

OPTIONAL_FILES = [
    "LICENSE",
    "checksum.sha256"
]

MAX_FILE_SIZES = {
    "plugin.lua": 10 * 1024 * 1024,  # 10MB
    "metadata.json": 100 * 1024,      # 100KB
    "README.md": 1 * 1024 * 1024,     # 1MB
    "CHANGELOG.md": 500 * 1024,       # 500KB
    "screenshot.png": 2 * 1024 * 1024 # 2MB
}

VALID_CATEGORIES = ["utility", "automation", "analysis", "enhancement", "debug"]
VALID_PERMISSIONS = ["read_save", "write_save", "ui_display", "events"]

FORBIDDEN_PATTERNS = [
    r'\bos\.',           # OS module access
    r'\bio\.',           # IO module access
    r'\bdebug\.',        # Debug module access
    r'\bpackage\.',      # Package module access
    r'\brequire\b',      # Require calls
    r'\bload\b',         # Load function
    r'\bloadfile\b',     # Loadfile function
    r'\bloadstring\b',   # Loadstring function
    r'\bdofile\b',       # Dofile function
    r'\bsetfenv\b',      # Setfenv function
    r'\bgetfenv\b',      # Getfenv function
]

# ============================================================================
# Validation Classes
# ============================================================================

class ValidationResult:
    """Result of a validation check"""
    def __init__(self, passed: bool, message: str, severity: str = "error"):
        self.passed = passed
        self.message = message
        self.severity = severity  # "error", "warning", "info"
    
    def __str__(self):
        icon = "✅" if self.passed else ("⚠️" if self.severity == "warning" else "❌")
        return f"{icon} {self.message}"

class PluginValidator:
    """Main plugin validator class"""
    
    def __init__(self, plugin_dir: Path):
        self.plugin_dir = plugin_dir
        self.plugin_id = plugin_dir.name
        self.results: List[ValidationResult] = []
        self.errors = 0
        self.warnings = 0
        
    def validate(self) -> bool:
        """Run all validations and return overall pass/fail"""
        print(f"\n🔍 Validating plugin: {self.plugin_id}")
        print("=" * 60)
        
        # Run validations
        self._validate_directory_structure()
        self._validate_file_sizes()
        metadata = self._validate_metadata()
        self._validate_lua_syntax()
        self._validate_lua_metadata()
        self._validate_security()
        self._validate_documentation()
        self._generate_checksum()
        
        # Print results
        self._print_results()
        
        return self.errors == 0
    
    def _add_result(self, passed: bool, message: str, severity: str = "error"):
        """Add validation result"""
        result = ValidationResult(passed, message, severity)
        self.results.append(result)
        if not passed:
            if severity == "error":
                self.errors += 1
            elif severity == "warning":
                self.warnings += 1
    
    def _validate_directory_structure(self):
        """Validate directory structure and required files"""
        print("\n📁 Checking directory structure...")
        
        # Check if directory exists
        if not self.plugin_dir.exists():
            self._add_result(False, f"Plugin directory not found: {self.plugin_dir}")
            return
        
        self._add_result(True, "Plugin directory exists")
        
        # Check if it's a directory
        if not self.plugin_dir.is_dir():
            self._add_result(False, f"Not a directory: {self.plugin_dir}")
            return
        
        # Check required files
        for filename in REQUIRED_FILES:
            filepath = self.plugin_dir / filename
            if filepath.exists():
                self._add_result(True, f"Required file found: {filename}")
            else:
                self._add_result(False, f"Missing required file: {filename}")
    
    def _validate_file_sizes(self):
        """Validate file sizes are within limits"""
        print("\n📏 Checking file sizes...")
        
        for filename, max_size in MAX_FILE_SIZES.items():
            filepath = self.plugin_dir / filename
            if filepath.exists():
                size = filepath.stat().st_size
                if size <= max_size:
                    size_mb = size / (1024 * 1024)
                    max_mb = max_size / (1024 * 1024)
                    self._add_result(True, f"{filename}: {size_mb:.2f}MB (max: {max_mb:.0f}MB)")
                else:
                    size_mb = size / (1024 * 1024)
                    max_mb = max_size / (1024 * 1024)
                    self._add_result(False, f"{filename}: {size_mb:.2f}MB exceeds max {max_mb:.0f}MB")
    
    def _validate_metadata(self) -> Optional[Dict]:
        """Validate metadata.json file"""
        print("\n📋 Validating metadata.json...")
        
        metadata_path = self.plugin_dir / "metadata.json"
        if not metadata_path.exists():
            self._add_result(False, "metadata.json not found")
            return None
        
        try:
            with open(metadata_path, 'r', encoding='utf-8') as f:
                metadata = json.load(f)
        except json.JSONDecodeError as e:
            self._add_result(False, f"Invalid JSON in metadata.json: {e}")
            return None
        
        self._add_result(True, "Valid JSON format")
        
        # Validate required fields
        required_fields = [
            "id", "name", "version", "author", "contact", "description",
            "longDescription", "category", "tags", "permissions", 
            "minEditorVersion", "homepage"
        ]
        
        for field in required_fields:
            if field in metadata:
                self._add_result(True, f"Required field present: {field}")
            else:
                self._add_result(False, f"Missing required field: {field}")
        
        # Validate plugin ID
        if "id" in metadata:
            plugin_id = metadata["id"]
            if plugin_id == self.plugin_id:
                self._add_result(True, f"Plugin ID matches directory name: {plugin_id}")
            else:
                self._add_result(False, f"Plugin ID '{plugin_id}' doesn't match directory '{self.plugin_id}'")
            
            # Validate ID format
            if re.match(r'^[a-z0-9]+(-[a-z0-9]+)*$', plugin_id):
                self._add_result(True, "Plugin ID format valid (lowercase, hyphens)")
            else:
                self._add_result(False, "Plugin ID must be lowercase alphanumeric with hyphens only")
        
        # Validate version format
        if "version" in metadata:
            version = metadata["version"]
            if re.match(r'^[0-9]+\.[0-9]+\.[0-9]+$', version):
                self._add_result(True, f"Version format valid: {version}")
            else:
                self._add_result(False, f"Version must be semantic (X.Y.Z): {version}")
        
        # Validate category
        if "category" in metadata:
            category = metadata["category"]
            if category in VALID_CATEGORIES:
                self._add_result(True, f"Category valid: {category}")
            else:
                self._add_result(False, f"Invalid category: {category}. Must be one of: {', '.join(VALID_CATEGORIES)}")
        
        # Validate permissions
        if "permissions" in metadata:
            permissions = metadata["permissions"]
            if isinstance(permissions, list) and len(permissions) > 0:
                invalid_perms = [p for p in permissions if p not in VALID_PERMISSIONS]
                if not invalid_perms:
                    self._add_result(True, f"Permissions valid: {', '.join(permissions)}")
                else:
                    self._add_result(False, f"Invalid permissions: {', '.join(invalid_perms)}")
            else:
                self._add_result(False, "Permissions must be non-empty array")
        
        # Validate description lengths
        if "description" in metadata:
            desc_len = len(metadata["description"])
            if 20 <= desc_len <= 200:
                self._add_result(True, f"Description length valid: {desc_len} chars")
            else:
                self._add_result(False, f"Description must be 20-200 characters: {desc_len} chars")
        
        if "longDescription" in metadata:
            long_desc_len = len(metadata["longDescription"])
            if 50 <= long_desc_len <= 2000:
                self._add_result(True, f"Long description length valid: {long_desc_len} chars")
            else:
                self._add_result(False, f"Long description must be 50-2000 characters: {long_desc_len} chars")
        
        # Validate tags
        if "tags" in metadata:
            tags = metadata["tags"]
            if isinstance(tags, list) and 1 <= len(tags) <= 5:
                self._add_result(True, f"Tag count valid: {len(tags)}")
            else:
                self._add_result(False, f"Tags must have 1-5 items: {len(tags) if isinstance(tags, list) else 0}")
        
        return metadata
    
    def _validate_lua_syntax(self):
        """Validate Lua syntax (basic check)"""
        print("\n🐛 Checking Lua syntax...")
        
        lua_path = self.plugin_dir / "plugin.lua"
        if not lua_path.exists():
            return
        
        try:
            with open(lua_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Check for main function
            if re.search(r'function\s+main\s*\(', content):
                self._add_result(True, "main() function defined")
            else:
                self._add_result(False, "main() function not found")
            
            # Check for balanced parentheses, brackets, and quotes (basic)
            if content.count('(') == content.count(')'):
                self._add_result(True, "Parentheses balanced")
            else:
                self._add_result(False, "Unbalanced parentheses")
            
            if content.count('[') == content.count(']'):
                self._add_result(True, "Brackets balanced")
            else:
                self._add_result(False, "Unbalanced brackets")
            
            # Check for common Lua keywords
            if 'function' in content and ('local' in content or 'return' in content):
                self._add_result(True, "Contains Lua code structure")
            else:
                self._add_result(False, "Doesn't appear to be valid Lua code", "warning")
                
        except Exception as e:
            self._add_result(False, f"Error reading Lua file: {e}")
    
    def _validate_lua_metadata(self):
        """Validate metadata comments in Lua file"""
        print("\n📝 Checking Lua metadata comments...")
        
        lua_path = self.plugin_dir / "plugin.lua"
        if not lua_path.exists():
            return
        
        try:
            with open(lua_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Check for required metadata comments
            required_metadata = ["@id:", "@name:", "@version:", "@author:", "@description:", "@permissions:"]
            
            for meta in required_metadata:
                if meta in content:
                    self._add_result(True, f"Metadata comment found: {meta}")
                else:
                    self._add_result(False, f"Missing metadata comment: {meta}", "warning")
                    
        except Exception as e:
            self._add_result(False, f"Error reading Lua metadata: {e}")
    
    def _validate_security(self):
        """Check for forbidden patterns in Lua code"""
        print("\n🔒 Checking security patterns...")
        
        lua_path = self.plugin_dir / "plugin.lua"
        if not lua_path.exists():
            return
        
        try:
            with open(lua_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            found_violations = []
            for pattern in FORBIDDEN_PATTERNS:
                matches = re.findall(pattern, content, re.MULTILINE)
                if matches:
                    found_violations.append(pattern)
            
            if not found_violations:
                self._add_result(True, "No forbidden patterns detected")
            else:
                for violation in found_violations:
                    self._add_result(False, f"Forbidden pattern found: {violation}")
                    
        except Exception as e:
            self._add_result(False, f"Error checking security: {e}")
    
    def _validate_documentation(self):
        """Validate documentation files"""
        print("\n📖 Checking documentation...")
        
        # Check README.md
        readme_path = self.plugin_dir / "README.md"
        if readme_path.exists():
            try:
                with open(readme_path, 'r', encoding='utf-8') as f:
                    content = f.read()
                
                # Check for required sections (basic check)
                required_sections = ["Installation", "Usage", "Permissions"]
                missing_sections = [s for s in required_sections if s.lower() not in content.lower()]
                
                if not missing_sections:
                    self._add_result(True, "README.md has required sections")
                else:
                    self._add_result(False, f"README.md missing sections: {', '.join(missing_sections)}", "warning")
                    
            except Exception as e:
                self._add_result(False, f"Error reading README.md: {e}")
        
        # Check CHANGELOG.md
        changelog_path = self.plugin_dir / "CHANGELOG.md"
        if changelog_path.exists():
            try:
                with open(changelog_path, 'r', encoding='utf-8') as f:
                    content = f.read()
                
                # Check for version entries
                if re.search(r'\[[\d.]+\]', content):
                    self._add_result(True, "CHANGELOG.md has version entries")
                else:
                    self._add_result(False, "CHANGELOG.md missing version entries", "warning")
                    
            except Exception as e:
                self._add_result(False, f"Error reading CHANGELOG.md: {e}")
    
    def _generate_checksum(self):
        """Generate SHA256 checksum for plugin.lua"""
        print("\n🔐 Generating checksum...")
        
        lua_path = self.plugin_dir / "plugin.lua"
        if not lua_path.exists():
            return
        
        try:
            with open(lua_path, 'rb') as f:
                content = f.read()
                checksum = hashlib.sha256(content).hexdigest()
            
            # Write checksum file
            checksum_path = self.plugin_dir / "checksum.sha256"
            with open(checksum_path, 'w') as f:
                f.write(checksum)
            
            self._add_result(True, f"Checksum generated: {checksum[:16]}...")
            
        except Exception as e:
            self._add_result(False, f"Error generating checksum: {e}")
    
    def _print_results(self):
        """Print validation results summary"""
        print("\n" + "=" * 60)
        print("📊 VALIDATION RESULTS")
        print("=" * 60)
        
        for result in self.results:
            print(result)
        
        print("\n" + "=" * 60)
        print(f"Total checks: {len(self.results)}")
        print(f"Errors: {self.errors}")
        print(f"Warnings: {self.warnings}")
        print("=" * 60)
        
        if self.errors == 0:
            print("\n✅ Plugin validation PASSED!")
            if self.warnings > 0:
                print(f"⚠️  {self.warnings} warning(s) detected. Review recommended.")
        else:
            print(f"\n❌ Plugin validation FAILED with {self.errors} error(s)")
            print("Fix the errors and run validation again.")

# ============================================================================
# Main
# ============================================================================

def main():
    """Main entry point"""
    if len(sys.argv) < 2:
        print("Usage: python validate-plugin.py <plugin-directory>")
        print("\nExample:")
        print("  python validate-plugin.py plugins/stats-display")
        sys.exit(1)
    
    plugin_dir = Path(sys.argv[1])
    
    validator = PluginValidator(plugin_dir)
    success = validator.validate()
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
