# FF6 Plugin Upgrade Plan - v1.1 Release Phase

**Plan Created:** January 16, 2026  
**Status:** PLANNING  
**Target:** Enhance 7 plugins from Phase 5-6 with v1.1 upgrades  
**Priority Focus:** Advanced features and quality-of-life improvements

---

## Executive Summary

This document outlines comprehensive upgrades for 7 existing plugins to v1.1, addressing known limitations, adding advanced features, and improving user experience. These upgrades will build on the successful 46-plugin ecosystem, taking proven plugins to the next level.

**Upgrade Scope:**
- 7 plugins receiving v1.1 updates
- ~50+ new functions across all plugins
- ~3,000+ additional lines of code
- ~25,000+ words of new documentation
- Estimated delivery: Phases 7A-7C (3 batches of 2-3 plugins each)

---

## Upgrade Philosophy

These v1.1 releases focus on:
- **Advanced Capabilities** - Features cut from v1.0 due to time constraints
- **User Empowerment** - Granular control and customization
- **Quality of Life** - Automation and convenience features
- **Integration** - Better cross-plugin compatibility
- **Professional Polish** - Industry-standard tool features

All upgrades maintain **backward compatibility** with v1.0 while adding powerful new systems.

---

## Upgrade Roadmap

### Phase 7A: Analysis & Tracking Tools (Batch 1 - 2 plugins)

**Timeline:** After Phase 6 completion  
**Target Plugins:** 
1. **Randomizer Assistant v1.1** - Advanced seed analysis
2. **Speedrun Timer v1.1** - Professional timing integration

---

## Plugin Upgrades (Detailed)

### 1. RANDOMIZER ASSISTANT v1.1

**Current State (v1.0):**
- Basic spoiler log parsing (JSON/CSV)
- Manual location tracking
- Simplified logic checking
- Beatability verification (basic)
- Static analysis only

**Known Limitations (v1.0):**
- Can't auto-detect obtained items from save file
- Simplified logic checking misses some softlocks
- No real-time accessibility calculation
- No visual map support
- No community features

**v1.1 Upgrade Plan:**

#### New Features (8 major additions)

**1. Auto-Detection System** (~200 LOC)
- `autoDetectObtainedItems()` - Scan save file for item acquisition
- `compareSeedToSaveState()` - Match save progress against spoiler log
- `trackItemProgression()` - Automatic location checking
- `syncSaveWithSeed()` - Sync progress state
- Returns: Completion percentage, next accessible locations

**2. Advanced Logic Solver** (~400 LOC)
- `analyzeProgressionChains()` - Build dependency graph
- `detectSoftlocks()` - Advanced detection algorithm
- `findAlternateRoutes()` - Multiple path solutions
- `calculateAccessibilityTree()` - All reachable states
- Returns: Complete walkability analysis

**3. Real-Time Accessibility** (~300 LOC)
- `calculateAccessible()` - Dynamic reachability calc
- `getPrevented()` - What locations are blocked
- `suggestNextLocation()` - Recommended progression
- `estimateCompletion()` - Percent to goal
- Returns: Current game state analysis

**4. Visual Map System** (~350 LOC)
- `generateLocationMap()` - Create visual map layout
- `markObtainedLocations()` - Show collected items
- `highlightAccessible()` - Show reachable locations
- `showProgressPath()` - Display taken route
- `exportMapImage()` - Save map graphic
- Returns: Visual representation

**5. Hint System** (~250 LOC)
- `provideHint()` - Context-aware suggestions
- `getSpoilerLevel()` - Adjust hint detail (no/minor/major spoilers)
- `getLocationHints()` - Tips for specific locations
- `trackHintsUsed()` - Count hints for scoring
- Returns: Player-friendly hints

**6. Community Features** (~300 LOC)
- `saveSeedProfile()` - Upload seed info
- `searchCommunitySeeds()` - Find similar seeds
- `viewCommunityStats()` - Global statistics
- `shareProgressionSolution()` - Post solution
- `downloadCommunityGuide()` - Get community guides
- Returns: Cloud-integrated data

**7. Seed Generation** (~250 LOC)
- `generateCustomSeed()` - Create new randomization
- `tweakSeedDifficulty()` - Adjust complexity
- `applySeedPreset()` - Use template
- `validateSeedLogic()` - Check for softlocks
- Returns: New seed with validation

**8. Real-Time Progression Analytics** (~200 LOC)
- `trackItemCollection()` - Monitor items obtained
- `calculateOptimalPath()` - Best route to completion
- `estimateTimeToCompletion()` - Speedrun projection
- `analyzeTimingBottlenecks()` - Where time lost
- Returns: Advanced statistics

#### Implementation Details

**New Functions:** 30+
**Estimated LOC:** ~1,850  
**Documentation:** ~7,000 words
**Testing:** Logic validation suite

**API Dependencies:**
- `save.getInventory()` - Read collected items
- `save.getProgress()` - Read completion state
- `map.generateLayout()` - Create visual maps
- `cloud.uploadData()` - Community sharing

**User Workflows:**

```lua
-- Auto-detect workflow
local seed_file = "beyond_chaos_seed.json"
importSeedFromFile(seed_file)
autoDetectObtainedItems()  -- Scans current save

-- Shows: "You've collected 45% of items"
-- Suggests: "Try Location A next (has 3 accessible items)"

-- Real-time accessibility
local accessible = calculateAccessible()
-- Returns: {locations: [...], items: [...], bosses: [...]}

-- Community sharing
local profile = saveSeedProfile()
-- Uploads to community, gets difficulty rating, etc.
```

**Backward Compatibility:** ✅ All v1.0 functions retained

---

### 2. SPEEDRUN TIMER v1.1

**Current State (v1.0):**
- Basic lap/segment timing
- Single split per category
- Manual time entry
- Static comparison
- File export only

**Known Limitations (v1.0):**
- Limited split management
- No personal record comparison
- No streaming integration
- Basic pace calculation
- No category presets

**v1.1 Upgrade Plan:**

#### New Features (7 major additions)

**1. Advanced Split Management** (~300 LOC)
- `createSplitProfile()` - Save split layout
- `saveSplitPreset()` - Store multiple presets
- `loadSplitPreset()` - Load saved layout
- `compareSplits()` - Split-by-split analysis
- `manageSplitLibrary()` - Organize splits
- `importExternalSplits()` - Load from files
- Returns: Complete split infrastructure

**2. Personal Record System** (~250 LOC)
- `setPBComparison()` - Track personal best
- `compareToRecord()` - Current vs. PB
- `calculateDeltaTime()` - Ahead/behind display
- `predictFinalTime()` - Projection based on current pace
- `archiveRecord()` - Save historical runs
- Returns: Record tracking

**3. Pace Calculator v2** (~280 LOC)
- `calculatePaceRealistic()` - Expected pace
- `detectPaceChanges()` - When pace shifted
- `projectedCompletion()` - Finish time estimate
- `paceBreakdown()` - Per-segment analysis
- `compareSegmentPace()` - Split comparison
- Returns: Advanced pacing metrics

**4. Category Presets** (~200 LOC)
- `createCategory()` - Define speedrun category
- `addCategoryPreset()` - Pre-configured setups
- Standard categories:
  - Any% (current support)
  - 100% (all items/espers)
  - Any% Glitchless
  - Low% (minimize items)
  - No-Save Run
  - Blind Race
- `loadCategoryPreset()` - Apply preset
- Returns: Category management

**5. Streaming Integration** (~350 LOC)
- `enableStreamMode()` - Optimize for streaming
- `streamTimer()` - OBS/Twitch integration
- `createStreamOverlay()` - Custom timer display
- `autoUpdateTwitch()` - Real-time Twitch update
- `generateStreamGraphic()` - Current pace graphic
- `trackChatPredictions()` - Chat engagement
- Returns: Full streaming support

**6. Race Mode** (~300 LOC)
- `startRaceMode()` - Multiplayer timing
- `trackMultiplePlayers()` - Side-by-side times
- `calculateLeadStatus()` - Who's ahead
- `raceSynchronization()` - Sync between runners
- `generateRaceStats()` - Head-to-head analysis
- `recordRaceVideo()` - Save race VOD metadata
- Returns: Competitive timing

**7. Time Tracking Analytics** (~280 LOC)
- `trackRunConsistency()` - Variance analysis
- `identifyTimeoutPoints()` - Where consistency breaks
- `generatePerformanceGraph()` - Visual trend
- `benchmarkAgainstCommunity()` - World record comparison
- `analyzeEfficiency()` - Route optimization
- Returns: Performance metrics

#### Implementation Details

**New Functions:** 35+
**Estimated LOC:** ~1,960
**Documentation:** ~7,500 words
**Testing:** Timing precision verification

**API Dependencies:**
- `streaming.connectOBS()` - OBS integration
- `twitch.sendUpdate()` - Twitch connection
- `system.getTime()` - Precise timing
- `file.importSplits()` - Load external splits

**User Workflows:**

```lua
-- Stream setup workflow
enableStreamMode()
loadCategoryPreset("Any%")
createStreamOverlay("horizontal")  -- Timer display style
connectOBS("localhost:4444")
generateStreamGraphic()

-- Race workflow
startRaceMode()
trackMultiplePlayers({player1, player2})
-- Shows: "Player 1 ahead by 0:32"

-- Analytics workflow
local consistency = trackRunConsistency()
benchmarkAgainstCommunity()
-- Shows: "You're 15% faster than last week"
```

**Backward Compatibility:** ✅ All v1.0 functions retained

---

### 3. CHALLENGE MODE VALIDATOR v1.1

**Current State (v1.0):**
- Basic rule checking
- Snapshot-based proof
- Static validation
- Limited tracking

**Known Limitations (v1.0):**
- Real-time tracking incomplete
- Violations logged at check-time only
- No continuous monitoring
- Limited proof system

**v1.1 Upgrade Plan:**

#### New Features (6 major additions)

**1. Real-Time Event Tracking** (~280 LOC)
- `enableEventHooks()` - Start monitoring
- `trackBattleWins()` - Battle counting
- `trackGilSpent()` - Money tracking
- `trackItemUsage()` - Item consumption
- `trackMapMovement()` - Location tracking
- `trackLevelUp()` - Character progression
- Returns: Live event stream

**2. Continuous Violation Monitoring** (~300 LOC)
- `startContinuousMonitoring()` - Activate monitoring
- `detailViolationEvent()` - Log violation details
- `preventViolation()` - Block rule breaks (sandboxing)
- `warnViolation()` - Alert player
- `logViolationChain()` - Track related violations
- Returns: Complete violation history

**3. Multi-Save Challenge Progression** (~250 LOC)
- `trackProgressionAcrossSaves()` - Chain saves
- `mergeRunSegments()` - Combine separate sessions
- `validateContinuity()` - Check no cheating between saves
- `generateProgressReport()` - Segment analysis
- Returns: Multi-save validation

**4. Advanced Proof System** (~280 LOC)
- `createProofSnapshot()` - Snapshot at intervals
- `generateProofChain()` - Linked proofs
- `createVideoProof()` - Recording metadata
- `generateScreenshotProof()` - Image evidence
- `createRepayProof()` - Replay data
- Returns: Professional proof system

**5. Community Leaderboard** (~320 LOC)
- `submitChallengeResult()` - Post to leaderboard
- `viewChallengeLeaderboard()` - Global rankings
- `compareToLeaderboard()` - Your rank
- `verifyLeaderboardEntry()` - Anti-cheat check
- `generateLeaderboardProof()` - Verify submission
- Returns: Competitive leaderboard

**6. Challenge AI Difficulty** (~250 LOC)
- `rateChallengeComplexity()` - Difficulty rating
- `compareToLeaderboard()` - How hard vs. others
- `suggestAlternateRules()` - Similar difficulty
- `generateChallengeRating()` - Community rating
- Returns: Complexity metrics

#### Implementation Details

**New Functions:** 25+
**Estimated LOC:** ~1,680
**Documentation:** ~6,500 words
**Testing:** Event hook verification

**API Dependencies:**
- `events.hookBattle()` - Battle events
- `events.hookGilTransaction()` - Money events
- `events.hookItemUsage()` - Item events
- `cloud.submitLeaderboardEntry()` - Leaderboard

---

### 4. INSTANT MASTERY SYSTEM v1.1

**Current State (v1.0):**
- All-or-nothing mastery
- Basic presets
- Limited control

**v1.1 Upgrade Plan:**

#### New Features (4 major additions)

**1. Granular Stat Control** (~200 LOC)
- `selectiveStatMastery()` - Choose specific stats
- `customStatPreset()` - User-defined setup
- `buildStatTemplate()` - Save configuration
- `applyPartialMastery()` - Selective unlock
- Returns: Fine-grained control

**2. Preset Templates Library** (~180 LOC)
- `createPresetTemplate()` - Save builds
- `manageTemplateLibrary()` - Organize presets
- Pre-built templates:
  - Speedrunner Build
  - Solo Challenge Build
  - Pacifist Build
  - All Equipment Build
  - Level-Capped Build
- `applyTemplatePreset()` - Load preset
- Returns: Template system

**3. Stat Limit Calculator** (~150 LOC)
- `calculateStatLimits()` - Game engine limits
- `detectHardCaps()` - Hardware constraints
- `optimizeWithinLimits()` - Max safe stats
- `warnAboutLimits()` - Alert on issues
- Returns: Safe configuration

**4. Build Analyzer** (~200 LOC)
- `analyzeBuildSynergies()` - Stat combinations
- `validateBuildLogic()` - Check dependencies
- `suggestOptimalBuild()` - Recommendations
- `detectBuildConflicts()` - Problem detection
- Returns: Build analysis

#### Implementation Details

**New Functions:** 15+
**Estimated LOC:** ~730
**Documentation:** ~4,500 words

---

### 5. CHARACTER ABILITY SWAP v1.1

**Current State (v1.0):**
- Basic ability swapping
- No preview system
- Limited conflict checking

**v1.1 Upgrade Plan:**

#### New Features (3 major additions)

**1. Full Ability Preview** (~250 LOC)
- `previewAbilitySwap()` - Show results before apply
- `compareAbilities()` - Side-by-side analysis
- `simulateBattle()` - Test in combat scenario
- `showAbilitySynergies()` - Stat combinations
- Returns: Detailed preview

**2. Ability Synergy Analyzer** (~220 LOC)
- `analyzeAbilitySynergy()` - Compatibility scoring
- `suggestOptimalAbilities()` - Recommendations
- `detectSynergyBonus()` - Combination bonuses
- `rateAbilitySetup()` - Overall effectiveness
- Returns: Synergy metrics

**3. Swap Conflict Detection** (~180 LOC)
- `detectConflicts()` - Incompatible combinations
- `warnAboutIssues()` - Alert player
- `suggestAlternatives()` - Workarounds
- `preventInvalidSwaps()` - Blocking system
- Returns: Conflict analysis

#### Implementation Details

**New Functions:** 12+
**Estimated LOC:** ~650
**Documentation:** ~4,000 words

---

### 6. STORYLINE EDITOR v1.1

**Current State (v1.0):**
- Basic dialogue/event editing
- Story branch control
- No visualization

**v1.1 Upgrade Plan:**

#### New Features (3 major additions)

**1. Dialogue Preview System** (~200 LOC)
- `previewDialogueInGame()` - Context display
- `simulateConversation()` - Show full scene
- `checkTextboxFit()` - Verify rendering
- `previewWithCharacterSprites()` - Visual context
- Returns: Dialogue preview

**2. Story Arc Visualization** (~280 LOC)
- `visualizeStoryArc()` - Graph/tree display
- `showBranchDependencies()` - Relationship mapping
- `highlightProgressionPath()` - Current route
- `exportStoryMap()` - Save visualization
- Returns: Visual story structure

**3. Advanced Event System** (~250 LOC)
- `createConditionalEvents()` - If/then logic
- `chainEventSequences()` - Complex chains
- `createEventBranches()` - Conditional paths
- `validateEventLogic()` - Check for loops
- Returns: Advanced event control

#### Implementation Details

**New Functions:** 13+
**Estimated LOC:** ~730
**Documentation:** ~4,500 words

---

### 7. SPRITE SWAPPER v1.1

**Current State (v1.0):**
- Basic sprite swaps
- Custom import
- Simple presets

**v1.1 Upgrade Plan:**

#### New Features (4 major additions)

**1. Sprite Preview Gallery** (~220 LOC)
- `createPreviewGallery()` - Visual browser
- `showSpriteComparison()` - Side-by-side view
- `previewCustomSprites()` - Show imports
- `tagAndOrganizeSprites()` - Library management
- `searchSpriteLibrary()` - Find sprites
- Returns: Visual sprite browser

**2. Animation Swap System** (~280 LOC)
- `swapCharacterAnimations()` - Change animations
- `combineSpritesAndAnimations()` - Mix sprite + anim
- `previewAnimation()` - See it moving
- `validateAnimationCompatibility()` - Check fit
- Returns: Animation control

**3. Color Palette System** (~250 LOC)
- `recolorSprite()` - Palette shift
- `createCustomPalette()` - New colors
- `applyPalettePreset()` - Template colors
- `preserveBackgroundColors()` - Selective recolor
- Returns: Color customization

**4. Sprite Collision Adjustment** (~200 LOC)
- `adjustCollisionBox()` - Hit box sizing
- `testCollisionInGame()` - Verify in combat
- `scaleCollisionProportional()` - Auto-scale
- `validateCollisionLogic()` - Check safety
- Returns: Collision control

#### Implementation Details

**New Functions:** 16+
**Estimated LOC:** ~950
**Documentation:** ~5,000 words

---

## Implementation Schedule

### Phase 7A: Advanced Analysis (2 plugins)
**Plugins:** Randomizer Assistant v1.1, Speedrun Timer v1.1
**Estimated LOC:** ~3,810
**Estimated Words:** ~14,500
**Timeline:** Weeks 1-2
**Files:** 6 (2 plugins × 3 files each - CHANGELOG only, others updated)

### Phase 7B: Validation & Verification (2 plugins)
**Plugins:** Challenge Mode Validator v1.1, Instant Mastery System v1.1
**Estimated LOC:** ~2,410
**Estimated Words:** ~11,000
**Timeline:** Weeks 3-4
**Files:** 4

### Phase 7C: Creative Tools (3 plugins)
**Plugins:** Character Ability Swap v1.1, Storyline Editor v1.1, Sprite Swapper v1.1
**Estimated LOC:** ~2,330
**Estimated Words:** ~13,500
**Timeline:** Weeks 5-6
**Files:** 6

---

## Overall Upgrade Statistics

**Total Upgrades:** 7 plugins
**Total New Functions:** 130+
**Total New LOC:** ~8,550
**Total Documentation:** ~54,500 words
**Total Files Modified:** 16+ (CHANGELOGs + plugin updates)

**New Project Totals (After v1.1):**
- 46 base plugins (v1.0) + 7 upgraded to v1.1
- 53 actively maintained plugins
- 26,465 total LOC (~8,550 new)
- 290,800+ total words (~54,500 new)
- 180+ files

---

## Quality Assurance

Each v1.1 upgrade will include:
- ✅ Unit testing for new functions
- ✅ Integration testing with existing v1.0 functions
- ✅ Backward compatibility verification
- ✅ Documentation review
- ✅ User workflow testing
- ✅ Edge case handling

---

## Backward Compatibility

**All v1.1 upgrades maintain 100% backward compatibility with v1.0:**
- All v1.0 functions remain unchanged
- New functions are additive only
- No breaking changes to existing APIs
- Users can upgrade safely without relearning

---

## Distribution Plan

**v1.1 Release Package:**
1. Release Notes (major features for each plugin)
2. Updated Changelogs (v1.0 → v1.1)
3. Enhanced Documentation (all README updates)
4. Migration Guide (optional new features explanation)
5. Upgrade Instructions (seamless update path)

---

## Success Criteria

✅ All 130+ functions implemented
✅ All functions documented with examples
✅ Backward compatibility 100%
✅ Testing coverage > 95%
✅ User workflows validated
✅ Cross-plugin integration working
✅ Master plan updated

---

## Next Steps

1. **User Approval** - Confirm upgrade priority
2. **Schedule** - Allocate time for Phase 7A-7C
3. **Development** - Build v1.1 upgrades
4. **Testing** - Verify quality
5. **Integration** - Update master plan
6. **Release** - Deliver to project

---

**Plan Status:** READY FOR APPROVAL

**Recommendation:** Begin with Phase 7A (Randomizer Assistant + Speedrun Timer) as they provide maximum value to users and showcase advanced features well.

---

*Upgrade Plan Created: January 16, 2026*  
*Target Completion: January 29, 2026 (2 weeks)*
