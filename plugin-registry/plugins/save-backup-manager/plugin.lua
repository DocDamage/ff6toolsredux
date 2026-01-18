--[[
    Save Backup Manager
    Version: 1.0.0

    Automated backup management with timestamped backups, restore functionality,
    and backup rotation. Creates safety net before modifications.

    Permissions: read_save, write_save, ui_display, file_io
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    backupsDirectory = "backups",
    maxBackups = 10,
    timestampFormat = "%Y%m%d_%H%M%S",
}

-- ============================================================================
-- Helpers
-- ============================================================================

local function safeCall(api, method, ...)
    if not api or not method then return nil end
    local ok, result = pcall(function(...)
        local fn = api[method]
        if type(fn) == "function" then
            return fn(api, ...)
        end
        return nil
    end, ...)
    if not ok then return nil end
    return result
end

local function getTimestamp()
    return os.date(config.timestampFormat)
end

local function createBackup(api, sourcePath, label)
    local timestamp = getTimestamp()
    local filename = string.format("%s_%s.sav", label or "backup", timestamp)
    local backupPath = string.format("%s/%s", config.backupsDirectory, filename)
    
    local content = safeCall(api, "ReadFile", sourcePath)
    if not content then return nil, "Failed to read source file" end
    
    local ok = safeCall(api, "WriteFile", backupPath, content)
    if not ok then return nil, "Failed to write backup file" end
    
    return backupPath, nil
end

local function listBackups(api)
    local files = safeCall(api, "ListFiles", config.backupsDirectory)
    if not files then return {} end
    
    local backups = {}
    for _, file in ipairs(files) do
        if file:match("%.sav$") then
            table.insert(backups, file)
        end
    end
    
    -- Sort by name (includes timestamp)
    table.sort(backups, function(a, b) return a > b end)
    
    return backups
end

local function restoreBackup(api, backupFilename, targetPath)
    local backupPath = string.format("%s/%s", config.backupsDirectory, backupFilename)
    local content = safeCall(api, "ReadFile", backupPath)
    if not content then return false, "Failed to read backup" end
    
    local ok = safeCall(api, "WriteFile", targetPath, content)
    if not ok then return false, "Failed to write to target" end
    
    return true, nil
end

local function deleteBackup(api, backupFilename)
    local backupPath = string.format("%s/%s", config.backupsDirectory, backupFilename)
    local ok = safeCall(api, "DeleteFile", backupPath)
    return ok ~= nil
end

local function rotateBackups(api)
    local backups = listBackups(api)
    if #backups <= config.maxBackups then return 0 end
    
    local deleted = 0
    for i = config.maxBackups + 1, #backups do
        if deleteBackup(api, backups[i]) then
            deleted = deleted + 1
        end
    end
    
    return deleted
end

-- ============================================================================
-- Views
-- ============================================================================

local function createBackupView(api)
    local sourcePath = api:ShowInput("Create Backup", "Enter path to save file:", "")
    if not sourcePath or sourcePath == "" then return end
    
    local label = api:ShowInput("Create Backup", "Enter backup label (optional):", "manual")
    if not label or label == "" then label = "manual" end
    
    local backupPath, err = createBackup(api, sourcePath, label)
    if not backupPath then
        api:ShowDialog("Error", "Backup failed: " .. (err or "unknown error"))
        return
    end
    
    api:ShowDialog("Success", string.format("Backup created:\n%s", backupPath))
    
    -- Auto-rotate
    local deleted = rotateBackups(api)
    if deleted > 0 then
        api:ShowDialog("Info", string.format("Rotated %d old backups (max: %d)", deleted, config.maxBackups))
    end
end

local function listBackupsView(api)
    local backups = listBackups(api)
    if #backups == 0 then
        api:ShowDialog("Backup Library", "No backups found.")
        return
    end
    
    local lines = {string.format("%d backups (showing latest %d):", #backups, math.min(#backups, 20))}
    for i = 1, math.min(20, #backups) do
        table.insert(lines, string.format("%d) %s", i, backups[i]))
    end
    
    api:ShowDialog("Backup Library", table.concat(lines, "\n"))
end

local function restoreBackupView(api)
    local backups = listBackups(api)
    if #backups == 0 then
        api:ShowDialog("Restore Backup", "No backups found.")
        return
    end
    
    local list = {}
    for i, file in ipairs(backups) do
        table.insert(list, string.format("%d) %s", i, file))
    end
    
    local choice = api:ShowInput("Restore Backup", "Select backup:\n" .. table.concat(list, "\n") .. "\n\nEnter number:", "1")
    if not choice then return end
    local idx = tonumber(choice)
    if not idx or idx < 1 or idx > #backups then
        api:ShowDialog("Error", "Invalid selection.")
        return
    end
    
    local targetPath = api:ShowInput("Restore Backup", "Enter target path to restore to:", "")
    if not targetPath or targetPath == "" then return end
    
    local confirm = api:ShowInput("Restore Backup", string.format("Restore '%s' to '%s'? (y/n)", backups[idx], targetPath), "n")
    if not confirm or confirm:lower() ~= "y" then return end
    
    local ok, err = restoreBackup(api, backups[idx], targetPath)
    if ok then
        api:ShowDialog("Success", string.format("Backup restored to:\n%s", targetPath))
    else
        api:ShowDialog("Error", "Restore failed: " .. (err or "unknown error"))
    end
end

local function deleteBackupView(api)
    local backups = listBackups(api)
    if #backups == 0 then
        api:ShowDialog("Delete Backup", "No backups found.")
        return
    end
    
    local list = {}
    for i, file in ipairs(backups) do
        table.insert(list, string.format("%d) %s", i, file))
    end
    
    local choice = api:ShowInput("Delete Backup", "Select backup to delete:\n" .. table.concat(list, "\n") .. "\n\nEnter number:", "1")
    if not choice then return end
    local idx = tonumber(choice)
    if not idx or idx < 1 or idx > #backups then
        api:ShowDialog("Error", "Invalid selection.")
        return
    end
    
    local confirm = api:ShowInput("Delete Backup", string.format("Delete '%s'? (y/n)", backups[idx]), "n")
    if not confirm or confirm:lower() ~= "y" then return end
    
    if deleteBackup(api, backups[idx]) then
        api:ShowDialog("Success", string.format("Backup '%s' deleted.", backups[idx]))
    else
        api:ShowDialog("Error", "Failed to delete backup.")
    end
end

local function helpView(api)
    local lines = {
        "Save Backup Manager",
        "- Create: Timestamped backups with optional labels.",
        "- List: View all backups sorted by date.",
        "- Restore: Apply backup to target path.",
        "- Delete: Remove old backups.",
        "- Auto-rotation: Keeps last " .. config.maxBackups .. " backups.",
        "Requires: read_save, write_save, file_io permissions.",
        "Backups stored in: " .. config.backupsDirectory .. "/ directory.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Save Backup Manager",
            "Select an option:\n1) Create Backup\n2) List Backups\n3) Restore Backup\n4) Delete Backup\n5) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            createBackupView(api)
        elseif sel == 2 then
            listBackupsView(api)
        elseif sel == 3 then
            restoreBackupView(api)
        elseif sel == 4 then
            deleteBackupView(api)
        elseif sel == 5 then
            helpView(api)
        elseif sel == 0 then
            return
        else
            api:ShowDialog("Invalid", "Please choose a valid option.")
        end
    end
end

local function safeMain(api)
    local ok, err = pcall(function()
        mainMenu(api)
    end)
    if not ok and api and api.ShowDialog then
        api:ShowDialog("Error", "Save Backup Manager failed: " .. tostring(err))
    end
end

return {
    id = "save-backup-manager",
    name = "Save Backup Manager",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "write_save", "ui_display", "file_io"},
    run = safeMain
}
