package gohelper

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

type TestStruct struct {
	ID    int
	Value string
}

func TestInArray(t *testing.T) {
	t.Run("testing Int", testInt)
	t.Run("testing Uint", testUint)
	t.Run("testing Float", testFloat)
	t.Run("testing Complex", testComplex)
	t.Run("testing String", testString)
	t.Run("testing Uintptr", testUintptr)
	t.Run("testing Struct", testStruct)
}

func testInt(t *testing.T) {
	// simple int exists
	var haystack = []int{0, 1, 2, 3, 4, 5}
	var needle int = 3

	isExistsIntArray, _, _ := InArray(haystack, needle)
	if !isExistsIntArray {
		t.Log("can't find that int value")
		t.FailNow()
	}

	// simple int not exists
	haystack = []int{0, 1, 2, 3, 4, 5}
	needle = 6

	isExistsIntArray, _, _ = InArray(haystack, needle)
	if isExistsIntArray {
		t.Log("shouldn't have find that int value")
		t.FailNow()
	}
}

func testUint(t *testing.T) {
	// simple uint exists
	var haystack = []uint{0, 1, 2, 3, 4, 5}
	var needle uint = 3

	isExistsIntArray, _, _ := InArray(haystack, needle)
	if !isExistsIntArray {
		t.Log("can't find that uint value")
		t.FailNow()
	}

	// simple uint not exists
	haystack = []uint{0, 1, 2, 3, 4, 5}
	needle = 6

	isExistsIntArray, _, _ = InArray(haystack, needle)
	if isExistsIntArray {
		t.Log("shouldn't have find that uint value")
		t.FailNow()
	}
}

func testFloat(t *testing.T) {
	// simple float32 exists
	var haystack = []float32{0, 1, 2, 3, 4, 5}
	var needle float32 = 3

	isExistsIntArray, _, _ := InArray(haystack, needle)
	if !isExistsIntArray {
		t.Log("can't find that float32 value")
		t.FailNow()
	}

	// simple float32 not exists
	haystack = []float32{0, 1, 2, 3, 4, 5}
	needle = 6

	isExistsIntArray, _, _ = InArray(haystack, needle)
	if isExistsIntArray {
		t.Log("shouldn't have find that float32 value")
		t.FailNow()
	}
}

func testComplex(t *testing.T) {
	// simple complex64 exists
	var haystack = []complex64{0, 1, 2, 3, 4, 5}
	var needle complex64 = 3

	isExistsIntArray, _, _ := InArray(haystack, needle)
	if !isExistsIntArray {
		t.Log("can't find that complex64 value")
		t.FailNow()
	}

	// simple int complex64 exists
	haystack = []complex64{0, 1, 2, 3, 4, 5}
	needle = 6

	isExistsIntArray, _, _ = InArray(haystack, needle)
	if isExistsIntArray {
		t.Log("shouldn't have find that complex64 value")
		t.FailNow()
	}
}

func testString(t *testing.T) {
	// simple string exists
	var haystack = []string{"0", "1", "2", "3", "4", "5"}
	var needle string = "3"

	isExistsIntArray, _, _ := InArray(haystack, needle)
	if !isExistsIntArray {
		t.Log("can't find that string value")
		t.FailNow()
	}

	// simple int string exists
	haystack = []string{"0", "1", "2", "3", "4", "5"}
	needle = "6"

	isExistsIntArray, _, _ = InArray(haystack, needle)
	if isExistsIntArray {
		t.Log("shouldn't have find that string value")
		t.FailNow()
	}
}

func testUintptr(t *testing.T) {
	// simple uintptr exists
	var haystack = []uintptr{0, 1, 2, 3, 4, 5}
	var needle uintptr = 3

	isExistsIntArray, _, _ := InArray(haystack, needle)
	if !isExistsIntArray {
		t.Log("can't find that uintptr value")
		t.FailNow()
	}

	// simple int uintptr exists
	haystack = []uintptr{0, 1, 2, 3, 4, 5}
	needle = 6

	isExistsIntArray, _, _ = InArray(haystack, needle)
	if isExistsIntArray {
		t.Log("shouldn't have find that uintptr value")
		t.FailNow()
	}
}

func testStruct(t *testing.T) {
	var haystackStruct = []any{
		TestStruct{ID: 0, Value: "0"},
		TestStruct{ID: 1, Value: "1"},
		TestStruct{ID: 2, Value: "2"},
		TestStruct{ID: 3, Value: "3"},
		TestStruct{ID: 4, Value: "4"},
		TestStruct{ID: 5, Value: "5"},
	}

	needleExistsFunction := func(data any) bool {
		var testStruct TestStruct
		dataMarshal, _ := json.Marshal(data)
		_ = json.Unmarshal(dataMarshal, &testStruct)

		return testStruct.ID == 3
	}

	idxExists := slices.IndexFunc(haystackStruct, needleExistsFunction)
	isExistsStruct, idxExistsStruct, _ := InArrayComplex(haystackStruct, needleExistsFunction)
	if !isExistsStruct {
		t.Log("can't find that value")
		t.Fail()
	}

	if idxExistsStruct != idxExists {
		t.Log("index not the same")
		t.Fail()
	}

	needleNotExistsFunction := func(data any) bool {
		var testStruct TestStruct
		dataMarshal, _ := json.Marshal(data)
		_ = json.Unmarshal(dataMarshal, &testStruct)

		return testStruct.ID == 6
	}

	idxNotExists := slices.IndexFunc(haystackStruct, needleNotExistsFunction)
	isNotExistsStruct, idxNotExistsStruct, _ := InArrayComplex(haystackStruct, needleNotExistsFunction)
	if isNotExistsStruct {
		t.Log("shouldn't find the value in the struct")
		t.Fail()
	}

	if idxNotExistsStruct != -1 && idxNotExists != -1 {
		t.Log("idx shouldn't be find")
		t.Fail()
	}
}

func TestGeneratePassword(t *testing.T) {
	t.Run("testing weak password", testWeakPassword)
	t.Run("testing strong password", testStrongPassword)
}

func testWeakPassword(t *testing.T) {
	password := Passwords{
		Username: "test",
		Email:    "test@gmail.com",
		Password: "test",
	}

	_, err := GeneratePassword(password, 15)
	if err != nil && err.Error() != "your password not good enough, please try again" {
		t.Log("password is strong, it should be weak")
		t.Fail()
	}

	password = Passwords{
		Username: "test",
		Email:    "test@gmail.com",
		Password: "B2g0!KKG4yC9",
	}

	_, err = GeneratePassword(password, 10)
	if err != nil && err.Error() != "password cost need to be more or equal than 15" {
		t.Log("password cost more than 15, it should be less than 15")
		t.Fail()
	}
}

func testStrongPassword(t *testing.T) {
	password := Passwords{
		Username: "test",
		Email:    "test@gmail.com",
		Password: "B2g0!KKG4yC9",
	}

	_, err := GeneratePassword(password, 15)
	if err != nil {
		t.Log("password is weak, it should be strong")
		t.Fail()
	}
}

const unixTime int64 = 1661009671 // 20 Aug 22 22:34:31
const formatDate string = "02 Jan 06 15:04:05"

func TestAddTime(t *testing.T) {
	log.Printf("unix time: %d", unixTime)
	log.Printf("format date: %s", formatDate)
	log.Printf("date: %s", time.Unix(unixTime, 0).Format(formatDate))

	t.Run("test adding/subtracting time -second-:", testAddTimeSecond)
	t.Run("test adding/subtracting time -minutes-:", testAddTimeMinutes)
	t.Run("test adding/subtracting time -hours-:", testAddTimeHours)
	t.Run("test adding/subtracting time -days-:", testAddTimeDays)
	t.Run("test adding/subtracting time -weeks-:", testAddTimeWeeks)
	t.Run("test adding/subtracting time -months-:", testAddTimeMonths)
	t.Run("test adding/subtracting time -years-:", testAddTimeYears)
}

func testAddTimeSecond(t *testing.T) {
	currentTime := time.Unix(unixTime, 0)
	resultTime, _ := AddTime(currentTime, 10, "seconds") // 1661009681 / 20 Aug 22 22:34:41
	if resultTime.Unix() != 1661009681 {
		t.Log("failed to add 10 seconds")
		t.Fail()
	}

	currentTime = time.Unix(unixTime, 0)
	resultTime, _ = AddTime(currentTime, -10, "seconds") // 1661009661 / 20 Aug 22 22:34:21
	if resultTime.Unix() != 1661009661 {
		t.Log("failed to sub 10 seconds")
		t.Fail()
	}
}

func testAddTimeMinutes(t *testing.T) {
	currentTime := time.Unix(unixTime, 0)
	resultTime, _ := AddTime(currentTime, 10, "minutes") // 1661010271 / 20 Aug 22 22:44:31
	if resultTime.Unix() != 1661010271 {
		t.Log("failed to add 10 minutes")
		t.Fail()
	}

	currentTime = time.Unix(unixTime, 0)
	resultTime, _ = AddTime(currentTime, -10, "minutes") // 1661009071 / 20 Aug 22 22:24:31
	if resultTime.Unix() != 1661009071 {
		t.Log("failed to sub 10 minutes")
		t.Fail()
	}
}

func testAddTimeHours(t *testing.T) {
	currentTime := time.Unix(unixTime, 0)
	resultTime, _ := AddTime(currentTime, 1, "hours") // 1661013271 / 20 Aug 22 23:34:31
	if resultTime.Unix() != 1661013271 {
		t.Log("failed to add 1 hours")
		t.Fail()
	}

	currentTime = time.Unix(unixTime, 0)
	resultTime, _ = AddTime(currentTime, -1, "hours") // 1661006071 / 20 Aug 22 21:34:31
	if resultTime.Unix() != 1661006071 {
		t.Log("failed to sub 1 hours")
		t.Fail()
	}
}

func testAddTimeDays(t *testing.T) {
	currentTime := time.Unix(unixTime, 0)
	resultTime, _ := AddTime(currentTime, 1, "days") // 1661096071 / 21 Aug 22 22:34:31
	if resultTime.Unix() != 1661096071 {
		t.Log("failed to add 1 days")
		t.Fail()
	}

	currentTime = time.Unix(unixTime, 0)
	resultTime, _ = AddTime(currentTime, -1, "days") // 1660923271 / 19 Aug 22 22:34:31
	if resultTime.Unix() != 1660923271 {
		t.Log("failed to sub 1 days")
		t.Fail()
	}
}

func testAddTimeWeeks(t *testing.T) {
	currentTime := time.Unix(unixTime, 0)
	resultTime, _ := AddTime(currentTime, 1, "weeks") // 1661614471 / 27 Aug 22 22:34:31
	if resultTime.Unix() != 1661614471 {
		t.Log("failed to add 1 weeks")
		t.Fail()
	}

	currentTime = time.Unix(unixTime, 0)
	resultTime, _ = AddTime(currentTime, -1, "weeks") // 1660404871 / 13 Aug 22 22:34:31
	if resultTime.Unix() != 1660404871 {
		t.Log("failed to sub 1 weeks")
		t.Fail()
	}
}

func testAddTimeMonths(t *testing.T) {
	currentTime := time.Unix(unixTime, 0)
	resultTime, _ := AddTime(currentTime, 1, "months") // 1663688071 / 20 Sep 22 22:34:31
	if resultTime.Unix() != 1663688071 {
		t.Log("failed to sub 1 months")
		t.Fail()
	}

	currentTime = time.Unix(unixTime, 0)
	resultTime, _ = AddTime(currentTime, -1, "months") // 1658331271 / 20 Jul 22 22:34:31
	if resultTime.Unix() != 1658331271 {
		t.Log("failed to sub 1 months")
		t.Fail()
	}
}

func testAddTimeYears(t *testing.T) {
	currentTime := time.Unix(unixTime, 0)
	resultTime, _ := AddTime(currentTime, 1, "years") // 1692545671 / 20 Aug 23 22:34:31
	if resultTime.Unix() != 1692545671 {
		t.Log("failed to add 1 years")
		t.Fail()
	}

	currentTime = time.Unix(unixTime, 0)
	resultTime, _ = AddTime(currentTime, -1, "years") // 1629473671 / 20 Aug 21 22:34:31
	if resultTime.Unix() != 1629473671 {
		t.Log("failed to sub 1 years")
		t.Fail()
	}
}
