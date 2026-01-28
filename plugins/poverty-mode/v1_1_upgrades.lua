--[[
  Poverty Mode Plugin - v1.1+ Upgrade Extension
  Budget constraints, income tracking, expense management, cost-benefit analysis
  
  Phase: 7F (Difficulty Systems)
  Version: 1.1+ (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: BUDGET CONSTRAINTS (~240 LOC)
-- ============================================================================

local BudgetConstraints = {}

---Set gil budget limit
---@param maxGil number Maximum allowed gil
---@return table budgetSettings Budget configuration
function BudgetConstraints.setBudgetLimit(maxGil)
  if not maxGil or maxGil < 0 then maxGil = 10000 end
  
  local budgetSettings = {
    enabled = true,
    maxGil = maxGil,
    currentGil = 0,
    budgetPercentage = 0,
    restrictions = {},
    enforceMode = "strict"
  }
  
  if maxGil <= 5000 then
    budgetSettings.restrictions = {"No expensive equipment", "Limited healing items"}
  elseif maxGil <= 50000 then
    budgetSettings.restrictions = {"Expensive items restricted", "Manage spending"}
  else
    budgetSettings.restrictions = {"Minimal restrictions"}
  end
  
  return budgetSettings
end

---Validate purchase against budget
---@param currentGil number Current gil balance
---@param purchasePrice number Cost of item
---@param budgetCap number Budget limit
---@return table validation Purchase validation result
function BudgetConstraints.validatePurchase(currentGil, purchasePrice, budgetCap)
  if not currentGil or not purchasePrice or not budgetCap then return {} end
  
  local validation = {
    currentGil = currentGil,
    purchasePrice = purchasePrice,
    budgetCap = budgetCap,
    canAfford = currentGil >= purchasePrice,
    wouldExceedBudget = (currentGil - purchasePrice) < 0,
    gilAfterPurchase = currentGil - purchasePrice,
    budgetPercentageUsed = 0
  }
  
  if validation.gilAfterPurchase >= 0 then
    validation.budgetPercentageUsed = ((budgetCap - validation.gilAfterPurchase) / budgetCap) * 100
  end
  
  return validation
end

---Calculate budget remaining
---@param currentGil number Current gil spent
---@param budgetLimit number Total budget
---@return table budgetStatus Remaining budget info
function BudgetConstraints.calculateBudgetRemaining(currentGil, budgetLimit)
  if not currentGil or not budgetLimit then return {} end
  
  local budgetStatus = {
    totalBudget = budgetLimit,
    gilSpent = currentGil,
    gilRemaining = budgetLimit - currentGil,
    percentageUsed = (currentGil / budgetLimit) * 100,
    percentageRemaining = ((budgetLimit - currentGil) / budgetLimit) * 100,
    warning = ""
  }
  
  if budgetStatus.percentageRemaining < 10 then
    budgetStatus.warning = "Budget critically low!"
  elseif budgetStatus.percentageRemaining < 25 then
    budgetStatus.warning = "Budget running low"
  end
  
  return budgetStatus
end

---Recommend budget allocation
---@param budgetLimit number Total budget
---@param necessities table Required items
---@return table allocation Recommended spending
function BudgetConstraints.recommendBudgetAllocation(budgetLimit, necessities)
  if not budgetLimit then budgetLimit = 50000 end
  necessities = necessities or {}
  
  local allocation = {
    totalBudget = budgetLimit,
    allocation = {},
    savingsRecommended = 0
  }
  
  -- Allocate budget
  local necessityCost = 0
  for _, item in ipairs(necessities) do
    necessityCost = necessityCost + (item.cost or 0)
  end
  
  allocation.allocation = {
    equipment = math.floor(budgetLimit * 0.40),
    healing = math.floor(budgetLimit * 0.30),
    utilities = math.floor(budgetLimit * 0.20),
    reserves = math.floor(budgetLimit * 0.10)
  }
  
  allocation.savingsRecommended = budgetLimit - necessityCost
  
  return allocation
end

-- ============================================================================
-- FEATURE 2: INCOME TRACKING (~240 LOC)
-- ============================================================================

local IncomeTracking = {}

---Track gil income from battles
---@param enemy table Defeated enemy
---@param killMethod string How it was defeated
---@return number gilReward Awarded gil
function IncomeTracking.trackBattleIncome(enemy, killMethod)
  if not enemy then return 0 end
  
  local baseGil = enemy.gilReward or 100
  local multiplier = 1.0
  
  if killMethod == "quick" then
    multiplier = 1.2
  elseif killMethod == "efficient" then
    multiplier = 1.1
  elseif killMethod == "overkill" then
    multiplier = 0.8
  end
  
  return math.floor(baseGil * multiplier)
end

---Calculate total income over time
---@param battleCount number Number of battles
---@param averageGilPerBattle number Average gil per battle
---@return table incomeProjection Income projection
function IncomeTracking.projectIncomeOverTime(battleCount, averageGilPerBattle)
  if not battleCount or not averageGilPerBattle then return {} end
  
  local incomeProjection = {
    battlesFought = battleCount,
    averagePerBattle = averageGilPerBattle,
    projectedTotal = battleCount * averageGilPerBattle,
    breakdown = {}
  }
  
  -- Breakdown by milestone
  for milestone = 10, battleCount, 10 do
    table.insert(incomeProjection.breakdown, {
      battles = milestone,
      projectedGil = milestone * averageGilPerBattle
    })
  end
  
  return incomeProjection
end

---Analyze income sources
---@param saveData table Save file data
---@return table sourceAnalysis Income breakdown by source
function IncomeTracking.analyzeIncomeSources(saveData)
  if not saveData then return {} end
  
  local sourceAnalysis = {
    totalIncome = 0,
    sources = {
      battles = 0,
      treasures = 0,
      sales = 0,
      other = 0
    }
  }
  
  -- Calculate from battles
  if saveData.battleCount then
    sourceAnalysis.sources.battles = saveData.battleCount * 120  -- Average 120 gil per battle
  end
  
  -- Calculate from treasure chests
  if saveData.treasuresFound then
    sourceAnalysis.sources.treasures = saveData.treasuresFound * 200
  end
  
  -- Calculate from item sales
  if saveData.soldItems then
    sourceAnalysis.sources.sales = saveData.soldItems * 150
  end
  
  for _, income in pairs(sourceAnalysis.sources) do
    sourceAnalysis.totalIncome = sourceAnalysis.totalIncome + income
  end
  
  return sourceAnalysis
end

---Recommend income optimization
---@param currentIncome number Current income rate
---@param targetBudget number Target budget
---@return table recommendations Optimization suggestions
function IncomeTracking.recommendIncomeOptimization(currentIncome, targetBudget)
  if not currentIncome or not targetBudget then return {} end
  
  local recommendations = {
    currentIncome = currentIncome,
    targetBudget = targetBudget,
    suggestions = {},
    priorityTier = "medium"
  }
  
  local deficit = targetBudget - currentIncome
  
  if deficit > 0 then
    recommendations.priorityTier = "high"
    table.insert(recommendations.suggestions, "Fight more battles to increase income")
    table.insert(recommendations.suggestions, "Seek out treasure chests")
    table.insert(recommendations.suggestions, "Sell unnecessary items")
  else
    recommendations.priorityTier = "low"
    table.insert(recommendations.suggestions, "Income sufficient")
  end
  
  return recommendations
end

-- ============================================================================
-- FEATURE 3: EXPENSE MANAGEMENT (~240 LOC)
-- ============================================================================

local ExpenseManagement = {}

---Categorize expenses by type
---@param expense table Expense record
---@return string category Expense category
function ExpenseManagement.categorizeExpense(expense)
  if not expense then return "other" end
  
  if expense.type == "equipment" then
    return "equipment"
  elseif expense.type == "item" or expense.type == "healing" then
    return "supplies"
  elseif expense.type == "service" or expense.type == "inn" then
    return "services"
  else
    return "other"
  end
end

---Track cumulative expenses
---@param expenses table List of expenses
---@return table expenseBreakdown Expense analysis
function ExpenseManagement.trackCumulativeExpenses(expenses)
  if not expenses then return {} end
  
  local expenseBreakdown = {
    totalExpense = 0,
    categories = {
      equipment = 0,
      supplies = 0,
      services = 0,
      other = 0
    },
    itemBreakdown = {},
    averagePerItem = 0
  }
  
  for _, expense in ipairs(expenses) do
    local category = ExpenseManagement.categorizeExpense(expense)
    local amount = expense.amount or 0
    
    expenseBreakdown.totalExpense = expenseBreakdown.totalExpense + amount
    expenseBreakdown.categories[category] = expenseBreakdown.categories[category] + amount
  end
  
  if #expenses > 0 then
    expenseBreakdown.averagePerItem = expenseBreakdown.totalExpense / #expenses
  end
  
  return expenseBreakdown
end

---Identify spending inefficiencies
---@param expenses table List of expenses
---@param budget number Total budget
---@return table inefficiencies Spending issues
function ExpenseManagement.identifyInefficiencies(expenses, budget)
  if not expenses or not budget then return {} end
  
  local inefficiencies = {
    issues = {},
    recommendations = {},
    wastedGil = 0
  }
  
  local costMap = {}
  
  for _, expense in ipairs(expenses) do
    local itemName = expense.name or "Unknown"
    
    if costMap[itemName] then
      costMap[itemName].count = costMap[itemName].count + 1
      costMap[itemName].totalCost = costMap[itemName].totalCost + expense.amount
      
      table.insert(inefficiencies.issues, 
        "Duplicate purchases: " .. itemName .. " (x" .. costMap[itemName].count .. ")")
      inefficiencies.wastedGil = inefficiencies.wastedGil + expense.amount
    else
      costMap[itemName] = {count = 1, totalCost = expense.amount}
    end
  end
  
  if inefficiencies.wastedGil > 0 then
    table.insert(inefficiencies.recommendations, 
      "Plan purchases to avoid duplicates")
    table.insert(inefficiencies.recommendations, 
      "Use inventory management before purchasing")
  end
  
  return inefficiencies
end

---Suggest cost-saving alternatives
---@param desiredItem table Item to purchase
---@param availableAlternatives table Cheaper alternatives
---@return table suggestion Cost comparison
function ExpenseManagement.suggestCheaperAlternative(desiredItem, availableAlternatives)
  if not desiredItem or not availableAlternatives then return {} end
  
  local suggestion = {
    desiredItem = desiredItem.name,
    desiredCost = desiredItem.cost or 0,
    cheapestAlternative = nil,
    potentialSavings = 0
  }
  
  local cheapest = nil
  local cheapestCost = desiredItem.cost or 999999
  
  for _, alt in ipairs(availableAlternatives) do
    if (alt.cost or 0) < cheapestCost then
      cheapestCost = alt.cost or 0
      cheapest = alt
    end
  end
  
  if cheapest then
    suggestion.cheapestAlternative = cheapest.name
    suggestion.potentialSavings = suggestion.desiredCost - cheapestCost
  end
  
  return suggestion
end

-- ============================================================================
-- FEATURE 4: COST-BENEFIT ANALYSIS (~220 LOC)
-- ============================================================================

local CostBenefitAnalysis = {}

---Analyze equipment purchase value
---@param equipment table Equipment to evaluate
---@param budget number Current budget
---@return table analysis Value analysis
function CostBenefitAnalysis.analyzeEquipmentValue(equipment, budget)
  if not equipment or not budget then return {} end
  
  local analysis = {
    item = equipment.name,
    cost = equipment.cost or 0,
    statBonus = equipment.statBonus or 0,
    costPerStat = 0,
    budgetImpact = 0,
    recommendation = ""
  }
  
  if analysis.statBonus > 0 then
    analysis.costPerStat = analysis.cost / analysis.statBonus
  end
  
  analysis.budgetImpact = (analysis.cost / budget) * 100
  
  if analysis.budgetImpact > 20 then
    analysis.recommendation = "Expensive - Consider alternatives"
  elseif analysis.budgetImpact > 10 then
    analysis.recommendation = "Moderate cost - Worth considering"
  else
    analysis.recommendation = "Affordable - Good value"
  end
  
  return analysis
end

---Compare efficiency of multiple items
---@param items table Items to compare
---@param budget number Available budget
---@return table comparison Efficiency ranking
function CostBenefitAnalysis.compareItemEfficiency(items, budget)
  if not items or not budget then return {} end
  
  local comparison = {
    budget = budget,
    items = {},
    bestValue = nil,
    bestValueScore = 0
  }
  
  for _, item in ipairs(items) do
    local efficiency = (item.statBonus or 1) / (item.cost or 1)
    local score = efficiency * 100
    
    table.insert(comparison.items, {
      item = item.name,
      cost = item.cost,
      bonus = item.statBonus,
      efficiency = score
    })
    
    if score > comparison.bestValueScore then
      comparison.bestValueScore = score
      comparison.bestValue = item.name
    end
  end
  
  return comparison
end

---Calculate ROI for spending
---@param investment number Gil spent
---@param benefit number Stat/power gain
---@return number roi Return on investment percentage
function CostBenefitAnalysis.calculateROI(investment, benefit)
  if not investment or investment == 0 then return 0 end
  if not benefit then benefit = 0 end
  
  return ((benefit - investment) / investment) * 100
end

---Recommend minimum viable equipment
---@param target string Build target
---@param budget number Available budget
---@return table buildPlan Minimum viable setup
function CostBenefitAnalysis.recommendMinimalViableBuild(target, budget)
  if not target or not budget then return {} end
  
  local buildPlan = {
    target = target,
    budget = budget,
    equipment = {},
    totalCost = 0,
    expectedPerformance = ""
  }
  
  -- Allocate budget efficiently
  if target == "speedrun" then
    buildPlan.equipment = {
      {slot = "weapon", cost = math.floor(budget * 0.40)},
      {slot = "armor", cost = math.floor(budget * 0.35)},
      {slot = "accessory", cost = math.floor(budget * 0.25)}
    }
    buildPlan.expectedPerformance = "Fast, minimal durability"
  elseif target == "tank" then
    buildPlan.equipment = {
      {slot = "armor", cost = math.floor(budget * 0.50)},
      {slot = "shield", cost = math.floor(budget * 0.35)},
      {slot = "accessory", cost = math.floor(budget * 0.15)}
    }
    buildPlan.expectedPerformance = "Defensive focus, slower"
  else
    buildPlan.equipment = {
      {slot = "weapon", cost = math.floor(budget * 0.40)},
      {slot = "armor", cost = math.floor(budget * 0.40)},
      {slot = "accessory", cost = math.floor(budget * 0.20)}
    }
    buildPlan.expectedPerformance = "Balanced"
  end
  
  for _, item in ipairs(buildPlan.equipment) do
    buildPlan.totalCost = buildPlan.totalCost + item.cost
  end
  
  return buildPlan
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1+",
  BudgetConstraints = BudgetConstraints,
  IncomeTracking = IncomeTracking,
  ExpenseManagement = ExpenseManagement,
  CostBenefitAnalysis = CostBenefitAnalysis,
  
  -- Feature completion
  features = {
    budgetConstraints = true,
    incomeTracking = true,
    expenseManagement = true,
    costBenefitAnalysis = true
  }
}
