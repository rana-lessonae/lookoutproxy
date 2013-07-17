/*
lookoutproxy is an IT management proxy
Planned or implemented features include :
	- SNMP trap forwarding
	- TCP proxying
	- IPv4 and IPv6 support
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

import "runtime"
import "fmt"
//import "reflect"


func main() {

	// Should load configuration files here

	var value string
	var present bool

	value, present = globalConfiguration[GLOBAL_CONFIGURATION_MAX_THREADS]
	if present {
		var threads int
		var valid bool

		threads, valid = intValue(value)

		if valid {
			if threads == 0 {
				threads = MAX_THREADS
			} else if threads < 0 {
				threads = MAX_THREADS + threads
			}

			if threads < 0 {
				threads = 1
			} else if threads > MAX_THREADS {
				threads = MAX_THREADS
			}

			runtime.GOMAXPROCS(threads)
			fmt.Printf("%d usable threads.\n", threads)
		} else {
			fmt.Printf("Non integer value for %s global parameter: \"%s\" !\n", GLOBAL_CONFIGURATION_MAX_THREADS, value)
		}
	}

	for object, configuration := range objectsConfiguration {
		value, present = configuration[CONFIGURATION_LOAD_ON_STARTUP]
		if present {
			if isTrue(value) {
				fmt.Printf("%s should be started now.\n", object)
			} else if isNotFalse(value) {
				fmt.Printf("Should %s be started or not ?\n", object)
			}
		}
	}
}

