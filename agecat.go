package agecat

import (
	"fmt"
	"sort"
	"time"
)

type Gender int

const (
	Universal Gender = iota
	Male
	Female
)

func (g Gender) Character() string {
	genders := [...]string{
		"",
		"M",
		"F",
	}

	return genders[g]
}

func (g Gender) String() string {
	genders := [...]string{
		"Universal",
		"Male",
		"Female",
	}

	return genders[g]
}

type AgeGroupType int

const (
	Juniors AgeGroupType = iota
	Masters
)

func (a AgeGroupType) isUnder() bool {
	return a == Juniors
}

var Location, _ = time.LoadLocation("UTC")

// AgeGroups are used to determine the category for junior and masters
// athletes.
//
// AgeGroups should be defined as per your competition rules; you can include
// or omit AgeGroups as necessary
type AgeGroups struct {
	// The gender category for the age groups
	// Groups exist for "universal", "male" and "female"
	Gender Gender
	// The type of age group
	// The type can be "Juniors" or "Masters", e.g. Under 13, Over 35, etc.
	AgeGroupType AgeGroupType
	// OperativeDate is the date used to determine the age
	// category for a junior athlete, as per the rules of your competition.
	//
	// RULE 141 S 1 UKA SUPPLEMENT
	// ADDITIONAL UKA AGE GROUPS
	// "The Competition Year for Road Running Events shall be from 1st
	// September each year to the following 31st August. The Competition Year
	// for all other disciplines shall be from 1st October each year to the
	// following 30th September.
	//
	// "The operative date for determining membership of age groups for all
	// athletes under the age of 17 shall be for Track and Field and Race
	// Walking, the 31st August at the end of the Competition Year, and for all
	// other disciplines, the 31st August prior to the commencement of the
	// Competition Year."
	//
	// So, a Track and Field event taking place in the competition year
	// starting on 1st October 2019 and finishing on 30th September 2020 would
	// have an OperativeDate for junior age groups of 31st August 2020, whereas
	// for a Road event taking place in the competition year starting on
	// 1st September 2019 and finishing on 31st August 2020, the OperativeDate
	// for junior age groups would be 31st August 2019.
	OperativeDate time.Time
	// CutOffDate is the date at which an athlete must be under
	// the age of the final age group in the Groups list below, when the
	// AgeGroupType is "Juniors".
	// If the value is not set, the OperativeDate applies.
	// If the AgeGroupType is "Masters" and the CutOffDate is set, it is ignored.
	//
	// This option exists because certain competitions use a different date
	// to the operative date to determine the end of the junior age groups; e.g:
	//
	// "Track and Field events for Under 20s shall be confined to competitors
	// who are 17 or over on 31st August within the Competition Year, but
	// Under 20 on 31st December in the calendar year of competition."
	CutOffDate *time.Time
	// The range of age groups available
	// e.g. []int{13, 15, 17, 20}
	// Combined with the AgeGroupType defined above to determine membership;
	// e.g. Under 15 would contain athletes aged 13 and 14
	Groups []int
}

// AgeCategory returns a string describing the age category for an athlete with
// the specified Gender and Date Of Birth, given the supplied AgeGroup
// constraints.
//
// If no AgeGroup constraints are supplied, the athlete will always be classed
// in the "Senior" category
//
func AgeCategory(gender Gender, birthYear int, birthMonth time.Month, birthDay int, ageGroups ...*AgeGroups) string {
	dob := time.Date(birthYear, birthMonth, birthDay, 0, 0, 0, 0, Location)

	for _, ageGroup := range ageGroups {
		if s := ageGroup.categorize(gender, dob); s != nil {
			return *s
		}
	}

	return fmt.Sprintf("%sSEN", gender.Character())
}

func (a AgeGroups) categorize(gender Gender, dob time.Time) *string {
	if a.Gender != gender {
		return nil
	}

	if a.AgeGroupType == Juniors {
		return a.categorizeUnder(gender, dob)
	}

	return a.categorizeOver(gender, dob)
}

func (a AgeGroups) categorizeUnder(gender Gender, dob time.Time) *string {
	sort.Ints(a.Groups)

	if a.CutOffDate != nil {
		if ageOnDate(*a.CutOffDate, dob) >= a.Groups[len(a.Groups)-1] {
			return nil
		}
	}

	age := ageOnDate(a.OperativeDate, dob)

	if age > a.Groups[len(a.Groups)-1] {
		return nil
	}

	if a.CutOffDate != nil &&
		age == a.Groups[len(a.Groups)-1] {
		s := fmt.Sprintf("J%s%d", gender.Character(), age)
		return &s
	}

	for _, ageGroup := range a.Groups {
		if age < ageGroup {
			s := fmt.Sprintf("J%s%d", gender.Character(), ageGroup)
			return &s
		}
	}

	return nil
}

func (a AgeGroups) categorizeOver(gender Gender, dob time.Time) *string {
	sort.Ints(a.Groups)

	age := ageOnDate(a.OperativeDate, dob)

	if age < a.Groups[0] {
		return nil
	}

	for i := len(a.Groups) - 1; i >= 0; i-- {
		ageGroup := a.Groups[i]
		if age >= ageGroup {
			s := fmt.Sprintf("%sV%d", gender.Character(), ageGroup)
			return &s
		}
	}

	return nil
}

func ageOnDate(opDate, dob time.Time) int {
	age := opDate.Year() - dob.Year()

	if dob.Month() > opDate.Month() ||
		dob.Month() == opDate.Month() &&
			dob.Day() > opDate.Day() {
		age -= 1
	}

	return age
}
