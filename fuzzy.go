package fuzzy

import (
	"math"
)

const (
	// LowCompetence is
	LowCompetence = 50
	// NotLowCompetence is
	NotLowCompetence = 55

	// NotMiddleLowCompetence is
	NotMiddleLowCompetence = 52.5
	// MiddleLowCompetence is
	MiddleLowCompetence = 60
	// MiddleHighCompetence is
	MiddleHighCompetence = 65
	// NotMiddleHighCompetence is
	NotMiddleHighCompetence = 67.5

	// NotHighCompetence is
	NotHighCompetence = 65
	// HighCompetence is
	HighCompetence = 70
)

const (
	// LowPersonality is
	LowPersonality = 50
	// NotLowPersonality is
	NotLowPersonality = 55

	// NotMiddleLowPersonality is
	NotMiddleLowPersonality = 52.5
	// MiddleLowPersonality is
	MiddleLowPersonality = 60
	// MiddleHighPersonality is
	MiddleHighPersonality = 65
	// NotMiddleHighPersonality is
	NotMiddleHighPersonality = 75

	// NotHighPersonality is
	NotHighPersonality = 70
	//HighPersonality is
	HighPersonality = 77.5
)

const (
	// AcceptedValue is
	AcceptedValue = 100
	// RejectedValue is
	RejectedValue = 50
)

// Fuzzy is
type Fuzzy interface {
	Fuzzification(number *Number) error
	Defuzzification(number *Number) error
	Inference(number *Number) error
}

// Interview is
type Interview struct {
	ID          string
	Competence  float64
	Personality float64
}

// Number is
type Number struct {
	Interview Interview

	CompetenceMembership  []float64
	PersonalityMembership []float64

	AcceptedInference float64
	RejectedInference float64

	CrispValue float64
	Inference  string
}

// EmployeeAcceptance is
type EmployeeAcceptance struct {
}

// Fuzzification is
func (e *EmployeeAcceptance) Fuzzification(number *Number) error {
	number.CompetenceMembership = append(number.CompetenceMembership, e.CompetenceLow(number.Interview.Competence))
	number.CompetenceMembership = append(number.CompetenceMembership, e.ComptenceMiddle(number.Interview.Competence))
	number.CompetenceMembership = append(number.CompetenceMembership, e.CompetenceHigh(number.Interview.Competence))

	number.PersonalityMembership = append(number.PersonalityMembership, e.PersonalityLow(number.Interview.Personality))
	number.PersonalityMembership = append(number.PersonalityMembership, e.PersonalityMiddle(number.Interview.Personality))
	number.PersonalityMembership = append(number.PersonalityMembership, e.PersonalityHigh(number.Interview.Personality))

	return nil
}

// Inference is
func (e *EmployeeAcceptance) Inference(number *Number) error {

	lc_lp := float64(math.Min(number.CompetenceMembership[0], number.PersonalityMembership[0]))
	lc_mp := float64(math.Min(number.CompetenceMembership[0], number.PersonalityMembership[1]))
	lc_hp := float64(math.Min(number.CompetenceMembership[0], number.PersonalityMembership[2]))

	mc_lp := float64(math.Min(number.CompetenceMembership[1], number.PersonalityMembership[0]))
	mc_mp := float64(math.Min(number.CompetenceMembership[1], number.PersonalityMembership[1]))
	mc_hp := float64(math.Min(number.CompetenceMembership[1], number.PersonalityMembership[2]))

	hc_lp := float64(math.Min(number.CompetenceMembership[2], number.PersonalityMembership[0]))
	hc_mp := float64(math.Min(number.CompetenceMembership[2], number.PersonalityMembership[1]))
	hc_hp := float64(math.Min(number.CompetenceMembership[2], number.PersonalityMembership[2]))

	number.AcceptedInference = math.Max(lc_hp, mc_hp)
	number.AcceptedInference = math.Max(number.AcceptedInference, hc_mp)
	number.AcceptedInference = math.Max(number.AcceptedInference, hc_hp)

	number.RejectedInference = math.Max(lc_lp, lc_mp)
	number.RejectedInference = math.Max(number.RejectedInference, mc_lp)
	number.RejectedInference = math.Max(number.RejectedInference, mc_mp)
	number.RejectedInference = math.Max(number.RejectedInference, hc_lp)

	return nil
}

// Defuzzification is
func (e *EmployeeAcceptance) Defuzzification(number *Number) error {
	number.CrispValue = 0
	number.CrispValue += number.AcceptedInference * AcceptedValue
	number.CrispValue += number.RejectedInference * RejectedValue
	number.CrispValue /= (number.AcceptedInference + number.RejectedInference)

	if number.CrispValue > 50.0 {
		number.Inference = "Ya"
	} else {
		number.Inference = "Tidak"
	}

	return nil
}

// CompetenceLow is
func (e *EmployeeAcceptance) CompetenceLow(competence float64) float64 {
	if competence <= LowCompetence {
		return 1
	} else if competence > NotLowCompetence {
		return 0
	}
	return 1 - (float64(competence-LowCompetence) / float64(NotLowCompetence-LowCompetence))
}

// ComptenceMiddle is
func (e *EmployeeAcceptance) ComptenceMiddle(competence float64) float64 {
	if competence > MiddleLowCompetence && competence <= MiddleHighCompetence {
		return 1
	} else if competence <= NotMiddleLowCompetence || competence > NotMiddleHighCompetence {
		return 0
	} else if competence < MiddleLowCompetence && competence >= NotMiddleLowCompetence {
		return float64(competence-NotMiddleLowCompetence) / float64(MiddleLowCompetence-NotMiddleLowCompetence)
	}
	return 1 - (float64(competence-MiddleHighCompetence) / float64(NotMiddleHighCompetence-MiddleHighCompetence))
}

// CompetenceHigh is
func (e *EmployeeAcceptance) CompetenceHigh(competence float64) float64 {
	if competence <= NotHighCompetence {
		return 0
	} else if competence > HighCompetence {
		return 1
	}
	return float64(competence-NotHighCompetence) / float64(HighCompetence-NotHighCompetence)
}

// PersonalityLow is
func (e *EmployeeAcceptance) PersonalityLow(personality float64) float64 {
	if personality <= LowPersonality {
		return 1
	} else if personality > NotLowPersonality {
		return 0
	}
	return 1 - (float64(personality-LowPersonality) / float64(NotLowPersonality-LowPersonality))
}

// PersonalityMiddle is
func (e *EmployeeAcceptance) PersonalityMiddle(personality float64) float64 {
	if personality > MiddleLowPersonality && personality <= MiddleHighPersonality {
		return 1
	} else if personality < NotMiddleLowPersonality || personality > NotMiddleHighPersonality {
		return 0
	} else if personality < MiddleLowPersonality && personality >= NotMiddleLowCompetence {
		return float64(personality-NotMiddleLowPersonality) / float64(MiddleLowPersonality-NotMiddleLowPersonality)
	}
	return 1 - (float64(personality-MiddleHighPersonality) / float64(NotMiddleHighPersonality-MiddleHighPersonality))
}

// PersonalityHigh is
func (e *EmployeeAcceptance) PersonalityHigh(personality float64) float64 {
	if personality <= NotHighPersonality {
		return 0
	} else if personality > HighPersonality {
		return 1
	}
	return float64(personality-NotHighPersonality) / float64(HighPersonality-NotHighPersonality)
}
