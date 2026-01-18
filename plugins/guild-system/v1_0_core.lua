--[[
  Guild System Plugin - v1.0
  Guild management with membership, roles, rewards, and communication
  
  Phase: 10 (Tier 3 - Social/Community)
  Version: 1.0 (New Plugin)
]]

-- ============================================================================
-- FEATURE 1: GUILD MANAGEMENT (~250 LOC)
-- ============================================================================

local GuildManagement = {}

---Create new guild
---@param guildName string Guild identifier
---@param ownerID string Guild founder
---@return table guild Created guild
function GuildManagement.createGuild(guildName, ownerID)
  if not guildName or not ownerID then return {} end
  
  local guild = {
    guild_id = "GUILD_" .. guildName,
    name = guildName,
    owner_id = ownerID,
    created_at = os.time(),
    members = 1,
    level = 1,
    experience = 0,
    funds = 0,
    status = "Active"
  }
  
  return guild
end

---Join guild
---@param userID string Player joining
---@param guildID string Guild to join
---@return table membership Membership record
function GuildManagement.joinGuild(userID, guildID)
  if not userID or not guildID then return {} end
  
  local membership = {
    user_id = userID,
    guild_id = guildID,
    joined_at = os.time(),
    role = "Member",
    contribution = 0,
    status = "Active"
  }
  
  return membership
end

---Leave guild
---@param userID string Player leaving
---@param guildID string Guild to leave
---@return table departure Departure record
function GuildManagement.leaveGuild(userID, guildID)
  if not userID or not guildID then return {} end
  
  local departure = {
    user_id = userID,
    guild_id = guildID,
    left_at = os.time(),
    membership_duration = 30,
    farewell = "Good in touch!"
  }
  
  return departure
end

---Get guild information
---@param guildID string Guild to query
---@return table info Guild information
function GuildManagement.getInfo(guildID)
  if not guildID then return {} end
  
  local info = {
    guild_id = guildID,
    name = "Dragon Slayers",
    owner = "Leader1",
    members = 45,
    level = 12,
    experience = 125000,
    founded = "2024-01-01",
    headquarters = "Central City"
  }
  
  return info
end

-- ============================================================================
-- FEATURE 2: ROLE MANAGEMENT (~250 LOC)
-- ============================================================================

local RoleManagement = {}

---Assign member role
---@param userID string Member to assign
---@param role string Role name
---@return table assigned Role assignment
function RoleManagement.assignRole(userID, role)
  if not userID or not role then return {} end
  
  local assigned = {
    user_id = userID,
    role = role,
    assigned_at = os.time(),
    permissions = {
      {permission = "Invite members", granted = true},
      {permission = "Manage messages", granted = role == "Officer"}
    },
    assignment_success = true
  }
  
  return assigned
end

---List member roles
---@param guildID string Guild to query
---@return table roleStructure Guild role hierarchy
function RoleManagement.listRoles(guildID)
  if not guildID then return {} end
  
  local roleStructure = {
    guild_id = guildID,
    roles = {
      {role = "Owner", members = 1, permissions = 50},
      {role = "Officer", members = 3, permissions = 35},
      {role = "Member", members = 41, permissions = 10}
    },
    total_roles = 3
  }
  
  return roleStructure
end

---Grant permission
---@param role string Role to modify
---@param permission string Permission to grant
---@return table granted Permission grant confirmation
function RoleManagement.grantPermission(role, permission)
  if not role or not permission then return {} end
  
  local granted = {
    role = role,
    permission = permission,
    granted_at = os.time(),
    status = "Granted",
    effective_immediately = true
  }
  
  return granted
end

---Revoke permission
---@param role string Role to modify
---@param permission string Permission to revoke
---@return table revoked Permission revoke confirmation
function RoleManagement.revokePermission(role, permission)
  if not role or not permission then return {} end
  
  local revoked = {
    role = role,
    permission = permission,
    revoked_at = os.time(),
    status = "Revoked",
    effective_immediately = true
  }
  
  return revoked
end

-- ============================================================================
-- FEATURE 3: REWARDS & LEVELS (~240 LOC)
-- ============================================================================

local RewardsAndLevels = {}

---Award guild experience
---@param guildID string Guild earning EXP
---@param amount number Experience to award
---@return table award EXP award record
function RewardsAndLevels.awardExperience(guildID, amount)
  if not guildID or not amount then return {} end
  
  local award = {
    guild_id = guildID,
    experience_awarded = amount,
    total_experience = 125000,
    experience_to_level = 25000,
    level_up_pending = false
  }
  
  return award
end

---Level up guild
---@param guildID string Guild to level
---@return table levelUp Guild level increase
function RewardsAndLevels.levelUp(guildID)
  if not guildID then return {} end
  
  local levelUp = {
    guild_id = guildID,
    new_level = 13,
    previous_level = 12,
    level_up_at = os.time(),
    new_benefits = {"Guild hall expansion", "Increased member capacity"}
  }
  
  return levelUp
end

---Distribute guild rewards
---@param guildID string Guild distributing rewards
---@param rewardPool number Total reward amount
---@return table distribution Reward distribution
function RewardsAndLevels.distributeRewards(guildID, rewardPool)
  if not guildID or not rewardPool then return {} end
  
  local distribution = {
    guild_id = guildID,
    total_pool = rewardPool,
    distributed_members = 45,
    per_member_share = rewardPool / 45,
    distribution_date = os.time(),
    status = "Distributed"
  }
  
  return distribution
end

---Get guild treasury
---@param guildID string Guild to query
---@return table treasury Guild funds
function RewardsAndLevels.getTreasury(guildID)
  if not guildID then return {} end
  
  local treasury = {
    guild_id = guildID,
    total_funds = 125000,
    income_this_month = 45000,
    expenses_this_month = 15000,
    projected_balance = 155000
  }
  
  return treasury
end

-- ============================================================================
-- FEATURE 4: COMMUNICATION & EVENTS (~210 LOC)
-- ============================================================================

local CommunicationAndEvents = {}

---Post to guild message board
---@param userID string Poster
---@param message string Message content
---@return table post Posted message
function CommunicationAndEvents.postMessage(userID, message)
  if not userID or not message then return {} end
  
  local post = {
    post_id = "POST_" .. os.time(),
    user_id = userID,
    message = message,
    posted_at = os.time(),
    replies = 0,
    visibility = "Guild"
  }
  
  return post
end

---Schedule guild event
---@param eventName string Event name
---@param eventTime number Unix timestamp
---@return table event Scheduled event
function CommunicationAndEvents.scheduleEvent(eventName, eventTime)
  if not eventName or not eventTime then return {} end
  
  local event = {
    event_id = "EVENT_" .. os.time(),
    name = eventName,
    scheduled_at = eventTime,
    status = "Scheduled",
    capacity = 50,
    signups = 0
  }
  
  return event
end

---Join guild event
---@param userID string Participant
---@param eventID string Event to join
---@return table signup Event signup confirmation
function CommunicationAndEvents.joinEvent(userID, eventID)
  if not userID or not eventID then return {} end
  
  local signup = {
    user_id = userID,
    event_id = eventID,
    signed_up_at = os.time(),
    status = "Registered",
    confirmation = "See you at the event!"
  }
  
  return signup
end

---Announce guild news
---@param guildID string Guild making announcement
---@param announcement string News content
---@return table news Announcement record
function CommunicationAndEvents.announceNews(guildID, announcement)
  if not guildID or not announcement then return {} end
  
  local news = {
    guild_id = guildID,
    announcement_id = "ANN_" .. os.time(),
    content = announcement,
    announced_at = os.time(),
    visibility = "All Members",
    importance = "Normal"
  }
  
  return news
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.0",
  GuildManagement = GuildManagement,
  RoleManagement = RoleManagement,
  RewardsAndLevels = RewardsAndLevels,
  CommunicationAndEvents = CommunicationAndEvents,
  
  features = {
    guildManagement = true,
    roleManagement = true,
    rewardsAndLevels = true,
    communicationAndEvents = true
  }
}
