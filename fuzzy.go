package fuzzy

import (
	"math"
)

const (
	// LowCompetence is const value for competence that is low
	LowCompetence = 50
	// NotLowCompetence is const value for competence that is not low
	NotLowCompetence = 55

	// NotMiddleLowCompetence is
	NotMiddleLowCompetence = 52.5
	// MiddleLowCompetence is
	MiddleLowCompetence = 60
	// MiddleHighCompetence is
	MiddleHighCompetence = 65
	// NotMiddleHighCompetence is
	NotMiddleHighCompetence = 67.5

	// NotHighCompetence is const value for competence that is not high
	NotHighCompetence = 65
	// HighCompetence is const value for competence that is high
	HighCompetence = 70
)

const (
	// LowPersonality is const value for personality that is low
	LowPersonality = 50
	// NotLowPersonality is const value for personality that is not low
	NotLowPersonality = 55

	// NotMiddleLowPersonality is
	NotMiddleLowPersonality = 52.5
	// MiddleLowPersonality is
	MiddleLowPersonality = 60
	// MiddleHighPersonality is
	MiddleHighPersonality = 65
	// NotMiddleHighPersonality is
	NotMiddleHighPersonality = 75

	// NotHighPersonality is const value for personality that is not high
	NotHighPersonality = 70
	//HighPersonality is const value for personality that is high
	HighPersonality = 77.5
)

const (
	// AcceptedValue is const value of Takagi-Sugeno Method, that the value is accepted
	AcceptedValue = 100
	// RejectedValue is const value of Takagi-Sugeno Method, that the value is rejected
	RejectedValue = 50
)

// Fuzzy is the main interface of a fuzzy logic algorithm
type Fuzzy interface {
	Fuzzification(number *Number) error
	Defuzzification(number *Number) error
	Inference(number *Number) error
}

// Interview is the struct for data needed for this program
type Interview struct {
	ID          string
	Competence  float64
	Personality float64
}

// Number is the struct that holds fuzzy data
type Number struct {
	Interview Interview

	CompetenceMembership  []float64
	PersonalityMembership []float64

	AcceptedInference float64
	RejectedInference float64

	CrispValue float64
	Inference  string
}

// EmployeeAcceptance is the struct that holds the fuzzy algorithm process
type EmployeeAcceptance struct {
}

// Fuzzification is a function that will transfer crisp data into linguistic
func (e *EmployeeAcceptance) Fuzzification(number *Number) error {
	number.CompetenceMembership = append(number.CompetenceMembership, e.CompetenceLow(number.Interview.Competence))
	number.CompetenceMembership = append(number.CompetenceMembership, e.ComptenceMiddle(number.Interview.Competence))
	number.CompetenceMembership = append(number.CompetenceMembership, e.CompetenceHigh(number.Interview.Competence))

	number.PersonalityMembership = append(number.PersonalityMembership, e.PersonalityLow(number.Interview.Personality))
	number.PersonalityMembership = append(number.PersonalityMembership, e.PersonalityMiddle(number.Interview.Personality))
	number.PersonalityMembership = append(number.PersonalityMembership, e.PersonalityHigh(number.Interview.Personality))

	return nil
}

// Inference is a function that change from raw linguistic into fuzzy linguistic
func (e *EmployeeAcceptance) Inference(number *Number) error {

	lcLp := float64(math.Min(number.CompetenceMembership[0], number.PersonalityMembership[0]))
	lcMp := float64(math.Min(number.CompetenceMembership[0], number.PersonalityMembership[1]))
	lcHp := float64(math.Min(number.CompetenceMembership[0], number.PersonalityMembership[2]))

	mcLp := float64(math.Min(number.CompetenceMembership[1], number.PersonalityMembership[0]))
	mcMp := float64(math.Min(number.CompetenceMembership[1], number.PersonalityMembership[1]))
	mcHp := float64(math.Min(number.CompetenceMembership[1], number.PersonalityMembership[2]))

	hcLp := float64(math.Min(number.CompetenceMembership[2], number.PersonalityMembership[0]))
	hcMp := float64(math.Min(number.CompetenceMembership[2], number.PersonalityMembership[1]))
	hcHp := float64(math.Min(number.CompetenceMembership[2], number.PersonalityMembership[2]))

	number.AcceptedInference = math.Max(lcHp, mcHp)
	number.AcceptedInference = math.Max(number.AcceptedInference, hcMp)
	number.AcceptedInference = math.Max(number.AcceptedInference, hcHp)

	number.RejectedInference = math.Max(lcLp, lcMp)
	number.RejectedInference = math.Max(number.RejectedInference, mcLp)
	number.RejectedInference = math.Max(number.RejectedInference, mcMp)
	number.RejectedInference = math.Max(number.RejectedInference, hcLp)

	return nil
}

// Defuzzification is a function that will transfer fuzzy linguistic into crisp data
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

// CompetenceLow is a function that determine low competence value
func (e *EmployeeAcceptance) CompetenceLow(competence float64) float64 {
	if competence <= LowCompetence {
		return 1
	} else if competence > NotLowCompetence {
		return 0
	}
	return 1 - (float64(competence-LowCompetence) / float64(NotLowCompetence-LowCompetence))
}

// ComptenceMiddle is a function that determine middle competence value
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

// CompetenceHigh is a function that determine high competence value
func (e *EmployeeAcceptance) CompetenceHigh(competence float64) float64 {
	if competence <= NotHighCompetence {
		return 0
	} else if competence > HighCompetence {
		return 1
	}
	return float64(competence-NotHighCompetence) / float64(HighCompetence-NotHighCompetence)
}

// PersonalityLow is a function that determine low personality value
func (e *EmployeeAcceptance) PersonalityLow(personality float64) float64 {
	if personality <= LowPersonality {
		return 1
	} else if personality > NotLowPersonality {
		return 0
	}
	return 1 - (float64(personality-LowPersonality) / float64(NotLowPersonality-LowPersonality))
}

// PersonalityMiddle is a function that determine middle personality value
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

// PersonalityHigh is a function that determine high personality value
func (e *EmployeeAcceptance) PersonalityHigh(personality float64) float64 {
	if personality <= NotHighPersonality {
		return 0
	} else if personality > HighPersonality {
		return 1
	}
	return float64(personality-NotHighPersonality) / float64(HighPersonality-NotHighPersonality)
}
