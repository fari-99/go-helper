package gohelper

import (
	"fmt"
	"reflect"
	"time"

	"github.com/nbutton23/zxcvbn-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

// InArray checks whether needle is in haystack.
// if haystack is an array of struck, then needle need to be a function
// example : https://stackoverflow.com/questions/38654383/how-to-search-for-an-element-in-a-golang-slice
func InArray(needle any, haystack []any) (bool, int, error) {
	haystackValue := reflect.ValueOf(haystack)
	checkData := haystackValue.Index(0)
	checkDataKind := reflect.ValueOf(checkData).Kind()
	needleValue := reflect.ValueOf(needle)
	needleKind := needleValue.Kind()

	// haystack is empty array
	if !checkData.IsValid() {
		return false, -1, nil
	}

	isStruct := false
	if checkDataKind == reflect.Struct && needleKind == reflect.Func {
		isStruct = true
	}

	if !isStruct {
		for i := 0; i < haystackValue.Len(); i++ {
			hayVal := haystackValue.Index(i).Interface()

			if reflect.DeepEqual(hayVal, needle) {
				return true, i, nil
			}
		}
	} else {
		if needleValue.Kind() != reflect.Func {
			return false, -1, fmt.Errorf("must be a function")
		}

		index := slices.IndexFunc(haystack, needle.(func(any) bool))
		if index != -1 {
			return true, index, nil
		}
	}

	return false, -1, nil
}

type Passwords struct {
	Email    string
	Username string
	Password string
}

// GeneratePassword generates bcrypt hash string of the given plaintext password
// because generate password is by using bcrypt, to check hash and password,
// please use bycrypt CompareHashAndPassword
func GeneratePassword(password Passwords, passwordCost int8) (hash *string, err error) {
	userInput := []string{password.Username, password.Email}
	secret := password.Password

	checkPassword := zxcvbn.PasswordStrength(secret, userInput)
	if checkPassword.Score <= 2 {
		return nil, fmt.Errorf("your password not good enough, please try again")
	}

	if passwordCost < 15 { // minimum cost is 15 (default bycrypt minimum is 4)
		return nil, fmt.Errorf("password cost need to be more or equal than 15")
	}

	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(secret), int(passwordCost))
	if err != nil {
		return nil, err
	}

	hashedPassword := string(hashedPasswordByte)
	return &hashedPassword, nil
}

const timeSeconds = "seconds"
const timeMinutes = "minutes"
const timeHours = "hours"
const timeDays = "days"
const timeWeeks = "weeks"
const timeMonths = "months"
const timeYears = "years"

// AddTime adding additional time to your time.
func AddTime(yourTime time.Time, addedTime int64, timeType string) (resultTime time.Time, err error) {
	switch timeType {
	case timeSeconds:
		resultTime = yourTime.Add(time.Second * time.Duration(addedTime))
	case timeMinutes:
		resultTime = yourTime.Add(time.Minute * time.Duration(addedTime))
	case timeHours:
		resultTime = yourTime.Add(time.Hour * time.Duration(addedTime))
	case timeDays:
		resultTime = yourTime.AddDate(0, 0, int(addedTime))
	case timeWeeks:
		addedTime = addedTime * 7 // 1 week = 7 days
		resultTime = yourTime.AddDate(0, 0, int(addedTime))
	case timeMonths:
		resultTime = yourTime.AddDate(0, int(addedTime), 0)
	case timeYears:
		resultTime = yourTime.AddDate(int(addedTime), 0, 0)
	default:
		err = fmt.Errorf("time type is not supported, please pick (seconds, minutes, hours, days, months, years)")
	}

	return resultTime, err
}
