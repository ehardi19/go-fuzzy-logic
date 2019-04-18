package fuzzy

import "math"

const (
	// MustBeLowCompetence is
	MustBeLowCompetence = 30
	// MustNotBeLowCompetence is
	MustNotBeLowCompetence = 50

	// MustNotBeMiddleLowCompetence is
	MustNotBeMiddleLowCompetence = 50
	// MustBeMiddleLowCompetence is
	MustBeMiddleLowCompetence = 65
	// MustBeMiddleHighCompetence is
	MustBeMiddleHighCompetence = 75
	// MustNotBeMiddleHighCompetence is
	MustNotBeMiddleHighCompetence = 80

	// MustNotBeHighCompetence is
	MustNotBeHighCompetence = 65
	// MustBeHighCompetence is
	MustBeHighCompetence = 85
)

const (
	// MustBeLowPersonality is
	MustBeLowPersonality = 15
	// MustNotBeLowPersonality is
	MustNotBeLowPersonality = 35

	// MustNotBeMiddleLowPersonality is
	MustNotBeMiddleLowPersonality = 10
	// MustBeMiddleLowPersonality is
	MustBeMiddleLowPersonality = 40
	// MustBeMiddleHighPersonality is
	MustBeMiddleHighPersonality = 50
	// MustNotBeMiddleHighPersonality is
	MustNotBeMiddleHighPersonality = 75

	// MustNotBeHighPersonality is
	MustNotBeHighPersonality = 45
	//MustBeHighPersonality is
	MustBeHighPersonality = 70
)

const (
	// AcceptedValue is
	AcceptedValue = 100
	// ConsideredValue is
	ConsideredValue = 70
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

	AcceptedInference   float64
	ConsideredInference float64
	RejectedInference   float64

	CrispValue float64
}

// BLT is
type BLT struct {
}

// Fuzzification is
func (b *BLT) Fuzzification(number *Number) error {
	number.CompetenceMembership = append(number.CompetenceMembership, b.CompetenceLow(number.Interview.Competence))
	number.CompetenceMembership = append(number.CompetenceMembership, b.ComptenceMiddle(number.Interview.Competence))
	number.CompetenceMembership = append(number.CompetenceMembership, b.CompetenceHigh(number.Interview.Competence))

	number.PersonalityMembership = append(number.PersonalityMembership, b.PersonalityLow(number.Interview.Personality))
	number.PersonalityMembership = append(number.PersonalityMembership, b.PersonalityMiddle(number.Interview.Personality))
	number.PersonalityMembership = append(number.PersonalityMembership, b.PersonalityHigh(number.Interview.Personality))

	return nil
}

// Defuzzification is
func (b *BLT) Defuzzification(number *Number) error {
	number.CrispValue = 0
	number.CrispValue += number.AcceptedInference * AcceptedValue
	number.CrispValue += number.ConsideredInference * ConsideredValue
	number.CrispValue += number.RejectedInference * RejectedValue
	number.CrispValue /= (number.AcceptedInference + number.ConsideredInference + number.RejectedInference)

	return nil
}

// Inference is
func (b *BLT) Inference(number *Number) error {
	number.AcceptedInference = math.Max(math.Min(number.CompetenceMembership[0], number.PersonalityMembership[2]), math.Min(number.CompetenceMembership[0], number.PersonalityMembership[1]))
	number.AcceptedInference = math.Max(number.AcceptedInference, math.Min(number.CompetenceMembership[0], number.PersonalityMembership[0]))

	number.ConsideredInference = math.Max(math.Min(number.CompetenceMembership[1], number.PersonalityMembership[1]), math.Min(number.CompetenceMembership[1], number.PersonalityMembership[2]))
	number.ConsideredInference = math.Max(number.ConsideredInference, math.Min(number.CompetenceMembership[2], number.PersonalityMembership[2]))

	number.RejectedInference = math.Max(math.Min(number.CompetenceMembership[2], number.PersonalityMembership[0]), math.Min(number.CompetenceMembership[2], number.PersonalityMembership[1]))
	number.RejectedInference = math.Max(number.RejectedInference, math.Min(number.CompetenceMembership[1], number.PersonalityMembership[0]))

	return nil
}

// CompetenceLow is
func (b *BLT) CompetenceLow(competence float64) float64 {
	if competence <= MustBeLowCompetence {
		return 1
	} else if competence > MustNotBeLowCompetence {
		return 0
	}
	return 1 - (float64(competence-MustBeLowCompetence) / float64(MustNotBeLowPersonality-MustBeLowCompetence))
}

// ComptenceMiddle is
func (b *BLT) ComptenceMiddle(competence float64) float64 {
	if competence > MustBeMiddleLowCompetence && competence <= MustBeMiddleHighCompetence {
		return 1
	} else if competence < MustNotBeMiddleLowCompetence || competence > MustNotBeMiddleHighCompetence {
		return 0
	} else if competence < MustBeMiddleLowCompetence && competence >= MustNotBeMiddleLowCompetence {
		return float64(competence-MustBeMiddleLowCompetence) / float64(MustBeLowCompetence-MustNotBeLowCompetence)
	}
	return 1 - (float64(competence-MustBeMiddleHighCompetence) / float64(MustNotBeMiddleHighCompetence-MustBeMiddleHighCompetence))
}

// CompetenceHigh is
func (b *BLT) CompetenceHigh(competence float64) float64 {
	if competence <= MustNotBeHighCompetence {
		return 0
	} else if competence > MustBeHighCompetence {
		return 1
	}
	return float64(competence-MustNotBeHighCompetence) / float64(MustBeHighCompetence-MustNotBeHighCompetence)
}

// PersonalityLow is
func (b *BLT) PersonalityLow(personality float64) float64 {
	if personality <= MustBeLowPersonality {
		return 1
	} else if personality > MustNotBeLowPersonality {
		return 0
	}
	return 1 - (float64(personality-MustBeLowPersonality) / float64(MustNotBeLowPersonality-MustBeLowPersonality))
}

// PersonalityMiddle is
func (b *BLT) PersonalityMiddle(personality float64) float64 {
	if personality > MustBeMiddleLowPersonality && personality <= MustBeMiddleHighPersonality {
		return 1
	} else if personality < MustNotBeMiddleLowPersonality || personality > MustNotBeMiddleHighPersonality {
		return 0
	} else if personality < MustBeMiddleLowPersonality && personality >= MustNotBeMiddleLowPersonality {
		return float64(personality-MustBeMiddleLowPersonality) / float64(MustBeLowPersonality-MustNotBeLowPersonality)
	}
	return 1 - (float64(personality-MustBeMiddleHighPersonality) / float64(MustNotBeMiddleHighPersonality-MustBeMiddleHighPersonality))
}

// PersonalityHigh is
func (b *BLT) PersonalityHigh(personality float64) float64 {
	if personality <= MustNotBeHighPersonality {
		return 0
	} else if personality > MustBeHighPersonality {
		return 1
	}
	return float64(personality-MustNotBeHighPersonality) / float64(MustBeHighPersonality-MustNotBeHighPersonality)
}
