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

type CategoryGroupType int

const (
	Juniors CategoryGroupType = iota
	Masters
)

func (a CategoryGroupType) isJunior() bool {
	return a == Juniors
}

var Location, _ = time.LoadLocation("UTC")

// Category groups are used to determine the category for junior and masters
// athletes.
//
// Category groups should be defined as per your competition rules; you can
// include/omit categoryGroup as necessary; e.g. if your competition doesn't
// cater for juniors, only include CategoryGroups for masters.
//
// Create a category group with the NewCategoryGroup function.
type categoryGroup struct {
	// The gender category for the age groups
	// Groups exist for "universal", "male" and "female"
	Gender Gender
	// The type of categoryGroup
	// The type can be "Juniors" or "Masters", e.g. Under 13, Over 35, etc.
	CategoryGroupType CategoryGroupType
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
	// the age of the final age group in the Ages list below, when the
	// CategoryGroupType is "Juniors".
	// If the value is not set, the OperativeDate applies.
	// If the CategoryGroupType is "Masters" and the CutOffDate is set, it is ignored.
	//
	// This option exists because certain competitions use a different date
	// to the operative date to determine the end of the junior age groups; e.g:
	//
	// "Track and Field events for Under 20s shall be confined to competitors
	// who are 17 or over on 31st August within the Competition Year, but
	// Under 20 on 31st December in the calendar year of competition."
	CutOffDate *time.Time
	// The range of ages in the group
	// e.g. []int{13, 15, 17, 20}
	// Combined with the CategoryGroupType defined above to determine membership;
	// e.g. Under 15 would contain athletes aged 13 and 14
	Ages []int
}

// AgeCategory returns a string describing the age category for an athlete with
// the specified Gender and Date Of Birth, given the constraints in the
// supplied category group(s).
//
// If no category groups are supplied, the athlete will always be classed
// in the "Senior" category.
//
// Junior formats appear as J + {Gender} + {Age Category},
// e.g. "JF13" for Under 13 Girls
//
// Masters formats appear as {Gender} + V + {Age Category},
// e.g. "MV40" for Men Over 40
//
// Senior categories are returned as {Gender} + SEN,
// e.g. "FSEN" for Senior Women
//
// For Universal races (i.e. where there is no classification by gender), the
// gender element is omitted.
// e.g. "V60" for People Over 60
func AgeCategory(gender Gender, dateOfBirth time.Time, categoryGroups ...*categoryGroup) string {
	for _, categoryGroup := range categoryGroups {
		if s := categoryGroup.categorize(gender, dateOfBirth); s != nil {
			return *s
		}
	}

	return fmt.Sprintf("%sSEN", gender.Character())
}

// NewCategoryGroup creates a categoryGroup and returns its reference
// The meanings of each parameter are explained in the categoryGroup struct
func NewCategoryGroup(gender Gender, categoryGroupType CategoryGroupType, operativeDate time.Time, cutOffDate *time.Time, ages []int) *categoryGroup {
	return &categoryGroup{
		Gender:            gender,
		CategoryGroupType: categoryGroupType,
		OperativeDate:     operativeDate,
		CutOffDate:        cutOffDate,
		Ages:              ages,
	}
}

func (c *categoryGroup) categorize(gender Gender, dob time.Time) *string {
	if c.Gender != gender {
		return nil
	}

	if c.CategoryGroupType.isJunior() {
		return c.categorizeJuniors(gender, dob)
	}

	return c.categorizeMasters(gender, dob)
}

func (c *categoryGroup) categorizeJuniors(gender Gender, dob time.Time) *string {
	sort.Ints(c.Ages)

	if c.CutOffDate != nil {
		if ageOnDate(*c.CutOffDate, dob) >= c.Ages[len(c.Ages)-1] {
			return nil
		}
	}

	age := ageOnDate(c.OperativeDate, dob)

	if age > c.Ages[len(c.Ages)-1] {
		return nil
	}

	if c.CutOffDate != nil &&
		age == c.Ages[len(c.Ages)-1] {
		s := fmt.Sprintf("J%s%d", gender.Character(), age)
		return &s
	}

	for _, ageGroup := range c.Ages {
		if age < ageGroup {
			s := fmt.Sprintf("J%s%d", gender.Character(), ageGroup)
			return &s
		}
	}

	return nil
}

func (c *categoryGroup) categorizeMasters(gender Gender, dob time.Time) *string {
	sort.Ints(c.Ages)

	age := ageOnDate(c.OperativeDate, dob)

	if age < c.Ages[0] {
		return nil
	}

	for i := len(c.Ages) - 1; i >= 0; i-- {
		ageGroup := c.Ages[i]
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
