package common

import (
	"math/rand"
	"reflect"
	"testing"
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

	if !reflect.DeepEqual(testEmailSet, summedEmailSets) {
		t.Fatalf(`Added email sets do not match expected email set: 
			Expected %+v
			Received %+v`, testEmailSet, summedEmailSets)
	}
}

func TestSubtractEmailSets(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	subbedEmailSets, _ := SubtractEmailSet(emailSetA, emailSetB)
	testEmailSet := emailSetA

	for email := range emailSetB {
		delete(testEmailSet, email)
	}

	if !reflect.DeepEqual(testEmailSet, subbedEmailSets) {
		t.Fatalf(`Subtracted email sets do not match expected email set: 
			Expected %+v
			Received %+v`, testEmailSet, subbedEmailSets)
	}
}

// YearlyEmailMap
func TestAddEmailSetToYearlyEmailMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYear := 2023
	yem := make(YearlyEmailMap, 0)

	yem.AddEmailSet(emailSetA, testYear)
	if yemAEmailSetA := yem[testYear]; !reflect.DeepEqual(yemAEmailSetA, emailSetA) {
		t.Fatalf(`Added email set to yearly emails map when year not already in map does not match expected changes:
			Expected %+v
			Received %+v`, emailSetA, yemAEmailSetA)
	}

	yem.AddEmailSet(emailSetB, testYear)
	summedEmailSets := AddEmailSet(emailSetA, emailSetB)

	if yemASummedEmailSets := yem[testYear]; !reflect.DeepEqual(yemASummedEmailSets, summedEmailSets) {
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

	if subEmailSet := yem[testYearA]; !reflect.DeepEqual(subEmailSet, expectedSubEmailSet) {
		t.Fatalf(`Subtracted email set from yearly email map does not match expected changes:
			Expected %+v
			Received %+v`, expectedSubEmailSet, subEmailSet)
	}
}