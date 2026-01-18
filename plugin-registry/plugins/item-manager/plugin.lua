-- ============================================================================
-- Item Manager Plugin
-- ============================================================================
-- @id: item-manager
-- @name: Item Manager
-- @version: 1.0.0
-- @author: FF6 Editor Team
-- @description: Batch inventory operations including adding, removing, and managing item quantities
-- @permissions: read_save, write_save, ui_display

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    maxQuantity = 99,           -- Maximum item quantity
    safeMode = true,            -- Confirm before applying changes
    showPreview = true,         -- Show preview of changes
    defaultCategory = "all"     -- Default item category filter
}

-- ============================================================================
-- Item Categories (Simplified)
-- ============================================================================
local itemCategories = {
    consumables = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},  -- Potions, etc.
    weapons = {100, 101, 102, 103, 104, 105},         -- Weapons
    armor = {200, 201, 202, 203, 204, 205},           -- Armor
    relics = {250, 251, 252, 253, 254, 255},          -- Relics
    keyItems = {300, 301, 302, 303, 304, 305}         -- Key items
}

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- Format item for display
local function formatItem(item)
    return string.format("Item #%d x%d", item.ItemID, item.Count)
end

-- Count total items
local function countItems(inventory)
    local total = 0
    for _, item in ipairs(inventory.Items) do
        if item.ItemID > 0 and item.Count > 0 then
            total = total + 1
        end
    end
    return total
end

-- Count total quantity
local function countQuantity(inventory)
    local total = 0
    for _, item in ipairs(inventory.Items) do
        if item.ItemID > 0 and item.Count > 0 then
            total = total + item.Count
        end
    end
    return total
end

-- Check if item exists in inventory
local function hasItem(inventory, itemID)
    for _, item in ipairs(inventory.Items) do
        if item.ItemID == itemID then
            return true
        end
    end
    return false
end

-- Get item index in inventory
local function getItemIndex(inventory, itemID)
    for i, item in ipairs(inventory.Items) do
        if item.ItemID == itemID then
            return i
        end
    end
    return nil
end

-- ============================================================================
-- Batch Operations
-- ============================================================================

-- Max all consumables
local function maxConsumables(inventory)
    local changes = 0
    for _, itemID in ipairs(itemCategories.consumables) do
        local idx = getItemIndex(inventory, itemID)
        if idx then
            if inventory.Items[idx].Count < config.maxQuantity then
                inventory.Items[idx].Count = config.maxQuantity
                changes = changes + 1
            end
        end
    end
    return inventory, changes
end

-- Add specific item
local function addItem(inventory, itemID, quantity)
    quantity = quantity or 1
    
    -- Check if item already exists
    local idx = getItemIndex(inventory, itemID)
    if idx then
        -- Update existing item
        inventory.Items[idx].Count = math.min(inventory.Items[idx].Count + quantity, config.maxQuantity)
    else
        -- Add new item (find empty slot)
        for i, item in ipairs(inventory.Items) do
            if item.ItemID == 0 or item.Count == 0 then
                inventory.Items[i].ItemID = itemID
                inventory.Items[i].Count = math.min(quantity, config.maxQuantity)
                break
            end
        end
    end
    
    return inventory
end

-- Remove specific item
local function removeItem(inventory, itemID)
    local idx = getItemIndex(inventory, itemID)
    if idx then
        inventory.Items[idx].ItemID = 0
        inventory.Items[idx].Count = 0
        return inventory, true
    end
    return inventory, false
end

-- Set item quantity
local function setItemQuantity(inventory, itemID, quantity)
    quantity = math.min(quantity, config.maxQuantity)
    local idx = getItemIndex(inventory, itemID)
    if idx then
        if quantity > 0 then
            inventory.Items[idx].Count = quantity
        else
            inventory.Items[idx].ItemID = 0
            inventory.Items[idx].Count = 0
        end
        return inventory, true
    end
    return inventory, false
end

-- Clear all items
local function clearAllItems(inventory)
    local changes = 0
    for i, item in ipairs(inventory.Items) do
        if item.ItemID > 0 then
            inventory.Items[i].ItemID = 0
            inventory.Items[i].Count = 0
            changes = changes + 1
        end
    end
    return inventory, changes
end

-- Duplicate item (max out quantity)
local function duplicateItem(inventory, itemID)
    local idx = getItemIndex(inventory, itemID)
    if idx then
        inventory.Items[idx].Count = config.maxQuantity
        return inventory, true
    end
    return inventory, false
end

-- ============================================================================
-- UI Functions
-- ============================================================================

-- Show operation menu
local function showOperationMenu(api)
    local menu = [[
Item Manager Operations:

1. Max All Consumables (set all potions/items to 99)
2. Add Specific Item (enter item ID and quantity)
3. Remove Specific Item (enter item ID)
4. Set Item Quantity (modify existing item quantity)
5. Duplicate Item (max out existing item quantity)
6. Clear All Items (remove all items - dangerous!)
7. View Inventory Summary
8. Cancel

Enter operation number (1-8):]]
    
    local choice = api:ShowInput("Item Manager", menu, "1")
    if not choice then
        return nil
    end
    
    return tonumber(choice)
end

-- Get item ID from user
local function getItemID(api, prompt)
    prompt = prompt or "Enter item ID (1-999):"
    local input = api:ShowInput("Item ID", prompt, "")
    if not input then
        return nil
    end
    
    local itemID = tonumber(input)
    if not itemID or itemID < 1 or itemID > 999 then
        api:ShowDialog("Error", "Invalid item ID. Must be between 1 and 999.")
        return nil
    end
    
    return itemID
end

-- Get quantity from user
local function getQuantity(api)
    local input = api:ShowInput("Quantity", "Enter quantity (1-99):", "1")
    if not input then
        return nil
    end
    
    local qty = tonumber(input)
    if not qty or qty < 0 or qty > 99 then
        api:ShowDialog("Error", "Invalid quantity. Must be between 0 and 99.")
        return nil
    end
    
    return qty
end

-- Show inventory summary
local function showInventorySummary(api, inventory)
    local itemCount = countItems(inventory)
    local totalQty = countQuantity(inventory)
    
    local summary = string.format([[
INVENTORY SUMMARY
================

Total Unique Items: %d
Total Item Quantity: %d

Top 10 Items:
]], itemCount, totalQty)
    
    -- Show first 10 items
    local count = 0
    for _, item in ipairs(inventory.Items) do
        if item.ItemID > 0 and item.Count > 0 and count < 10 then
            summary = summary .. string.format("\n  %s", formatItem(item))
            count = count + 1
        end
    end
    
    if itemCount > 10 then
        summary = summary .. string.format("\n  ... and %d more items", itemCount - 10)
    end
    
    api:ShowDialog("Inventory Summary", summary)
end

-- Show preview of changes
local function showPreview(api, operation, changes)
    local message = string.format([[
Operation: %s
Changes: %d item(s) will be modified

Apply these changes?
]], operation, changes)
    
    return api:ShowConfirm("Confirm Changes", message)
end

-- ============================================================================
-- Main Entry Point
-- ============================================================================
function main(api)
    -- Verify permissions
    if not api:HasPermission("read_save") then
        api:Log("error", "Missing required permission: read_save")
        api:ShowDialog("Permission Error", "This plugin requires 'read_save' permission.")
        return false
    end
    
    if not api:HasPermission("write_save") then
        api:Log("error", "Missing required permission: write_save")
        api:ShowDialog("Permission Error", "This plugin requires 'write_save' permission.")
        return false
    end
    
    if not api:HasPermission("ui_display") then
        api:Log("error", "Missing required permission: ui_display")
        return false
    end
    
    api:Log("info", "Item Manager plugin started")
    
    -- Get current inventory
    local inventory = api:GetInventory()
    if not inventory or not inventory.Items then
        api:ShowDialog("Error", "Failed to retrieve inventory data.")
        api:Log("error", "GetInventory returned invalid data")
        return false
    end
    
    -- Show operation menu
    local operation = showOperationMenu(api)
    if not operation then
        api:Log("info", "User canceled operation")
        return false
    end
    
    local modified = false
    local changes = 0
    local operationName = "Unknown"
    
    -- Execute selected operation
    if operation == 1 then
        -- Max all consumables
        operationName = "Max All Consumables"
        inventory, changes = maxConsumables(inventory)
        modified = true
        
    elseif operation == 2 then
        -- Add specific item
        operationName = "Add Item"
        local itemID = getItemID(api, "Enter item ID to add (1-999):")
        if not itemID then return false end
        
        local qty = getQuantity(api)
        if not qty then return false end
        
        inventory = addItem(inventory, itemID, qty)
        changes = 1
        modified = true
        
    elseif operation == 3 then
        -- Remove specific item
        operationName = "Remove Item"
        local itemID = getItemID(api, "Enter item ID to remove (1-999):")
        if not itemID then return false end
        
        local success
        inventory, success = removeItem(inventory, itemID)
        if success then
            changes = 1
            modified = true
        else
            api:ShowDialog("Error", string.format("Item #%d not found in inventory.", itemID))
            return false
        end
        
    elseif operation == 4 then
        -- Set item quantity
        operationName = "Set Item Quantity"
        local itemID = getItemID(api, "Enter item ID to modify (1-999):")
        if not itemID then return false end
        
        local qty = getQuantity(api)
        if not qty then return false end
        
        local success
        inventory, success = setItemQuantity(inventory, itemID, qty)
        if success then
            changes = 1
            modified = true
        else
            api:ShowDialog("Error", string.format("Item #%d not found in inventory.", itemID))
            return false
        end
        
    elseif operation == 5 then
        -- Duplicate item
        operationName = "Duplicate Item"
        local itemID = getItemID(api, "Enter item ID to duplicate (1-999):")
        if not itemID then return false end
        
        local success
        inventory, success = duplicateItem(inventory, itemID)
        if success then
            changes = 1
            modified = true
        else
            api:ShowDialog("Error", string.format("Item #%d not found in inventory.", itemID))
            return false
        end
        
    elseif operation == 6 then
        -- Clear all items
        operationName = "Clear All Items"
        local confirm = api:ShowConfirm("DANGER", 
            "This will remove ALL items from your inventory!\n\nAre you absolutely sure?")
        if not confirm then
            api:Log("info", "User canceled clear all operation")
            return false
        end
        
        inventory, changes = clearAllItems(inventory)
        modified = true
        
    elseif operation == 7 then
        -- View inventory summary
        showInventorySummary(api, inventory)
        return true
        
    else
        -- Cancel
        api:Log("info", "User canceled operation")
        return false
    end
    
    -- Show preview and confirm
    if modified and config.safeMode and config.showPreview then
        if not showPreview(api, operationName, changes) then
            api:Log("info", "User canceled changes")
            return false
        end
    end
    
    -- Apply changes
    if modified then
        api:SetInventory(inventory)
        api:ShowDialog("Success", 
            string.format("Operation complete!\n\n%s: %d item(s) modified", operationName, changes))
        api:Log("info", string.format("Applied %s: %d changes", operationName, changes))
    end
    
    return true
end

-- ============================================================================
-- Error Handling Wrapper
-- ============================================================================
local function safeMain(api)
    local success, result = pcall(main, api)
    if not success then
        local errorMsg = tostring(result)
        api:Log("error", "Plugin error: " .. errorMsg)
        api:ShowDialog("Plugin Error", 
            "An error occurred while running the plugin:\n\n" .. errorMsg .. 
            "\n\nCheck the plugin log for details.")
        return false
    end
    return result
end

-- Execute with error handling
return safeMain(api)
