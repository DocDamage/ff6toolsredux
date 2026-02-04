package pr

// Job ID constants for character job types
// These are used to determine character-specific abilities
const (
	JobIDNone   = 0  // Default/No job
	JobIDCeles  = 1  // Celes - Runic
	JobIDCyan   = 3  // Cyan - Bushido
	JobIDEdgar  = 4  // Edgar - Tools
	JobIDLocke  = 5  // Locke - Steal/Mug
	JobIDSabin  = 6  // Sabin - Blitz
	JobIDShadow = 7  // Shadow - Throw
	JobIDStrago = 8  // Strago - Lore
	JobIDRelm   = 9  // Relm - Sketch/Control
	JobIDSetzer = 10 // Setzer - Slot/GP Toss
	JobIDGau    = 12 // Gau - Rage
	JobIDGogo   = 13 // Gogo - Mimic
)

// CharacterID constants for special character checks
const (
	CharacterIDTerra  = 0
	CharacterIDLocke  = 1
	CharacterIDCyan   = 2
	CharacterIDShadow = 3
	CharacterIDEdgar  = 4
	CharacterIDSabin  = 5
	CharacterIDCeles  = 6
	CharacterIDStrago = 7
	CharacterIDRelm   = 8
	CharacterIDSetzer = 9
	CharacterIDMog    = 10
	CharacterIDGau    = 11
	CharacterIDGogo   = 12
	CharacterIDUmaro  = 13
	CharacterIDMogMog = 16 // Mog in Moogle party (has Dance ability)
)

// IsJobWithBushido returns true if the job ID has Bushido ability (Cyan)
func IsJobWithBushido(jobID int) bool {
	return jobID == JobIDCyan
}

// IsJobWithBlitz returns true if the job ID has Blitz ability (Sabin)
func IsJobWithBlitz(jobID int) bool {
	return jobID == JobIDSabin
}

// IsJobWithLore returns true if the job ID has Lore ability (Strago)
func IsJobWithLore(jobID int) bool {
	return jobID == JobIDStrago
}

// IsJobWithRage returns true if the job ID has Rage ability (Gau)
func IsJobWithRage(jobID int) bool {
	return jobID == JobIDGau
}

// IsCharacterWithDance returns true if the character ID has Dance ability (Mog)
func IsCharacterWithDance(charID int) bool {
	return charID == CharacterIDMogMog
}
