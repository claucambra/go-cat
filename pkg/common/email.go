package common

type EmailSet map[string]bool
type YearlyEmailMap map[int]EmailSet

func AddEmailSet(a EmailSet, b EmailSet) EmailSet {
	newSet := CopyMap(a)

	for email := range b {
		newSet[email] = true
	}

	return newSet
}

func SubtractEmailSet(a EmailSet, b EmailSet) (EmailSet, bool) {
	newSet := CopyMap(a)

	for email := range b {
		delete(newSet, email)
	}

	empty := len(newSet) == 0
	return newSet, empty
}

func (yem *YearlyEmailMap) CountArray(years []int) []int {
	yearsToReturn := years

	if yearsToReturn == nil {
		yearsToReturn = SortedMapKeys(*yem)
	}

	countArray := make([]int, len(yearsToReturn))

	for i, year := range yearsToReturn {
		countArray[i] = len((*yem)[year])
	}

	return countArray
}

func (yem *YearlyEmailMap) AddEmailSet(emailSetToAdd EmailSet, year int) {
	AdditiveValueMapInsert[int, EmailSet, YearlyEmailMap](*yem, year, AddEmailSet, emailSetToAdd)
}

func (yem *YearlyEmailMap) SubtractEmailSet(emailSetToSubtract EmailSet, year int) {
	SubtractiveValueMapRemove[int, EmailSet, YearlyEmailMap](*yem, year, SubtractEmailSet, emailSetToSubtract)
}

func (yem *YearlyEmailMap) AddYearlyEmailMap(yemToAdd YearlyEmailMap) {
	for year, emailsToAdd := range yemToAdd {
		yem.AddEmailSet(emailsToAdd, year)
	}
}

func (yem *YearlyEmailMap) SubtractYearlyEmailMap(yemToSubtract YearlyEmailMap) {
	for year, emailsToSubtract := range yemToSubtract {
		yem.SubtractEmailSet(emailsToSubtract, year)
	}
}
