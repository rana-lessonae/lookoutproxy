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
import "time"


func isTrue(value string) bool {
	switch strings.ToLower(value) {
		case "1":	return true
		case "yes":	return true
		case "true":	return true
		default:	return false
	}
}


func isFalse(value string) bool {
	switch strings.ToLower(value) {
		case "0":	return true
		case "no":	return true
		case "false":	return true
		default:	return false
	}
}


func isNotTrue(value string) bool {
	return ! isTrue(value)
}


func isNotFalse(value string) bool {
	return ! isFalse(value)
}


func boolConfig(configuration ObjectConfiguration, index string) (present, valid, value bool) {
	var text string

	text, present = configuration[index]

	if present {
		if isTrue (text) {
			valid = true
			value = true
		} else if isFalse (text) {
			valid = true
			value = false
		}
	}

	return
}


func int64Config(configuration ObjectConfiguration, index string, min, max int64) (present, valid bool, value int64) {
	var text string
	var problem error

	value = max
	text, present = configuration[index]

	if present {
		value, problem = strconv.ParseInt(text, 0, 64)

		if (problem == nil) && (min <= value) && (value <= max) {
			valid = true
		}
	}

	return
}


func durationConfig(configuration ObjectConfiguration, index string, min, max time.Duration) (present, valid bool, value time.Duration) {
	var text string
	var problem error

	value = max
	text, present = configuration[index]

	if present {
		value, problem = time.ParseDuration(text)

		if (problem == nil) && (min <= value) && (value <= max) {
			valid = true
		}
	}

	return
}


/*func funcConfig(configuration ObjectConfiguration, index string) (present, valid bool, function ObjectCreation) {
	var text string

	text, present = configuration[index]

	if present {
		function, valid = objectCreation[text]
	}

	return
}*/


func moduleConfig(configuration ObjectConfiguration, index string) (present, valid bool, module ObjectConfiguration) {
	var text string

	text, present = configuration[index]

	if present {
		module, valid = modulesConfiguration[text]
	}

	return
}

