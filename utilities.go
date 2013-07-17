/*
lookoutproxy utilities
*/


/*
   Copyright 2013 Rana Lessonae

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/


package main

import "strings"
import "strconv"


func isTrue(value string) bool {
	switch strings.ToUpper(value) {
		case "1":	return true
		case "YES":	return true
		case "TRUE":	return true
		default:	return false
	}
}


func isFalse(value string) bool {
	switch strings.ToUpper(value) {
		case "0":	return true
		case "NO":	return true
		case "FALSE":	return true
		default:	return false
	}
}


func isNotTrue(value string) bool {
	return ! isTrue(value)
}


func isNotFalse(value string) bool {
	return ! isFalse(value)
}


func intValue(value string) (result int, valid bool) {
	var number int64
	var err error

	number, err = strconv.ParseInt(value, 0, 0)

	if err == nil {
		return int(number), true
	} else {
		return 0, false
	}
}

