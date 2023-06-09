package common

import (
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const generateTestEmailCount = 10

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateRandomEmails(amount int) []string {
	emailArr := make([]string, amount)

	for i := 0; i < amount; i++ {
		email := randSeq(8) + "@" + randSeq(5) + "." + randSeq(3)
		emailArr[i] = email
	}

	return emailArr
}

func generateRandomEmailSets() (EmailSet, EmailSet) {
	commonEmail := "developer@claudiocambra.com"
	setAEmails := append(generateRandomEmails(generateTestEmailCount), commonEmail)
	setBEmails := append(generateRandomEmails(generateTestEmailCount), commonEmail)

	emailSetA := EmailSet{}
	emailSetB := EmailSet{}

	for _, email := range setAEmails {
		emailSetA[email] = true
	}

	for _, email := range setBEmails {
		emailSetB[email] = true
	}

	return emailSetA, emailSetB
}

// EmailSet
func TestAddEmailSet(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	summedEmailSets := AddEmailSet(emailSetA, emailSetB)
	testEmailSet := emailSetA

	for email := range emailSetB {
		testEmailSet[email] = true
	}

	if !cmp.Equal(testEmailSet, summedEmailSets) {
		t.Fatalf(`Added email sets do not match expected email set: %s`, cmp.Diff(testEmailSet, summedEmailSets))
	}
}

func TestSubtractEmailSets(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	subbedEmailSets, _ := SubtractEmailSet(emailSetA, emailSetB)
	testEmailSet := emailSetA

	for email := range emailSetB {
		delete(testEmailSet, email)
	}

	if !cmp.Equal(testEmailSet, subbedEmailSets) {
		t.Fatalf(`Subtracted email sets do not match expected email set: %s`, cmp.Diff(testEmailSet, subbedEmailSets))
	}
}

// YearlyEmailMap
func TestAddEmailSetToYearlyEmailMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYear := 2023
	yem := make(YearlyEmailMap, 0)

	yem.AddEmailSet(emailSetA, testYear)
	if yemAEmailSetA := yem[testYear]; !cmp.Equal(yemAEmailSetA, emailSetA) {
		t.Fatalf(`Added email set to yearly emails map when year not already in map does not match expected changes: %s`, cmp.Diff(emailSetA, yemAEmailSetA))
	}

	yem.AddEmailSet(emailSetB, testYear)
	summedEmailSets := AddEmailSet(emailSetA, emailSetB)

	if yemASummedEmailSets := yem[testYear]; !cmp.Equal(yemASummedEmailSets, summedEmailSets) {
		t.Fatalf(`Added email set to yearly emails map when year already in map does not match expected changes:
			Expected %+v
			Received %+v`, summedEmailSets, yemASummedEmailSets)
	}
}

func TestSubtractEmailSetInYearlyEmailsMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYearA := 2023
	testYearB := 2003
	yem := YearlyEmailMap{testYearA: emailSetA}

	yem.SubtractEmailSet(emailSetB, testYearB)
	if _, ok := yem[testYearB]; ok {
		t.Fatalf("Subtracting email set from a year not present in YEM should not add this year to YEM.")
	}

	yem = YearlyEmailMap{testYearA: emailSetA}
	yem.SubtractEmailSet(emailSetB, testYearA)

	expectedSubEmailSet, _ := SubtractEmailSet(emailSetA, emailSetB)

	if subEmailSet := yem[testYearA]; !cmp.Equal(subEmailSet, expectedSubEmailSet) {
		t.Fatalf(`Subtracted email set from yearly email map does not match expected changes: %s`, cmp.Diff(expectedSubEmailSet, subEmailSet))
	}
}

func TestAddYearlyEmailMapToYearlyEmailMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYearA := 2023
	testYearB := 2003
	yemA := YearlyEmailMap{testYearA: emailSetA}
	yemB := YearlyEmailMap{testYearB: emailSetB}

	yemA.AddYearlyEmailMap(yemB)

	if addedEmailSetB, ok := yemA[testYearB]; !ok {
		t.Fatalf("Adding email set from a year of YEM B not present in YEM A should add this year to YEM A.")
	} else if !cmp.Equal(addedEmailSetB, emailSetB) {
		t.Fatalf(`Added YLCM B to YLCM A does not match expected line changes: %s`, cmp.Diff(emailSetB, addedEmailSetB))
	}

	yemB = YearlyEmailMap{testYearA: emailSetB}
	yemA.AddYearlyEmailMap(yemB)

	expectedSummedEmailSets := AddEmailSet(emailSetA, emailSetB)

	if addedEmailSets := yemA[testYearA]; !cmp.Equal(addedEmailSets, expectedSummedEmailSets) {
		t.Fatalf(`Added YLCM B email set to YLCM A does not match expected email set: %s`, cmp.Diff(addedEmailSets, expectedSummedEmailSets))
	}
}

func TestSubtractYearlyEmailMapToYearlyEmailMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYearA := 2023
	testYearB := 2003
	yemA := YearlyEmailMap{testYearA: emailSetA}
	yemB := YearlyEmailMap{testYearB: emailSetB}

	yemA.SubtractYearlyEmailMap(yemB)

	if _, ok := yemA[testYearB]; ok {
		t.Fatalf("Subtracting YEM B email set from a year not present in YEM A should not add this year to YEM A.")
	}

	yemA.AddEmailSet(emailSetA, testYearB)
	yemA.SubtractYearlyEmailMap(yemB)

	expectedSubbedEmailSets, _ := SubtractEmailSet(emailSetA, emailSetB)

	if testYearBSubEmailSet, ok := yemA[testYearB]; !ok {
		t.Fatalf("Test year B should now be present in YLCM A")
	} else if !cmp.Equal(testYearBSubEmailSet, expectedSubbedEmailSets) {
		t.Fatalf(`Subtracted YLCM B email set to yearly email set map does not match expected changes: %s`, cmp.Diff(expectedSubbedEmailSets, testYearBSubEmailSet))
	}
}
