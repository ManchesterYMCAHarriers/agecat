package agecat

import (
	"fmt"
	"testing"
	"time"
)

type AgeCategoryTestCase struct {
	Gender              Gender
	DateOfBirth         time.Time
	ExpectedAgeCategory string
	CategoryGroups      []*categoryGroup
}

func TestAgeCategory(t *testing.T) {
	t.Run("no age groups", func(t *testing.T) {
		var testCases []*AgeCategoryTestCase
		var dob time.Time

		dob = time.Now().AddDate(-11, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FSEN",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "MSEN",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Universal,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "SEN",
		})

		dob = time.Now().AddDate(-90, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FSEN",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "MSEN",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Universal,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "SEN",
		})

		for _, testCase := range testCases {
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth)

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})

	t.Run("junior age groups", func(t *testing.T) {
		var testCases []*AgeCategoryTestCase
		var dob time.Time

		categoryGroups := make([]*categoryGroup, 2)

		categoryGroups[0] = NewCategoryGroup(Female, Juniors, time.Now().UTC().Truncate(time.Hour*24), nil, []int{13, 15, 17, 20})
		categoryGroups[1] = NewCategoryGroup(Male, Juniors, time.Now().UTC().Truncate(time.Hour*24), nil, []int{13, 15, 17, 20})

		dob = time.Now().AddDate(-12, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF13",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JM13",
		})

		dob = time.Now().AddDate(-10, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF13",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JM13",
		})

		for y := 13; y < 18; y += 2 {
			dob = time.Now().AddDate(-y+2, 0, 0).UTC().Truncate(time.Hour * 24)
			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Female,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JF%d", y),
			})

			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Male,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JM%d", y),
			})

			dob = time.Now().AddDate(-y, 0, 1).UTC().Truncate(time.Hour * 24)
			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Female,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JF%d", y),
			})

			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Male,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JM%d", y),
			})
		}

		dob = time.Now().AddDate(-17, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF20",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JM20",
		})

		dob = time.Now().AddDate(-20, 0, 1).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF20",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JM20",
		})

		dob = time.Now().AddDate(-20, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FSEN",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "MSEN",
		})

		for _, testCase := range testCases {
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth, categoryGroups...)

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})

	t.Run("masters age groups", func(t *testing.T) {
		var testCases []*AgeCategoryTestCase
		var dob time.Time

		categoryGroups := make([]*categoryGroup, 2)

		// Many races award prizes for women over 35, but not men over 35.
		// I don't know the reason why... misogyny? (Or misandry!)
		categoryGroups[0] = NewCategoryGroup(Female, Masters, time.Now().UTC().Truncate(time.Hour*24), nil, []int{35, 40, 45, 50, 55, 60, 65, 70, 75})
		categoryGroups[1] = NewCategoryGroup(Male, Masters, time.Now().UTC().Truncate(time.Hour*24), nil, []int{40, 45, 50, 55, 60, 65, 70, 75})

		dob = time.Now().AddDate(-35, 0, 1).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FSEN",
		})

		dob = time.Now().AddDate(-35, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FV35",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "MSEN",
		})

		for y := 40; y < 76; y += 5 {
			dob = time.Now().AddDate(-y, 0, 0).UTC().Truncate(time.Hour * 24)
			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Female,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("FV%d", y),
			})

			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Male,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("MV%d", y),
			})

			dob = time.Now().AddDate(-y-5, 0, 1).UTC().Truncate(time.Hour * 24)
			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Female,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("FV%d", y),
			})

			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Male,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("MV%d", y),
			})
		}

		dob = time.Now().AddDate(-80, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: fmt.Sprintf("FV75"),
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: fmt.Sprintf("MV75"),
		})

		for _, testCase := range testCases {
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth, categoryGroups...)

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})

	t.Run("junior age groups with cut off date", func(t *testing.T) {
		var testCases []*AgeCategoryTestCase
		var dob time.Time
		operativeDate := time.Date(2020, 8, 31, 0, 0, 0, 0, Location)
		cutOffDate := time.Date(2019, 12, 31, 0, 0, 0, 0, Location)

		categoryGroups := make([]*categoryGroup, 2)

		categoryGroups[0] = NewCategoryGroup(Female, Juniors, operativeDate, &cutOffDate, []int{13, 15, 17, 20})
		categoryGroups[1] = NewCategoryGroup(Male, Juniors, operativeDate, &cutOffDate, []int{13, 15, 17, 20})

		dob = operativeDate.AddDate(-10, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF13",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JM13",
		})

		for y := 13; y < 18; y += 2 {
			dob = operativeDate.AddDate(-y+2, 0, 0).UTC().Truncate(time.Hour * 24)
			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Female,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JF%d", y),
			})

			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Male,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JM%d", y),
			})

			dob = operativeDate.AddDate(-y, 0, 1).UTC().Truncate(time.Hour * 24)
			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Female,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JF%d", y),
			})

			testCases = append(testCases, &AgeCategoryTestCase{
				Gender:              Male,
				DateOfBirth:         dob,
				ExpectedAgeCategory: fmt.Sprintf("JM%d", y),
			})
		}

		dob = operativeDate.AddDate(-17, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF20",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JM20",
		})

		dob = cutOffDate.AddDate(-20, 0, 1).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF20",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JM20",
		})

		dob = cutOffDate.AddDate(-20, 0, 0).UTC().Truncate(time.Hour * 24)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FSEN",
		})

		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Male,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "MSEN",
		})

		for _, testCase := range testCases {
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth, categoryGroups...)

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})

	t.Run("timezones don't affect age calculations", func(t *testing.T) {
		var testCases []*AgeCategoryTestCase
		var dob, operativeDate, cutOffDate time.Time
		var midway, fiji *time.Location
		midway, _ = time.LoadLocation("Pacific/Midway")
		fiji, _ = time.LoadLocation("Pacific/Fiji")
		operativeDate = time.Date(2020, 8, 31, 0, 0, 0, 0, midway)
		cutOffDate = time.Date(2019, 12, 31, 23, 59, 59, 999, midway)

		categoryGroups := make([]*categoryGroup, 1)

		categoryGroups[0] = NewCategoryGroup(Female, Juniors, operativeDate, &cutOffDate, []int{13, 15, 17, 20})

		dob = time.Date(2007, 8, 31, 23, 59, 59, 999, fiji)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF15",
		})

		dob = time.Date(2007, 9, 1, 0, 0, 0, 0, fiji)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF13",
		})

		dob = time.Date(1999, 12, 31, 23, 59, 59, 999, fiji)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FSEN",
		})

		dob = time.Date(2000, 1, 1, 0, 0, 0, 0, fiji)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF20",
		})

		operativeDate = time.Date(2020, 8, 31, 0, 0, 0, 0, fiji)
		cutOffDate = time.Date(2019, 12, 31, 23, 59, 59, 999, fiji)

		categoryGroups[0] = NewCategoryGroup(Female, Juniors, operativeDate, &cutOffDate, []int{13, 15, 17, 20})

		dob = time.Date(2007, 8, 31, 23, 59, 59, 999, midway)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF15",
		})

		dob = time.Date(2007, 9, 1, 0, 0, 0, 0, midway)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF13",
		})

		dob = time.Date(1999, 12, 31, 23, 59, 59, 999, midway)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "FSEN",
		})

		dob = time.Date(2000, 1, 1, 0, 0, 0, 0, midway)
		testCases = append(testCases, &AgeCategoryTestCase{
			Gender:              Female,
			DateOfBirth:         dob,
			ExpectedAgeCategory: "JF20",
		})

		for _, testCase := range testCases {
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth, categoryGroups...)

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})
}
