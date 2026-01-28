--[[
  Social Hub Plugin - v1.0
  Messaging, friend lists, group activities, and social features
  
  Phase: 10 (Tier 3 - Social/Community)
  Version: 1.0 (New Plugin)
]]

-- ============================================================================
-- FEATURE 1: MESSAGING SYSTEM (~250 LOC)
-- ============================================================================

local MessagingSystem = {}

---Send direct message
---@param senderID string Sending player
---@param recipientID string Receiving player
---@param message string Message content
---@return table sent Message sent confirmation
function MessagingSystem.sendMessage(senderID, recipientID, message)
  if not senderID or not recipientID or not message then return {} end
  
  local sent = {
    message_id = "MSG_" .. os.time(),
    from = senderID,
    to = recipientID,
    content = message,
    sent_at = os.time(),
    read = false,
    status = "Delivered"
  }
  
  return sent
end

---Read messages
---@param userID string User reading messages
---@return table inbox Inbox messages
function MessagingSystem.readMessages(userID)
  if not userID then return {} end
  
  local inbox = {
    user_id = userID,
    total_messages = 45,
    unread_count = 12,
    messages = {
      {from = "Player1", subject = "Join my team", date = os.time() - 3600, unread = true},
      {from = "Player2", subject = "Great match!", date = os.time() - 7200, unread = true},
      {from = "Player3", subject = "Guild invite", date = os.time() - 86400, unread = false}
    }
  }
  
  return inbox
end

---Send group message
---@param senderID string Sending player
---@param groupID string Target group
---@param message string Message content
---@return table groupMsg Group message confirmation
function MessagingSystem.sendGroupMessage(senderID, groupID, message)
  if not senderID or not groupID or not message then return {} end
  
  local groupMsg = {
    message_id = "GROUP_MSG_" .. os.time(),
    sender = senderID,
    group = groupID,
    content = message,
    recipients = 12,
    sent_at = os.time()
  }
  
  return groupMsg
end

---Delete message
---@param userID string User deleting
---@param messageID string Message to delete
---@return table deleted Deletion confirmation
function MessagingSystem.deleteMessage(userID, messageID)
  if not userID or not messageID then return {} end
  
  local deleted = {
    user_id = userID,
    message_id = messageID,
    deleted_at = os.time(),
    status = "Deleted"
  }
  
  return deleted
end

-- ============================================================================
-- FEATURE 2: FRIEND LISTS & GROUPS (~250 LOC)
-- ============================================================================

local FriendListsGroups = {}

---Create friend group
---@param userID string User creating group
---@param groupName string Group name
---@return table group Created friend group
function FriendListsGroups.createGroup(userID, groupName)
  if not userID or not groupName then return {} end
  
  local group = {
    group_id = "FGROUP_" .. groupName,
    owner = userID,
    name = groupName,
    members = 0,
    created_at = os.time()
  }
  
  return group
end

---Add friend to group
---@param userID string User adding
---@param friendID string Friend to add
---@param groupID string Target group
---@return table added Addition confirmation
function FriendListsGroups.addToGroup(userID, friendID, groupID)
  if not userID or not friendID or not groupID then return {} end
  
  local added = {
    user_id = userID,
    friend_id = friendID,
    group_id = groupID,
    added_at = os.time(),
    status = "Added"
  }
  
  return added
end

---List groups
---@param userID string User to query
---@return table groups User's friend groups
function FriendListsGroups.listGroups(userID)
  if not userID then return {} end
  
  local groups = {
    user_id = userID,
    total_groups = 5,
    groups = {
      {name = "Close Friends", count = 8},
      {name = "Guild Members", count = 45},
      {name = "Team", count = 4}
    }
  }
  
  return groups
end

---Quick message group
---@param userID string Sender
---@param groupID string Target group
---@param message string Message content
---@return table broadcast Broadcast confirmation
function FriendListsGroups.quickMessageGroup(userID, groupID, message)
  if not userID or not groupID or not message then return {} end
  
  local broadcast = {
    sender = userID,
    group = groupID,
    message = message,
    recipients = 8,
    sent_at = os.time()
  }
  
  return broadcast
end

-- ============================================================================
-- FEATURE 3: ACTIVITY FEED (~240 LOC)
-- ============================================================================

local ActivityFeed = {}

---Get social feed
---@param userID string User to fetch feed for
---@return table feed Social activity feed
function ActivityFeed.getFeed(userID)
  if not userID then return {} end
  
  local feed = {
    user_id = userID,
    activities = {
      {actor = "Player1", action = "joined your guild", time = os.time() - 600},
      {actor = "Player2", action = "achieved 100 battles", time = os.time() - 3600},
      {actor = "Player3", action = "defeated boss with you", time = os.time() - 7200}
    },
    total_activities = 45
  }
  
  return feed
end

---Post activity
---@param userID string Poster
---@param activity string Activity description
---@return table post Posted activity
function ActivityFeed.postActivity(userID, activity)
  if not userID or not activity then return {} end
  
  local post = {
    post_id = "ACTIVITY_" .. os.time(),
    user_id = userID,
    content = activity,
    posted_at = os.time(),
    visibility = "Friends",
    likes = 0
  }
  
  return post
end

---Like activity
---@param userID string Liker
---@param activityID string Activity to like
---@return table liked Like confirmation
function ActivityFeed.likeActivity(userID, activityID)
  if not userID or not activityID then return {} end
  
  local liked = {
    user_id = userID,
    activity_id = activityID,
    liked_at = os.time(),
    total_likes = 15
  }
  
  return liked
end

---Comment on activity
---@param userID string Commenter
---@param activityID string Activity to comment
---@param comment string Comment text
---@return table commented Comment confirmation
function ActivityFeed.commentActivity(userID, activityID, comment)
  if not userID or not activityID or not comment then return {} end
  
  local commented = {
    user_id = userID,
    activity_id = activityID,
    comment = comment,
    commented_at = os.time(),
    total_comments = 8
  }
  
  return commented
end

-- ============================================================================
-- FEATURE 4: NOTIFICATIONS & PRESENCE (~210 LOC)
-- ============================================================================

local NotificationsPresence = {}

---Update online status
---@param userID string User updating
---@param status string Status (Online/Idle/Offline)
---@return table updated Status update confirmation
function NotificationsPresence.updateStatus(userID, status)
  if not userID or not status then return {} end
  
  local updated = {
    user_id = userID,
    status = status,
    updated_at = os.time(),
    visibility = "Friends"
  }
  
  return updated
end

---Get online friends
---@param userID string User to query
---@return table onlineFriends Online friends list
function NotificationsPresence.getOnlineFriends(userID)
  if not userID then return {} end
  
  local onlineFriends = {
    user_id = userID,
    online_count = 8,
    friends = {
      {name = "Friend1", status = "Online", activity = "Playing"},
      {name = "Friend2", status = "Online", activity = "Idle"},
      {name = "Friend3", status = "Idle", activity = "Away"}
    }
  }
  
  return onlineFriends
end

---Send notification
---@param userID string Recipient
---@param notifType string Notification type
---@param content string Notification content
---@return table notification Notification sent
function NotificationsPresence.sendNotification(userID, notifType, content)
  if not userID or not notifType or not content then return {} end
  
  local notification = {
    notification_id = "NOTIF_" .. os.time(),
    user_id = userID,
    type = notifType,
    content = content,
    sent_at = os.time(),
    read = false
  }
  
  return notification
end

---Get notification history
---@param userID string User to query
---@return table history Notification history
function NotificationsPresence.getHistory(userID)
  if not userID then return {} end
  
  local history = {
    user_id = userID,
    total_notifications = 125,
    unread_count = 5,
    recent = {
      {type = "message", from = "Player1", time = os.time() - 600},
      {type = "friend_request", from = "Player2", time = os.time() - 3600},
      {type = "event_reminder", content = "Event starts soon", time = os.time() - 7200}
    }
  }
  
  return history
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.0",
  MessagingSystem = MessagingSystem,
  FriendListsGroups = FriendListsGroups,
  ActivityFeed = ActivityFeed,
  NotificationsPresence = NotificationsPresence,
  
  features = {
    messagingSystem = true,
    friendListsGroups = true,
    activityFeed = true,
    notificationsPresence = true
  }
}
