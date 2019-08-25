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
	AgeGroups           []*AgeGroups
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
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth.Year(), testCase.DateOfBirth.Month(), testCase.DateOfBirth.Day())

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})

	t.Run("junior age groups", func(t *testing.T) {
		var testCases []*AgeCategoryTestCase
		var dob time.Time

		ageGroups := make([]*AgeGroups, 2)

		ageGroups[0] = &AgeGroups{
			Gender:        Female,
			AgeGroupType:  Juniors,
			OperativeDate: time.Now().UTC().Truncate(time.Hour * 24),
			CutOffDate:    nil,
			Groups:        []int{13, 15, 17, 20},
		}

		ageGroups[1] = &AgeGroups{
			Gender:        Male,
			AgeGroupType:  Juniors,
			OperativeDate: time.Now().UTC().Truncate(time.Hour * 24),
			CutOffDate:    nil,
			Groups:        []int{13, 15, 17, 20},
		}

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
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth.Year(), testCase.DateOfBirth.Month(), testCase.DateOfBirth.Day(), ageGroups...)

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})

	t.Run("masters age groups", func(t *testing.T) {
		var testCases []*AgeCategoryTestCase
		var dob time.Time

		ageGroups := make([]*AgeGroups, 2)

		// Many races award prizes for women over 35, but not men over 35.
		// I don't know the reason why... misogyny? (Or misandry!)
		ageGroups[0] = &AgeGroups{
			Gender:        Female,
			AgeGroupType:  Masters,
			OperativeDate: time.Now().UTC().Truncate(time.Hour * 24),
			CutOffDate:    nil,
			Groups:        []int{35, 40, 45, 50, 55, 60, 65, 70, 75},
		}

		ageGroups[1] = &AgeGroups{
			Gender:        Male,
			AgeGroupType:  Masters,
			OperativeDate: time.Now().UTC().Truncate(time.Hour * 24),
			CutOffDate:    nil,
			Groups:        []int{40, 45, 50, 55, 60, 65, 70, 75},
		}

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
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth.Year(), testCase.DateOfBirth.Month(), testCase.DateOfBirth.Day(), ageGroups...)

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

		ageGroups := make([]*AgeGroups, 2)

		ageGroups[0] = &AgeGroups{
			Gender:        Female,
			AgeGroupType:  Juniors,
			OperativeDate: operativeDate,
			CutOffDate:    &cutOffDate,
			Groups:        []int{13, 15, 17, 20},
		}

		ageGroups[1] = &AgeGroups{
			Gender:        Male,
			AgeGroupType:  Juniors,
			OperativeDate: operativeDate,
			CutOffDate:    &cutOffDate,
			Groups:        []int{13, 15, 17, 20},
		}

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
			got := AgeCategory(testCase.Gender, testCase.DateOfBirth.Year(), testCase.DateOfBirth.Month(), testCase.DateOfBirth.Day(), ageGroups...)

			if got != testCase.ExpectedAgeCategory {
				t.Errorf("got %s, want %s for date of birth %s", got, testCase.ExpectedAgeCategory, testCase.DateOfBirth.Format("2006-01-02"))
			}
		}
	})
}
