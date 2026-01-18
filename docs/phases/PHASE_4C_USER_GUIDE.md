# Phase 4C: Marketplace User Guide

**Date:** January 16, 2026  
**Version:** 1.0  
**For:** FF6 Save Editor v4.1.0+

---

## Table of Contents

1. [Overview](#overview)
2. [Getting Started](#getting-started)
3. [Discovering Plugins](#discovering-plugins)
4. [Installing Plugins](#installing-plugins)
5. [Managing Plugins](#managing-plugins)
6. [Rating & Reviewing](#rating--reviewing)
7. [Updates](#updates)
8. [Troubleshooting](#troubleshooting)

---

## Overview

The **Marketplace** is a community-driven plugin repository where users can discover, download, install, and rate plugins created by the FF6 Save Editor community.

### What You Can Find

- **Tools:** Advanced utilities for save editing (stat optimizers, damage calculators, etc.)
- **Utilities:** General-purpose helpers (batch operations, data exporters, etc.)
- **Game Enhancements:** Speedrun configurations, achievement trackers, etc.
- **Community:** User-created templates, presets, and guides

### Key Features

✅ Browse 50+ community plugins  
✅ Full-text search across plugin names, descriptions, and authors  
✅ Filter by category, tags, and ratings  
✅ One-click installation  
✅ Automatic update notifications  
✅ Community ratings and reviews  
✅ Secure downloads with checksum verification  

---

## Getting Started

### Opening the Marketplace

1. Open **FF6 Save Editor**
2. Click **Tools** in the menu bar
3. Select **Marketplace Browser**

The Marketplace Browser window will open with the latest plugins.

### First Time Setup

- The first load may take a few seconds as it downloads the plugin catalog
- An internet connection is required for browsing and installing
- Plugins are cached locally, so subsequent browsing is faster

---

## Discovering Plugins

### Browsing Available Plugins

The main Marketplace Browser shows:
- **Plugin List:** All available plugins with name, author, and rating
- **Search Bar:** Quick-search across all plugins
- **Category Filter:** Filter by plugin type
- **Sort Options:** Sort by rating, downloads, recently updated, or alphabetical

### Using Search

**To search for a plugin:**

1. Click the **Search** field at the top
2. Type your query (e.g., "speedrun", "damage", "tracker")
3. Press Enter or click Search
4. Results appear instantly

**Search tips:**
- Search by plugin name: "speedrun tracker"
- Search by author: "by:AuthorName"
- Search by category: "category:Tools"
- Combine filters: "speedrun rating:4.0+"

### Filtering by Category

**To filter by category:**

1. Click the **Category** dropdown
2. Select a category:
   - **All** - Show all plugins
   - **Tools** - Advanced utilities and calculators
   - **Utilities** - General helpers and converters
   - **Game** - Game-specific enhancements
   - **Community** - User-created content

### Sorting Options

**Available sort options:**

- **Rating** - Highest rated first (default)
- **Downloads** - Most popular first
- **Recent** - Recently updated first
- **Name** - Alphabetical order

### Viewing Plugin Details

**To see full plugin information:**

1. Click on a plugin in the list
2. The **Details Panel** opens showing:
   - Full description
   - Author and license information
   - Version history
   - Community ratings and reviews
   - Installation button

---

## Installing Plugins

### One-Click Installation

**To install a plugin:**

1. Find the plugin in the marketplace
2. Click the **Install** button
3. A progress bar shows the download status
4. Once complete, you'll see a success message

**That's it!** The plugin is automatically enabled and ready to use.

### Manual Installation

If you prefer to download plugins manually:

1. Click the plugin's **Download** button (if available)
2. Save the `.lua` file to your computer
3. Open **Tools → Plugin Manager**
4. Click **Load Plugin** and select your downloaded file

### Installing Specific Versions

**To install a previous version:**

1. Click on the plugin
2. In the Details Panel, expand **Version History**
3. Click **Install** next to the desired version
4. Confirm the downgrade

---

## Managing Plugins

### Viewing Installed Plugins

To see all installed plugins from the marketplace:

1. Open **Tools → Plugin Manager**
2. The **Installed** tab shows all active plugins
3. Each plugin shows:
   - Name and version
   - Status (enabled/disabled)
   - Auto-update setting

### Enabling/Disabling Plugins

**To temporarily disable a plugin:**

1. Open **Tools → Plugin Manager**
2. Find the plugin in the installed list
3. Click the **Disable** toggle
4. The plugin is now inactive but not uninstalled

**To re-enable:**
- Click the **Enable** toggle

### Auto-Update Setting

**To enable automatic updates:**

1. Open **Tools → Plugin Manager**
2. Find the plugin
3. Enable **Auto-Update** toggle
4. The plugin will update automatically (when available)

### Uninstalling Plugins

**To uninstall a plugin:**

1. Open **Tools → Plugin Manager**
2. Find the plugin in the installed list
3. Click the **Uninstall** button
4. Confirm the uninstallation

**Note:** Uninstalling does not affect your save files.

---

## Rating & Reviewing

### Why Rate Plugins?

Community ratings help other users find quality plugins and provide feedback to developers.

### Submitting a Rating

**To rate a plugin:**

1. Find the plugin in the marketplace
2. Click on it to open Details
3. Scroll to the **Ratings** section
4. Click the **star rating** (1-5 stars)
5. Optionally, add a written review
6. Click **Submit Rating**

**Your rating is anonymous** - no personal information is collected.

### What Makes a Good Review?

✅ **Be specific** - Describe what you like or dislike  
✅ **Be honest** - Share your genuine experience  
✅ **Be constructive** - Suggest improvements if applicable  
✅ **Be concise** - Keep it brief but informative  

### Examples

**Good reviews:**
- "Perfect speedrun tracking! Updates the UI in real-time."
- "Useful for batch stat edits, but crashed on large files."
- "Great UI, but could use dark theme support."

**Unhelpful reviews:**
- "Good" (too vague)
- "This sucks" (not constructive)
- "5 stars" (no explanation)

### Viewing Community Reviews

**To see what others think:**

1. Click on a plugin
2. Scroll to the **Community Reviews** section
3. Read recent reviews
4. Click **Helpful** on reviews that helped you
5. Click **Report** if a review is inappropriate

---

## Updates

### Automatic Update Notifications

When updates are available:

1. A notification appears in the **Status Bar**
2. The marketplace shows a "New Version Available" badge
3. Update information includes:
   - New version number
   - Release date
   - Change notes
   - Upgrade button

### Checking for Updates

**To manually check for updates:**

1. Open **Tools → Marketplace Browser**
2. Click the **Check for Updates** button
3. Available updates are highlighted
4. Click **Update** to install

### Updating Plugins

**To update a plugin:**

1. Click the **Update** button on the plugin
2. Or, in Plugin Manager, click **Update All**
3. Progress is shown as plugins download and install
4. Previously enabled plugins restart automatically

### Downgrading Plugins

**If you need an older version:**

1. Open plugin Details
2. Expand **Version History**
3. Click **Install** next to a previous version
4. Confirm the downgrade

**Your settings are preserved** when downgrading.

---

## Troubleshooting

### Plugin Won't Install

**Error: "Network Connection Failed"**
- Check your internet connection
- Wait a moment and try again
- If persistent, the marketplace server may be down

**Error: "Checksum Mismatch"**
- The plugin file may be corrupted
- Try installing again
- Contact plugin developer if it persists

**Error: "Incompatible Version"**
- The plugin requires a newer version of FF6 Editor
- Update FF6 Editor to the latest version
- Then try installing again

### Plugin Doesn't Work After Installation

**Plugin appears enabled but doesn't function:**

1. Check the plugin requirements in Details
2. Ensure FF6 Editor version meets minimum
3. Try disabling and re-enabling the plugin
4. Check the **Plugin Manager** for error messages

**Plugin causes crashes:**

1. Disable the plugin immediately
2. Open **Tools → Plugin Manager**
3. Click **Disable** on the problematic plugin
4. Report the issue to the plugin developer

### Search Isn't Finding Plugins

**Plugins not appearing in search:**

1. Try a broader search term
2. Check that you haven't set too restrictive filters
3. Click **Clear Filters** to reset
4. Try refreshing the plugin list

### Slow Performance

**Marketplace browser is slow:**

1. You may have a slow internet connection
2. The plugin catalog may be loading
3. Try limiting search results (set a category filter)
4. Reduce the number of concurrent downloads

**Solutions:**
- Close other applications using the internet
- Wait for the initial plugin list to cache
- Try again at a different time

### Rating Won't Submit

**Error submitting a rating:**

1. Check your internet connection
2. Ensure you completed the rating form
3. Try submitting again
4. Contact support if it continues

---

## Best Practices

### Plugin Selection

✅ **Do:**
- Read plugin descriptions carefully
- Check compatibility requirements
- Review community ratings
- Start with highly-rated plugins

❌ **Don't:**
- Install every plugin at once
- Ignore version requirements
- Disable built-in editor features

### Security

✅ **Do:**
- Only install from the official marketplace
- Keep plugins updated
- Report suspicious behavior

❌ **Don't:**
- Modify downloaded plugin files
- Run plugins from unknown sources
- Ignore security warnings

### Performance

✅ **Do:**
- Install only plugins you need
- Disable unused plugins
- Monitor editor performance
- Report performance issues

❌ **Don't:**
- Install too many plugins simultaneously
- Leave debugging plugins enabled
- Force-quit the editor

---

## Frequently Asked Questions

### Q: Are marketplace plugins safe?

**A:** Yes! Plugins are sandboxed and cannot access your file system or personal data. All plugins are vetted by the community through ratings and reviews.

### Q: Can I uninstall a plugin and keep using my saves?

**A:** Yes. Plugins only extend functionality—uninstalling them doesn't affect your save files or editor settings.

### Q: How do I know if a plugin works with my FF6 Editor version?

**A:** Each plugin shows compatibility information. The Plugin Manager also checks compatibility when installing.

### Q: Can I create and share my own plugins?

**A:** Yes! See the **Phase 4C Plugin Developer Guide** for instructions on creating and publishing plugins.

### Q: Why does a plugin need permissions?

**A:** Permissions control what a plugin can do. For example, "write_save" permission allows a plugin to modify your save data. You grant permissions when installing.

### Q: What happens if I disable auto-update?

**A:** You can manually update plugins later, but you'll miss security updates and new features. Auto-update is recommended.

### Q: How are plugin ratings calculated?

**A:** Ratings are the average of all community reviews (1-5 stars). Recent ratings weighted more heavily than old ones.

### Q: Can I suggest a plugin to be created?

**A:** Yes! The community forum has a "Plugin Requests" section where you can suggest ideas.

---

## Getting Help

### Resources

- **Documentation:** See PHASE_4C_API_REFERENCE.md for technical details
- **Plugin Examples:** Community examples available in each plugin's documentation
- **Community Forum:** ff6editor.community/forums/plugins
- **Issue Tracker:** github.com/ff6-save-editor/issues

### Reporting Issues

Found a problem with a plugin?

1. Collect error details:
   - FF6 Editor version
   - Plugin name and version
   - Steps to reproduce
   - Error messages

2. Report to plugin developer via:
   - Plugin's GitHub page
   - Community forum
   - Issue tracker

---

## Glossary

**Plugin** - A small program that extends FF6 Save Editor functionality  
**Registry** - Local list of installed plugins and their metadata  
**Sandbox** - Secure environment where plugins run with limited permissions  
**Checksum** - Verification code ensuring plugin files are not corrupted  
**Marketplace** - Central repository of community plugins  
**Rating** - User feedback (1-5 stars) about plugin quality  
**Version** - Specific release of a plugin (e.g., v1.2.0)  

---

## Change Log

### Version 1.0 (January 16, 2026)
- ✅ Initial Marketplace implementation
- ✅ Plugin discovery and search
- ✅ One-click installation
- ✅ Community ratings and reviews
- ✅ Automatic update detection
- ✅ Local registry management

---

**Last Updated:** January 16, 2026  
**Version:** 1.0  
**Status:** Production Ready ✅

For technical details, see [PHASE_4C_API_REFERENCE.md](PHASE_4C_API_REFERENCE.md)

---
