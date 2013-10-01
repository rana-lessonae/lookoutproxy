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

import "time"
import "runtime"
import "fmt"


func main() {
	// Should load configuration files here

	// Thread number configuration
	present, valid, threads := int64Config(globalConfiguration, GLOBAL_MAX_THREADS, 1 - MAX_THREADS, MAX_THREADS)
	if present {
		if valid {
			if threads <= 0 {
				threads = MAX_THREADS + threads
			}

			runtime.GOMAXPROCS(int(threads))
			fmt.Printf("%d usable thread(s).\n", threads)
		} else {
			fmt.Printf("Invalid value for %s global parameter !\n", GLOBAL_MAX_THREADS)
		}
	}

	// Check all objects in configuration
	objectControl := make(chan ControlMessage, CONTROL_QUEUE_SIZE)
	var invalidConfig	bool = false

	for objectName, configuration := range modulesConfiguration {
		present, valid, moduleDef := configuration.moduleDef(MODULE_NAME)
		if present && valid {
			object := moduleDef.create(objectName, configuration, objectControl)
			valid = moduleDef.check(object)
			if ! valid {
				fmt.Printf("Object \"%s\" configuration is invalid !\n", objectName)
				invalidConfig = true
			}
		} else {
			fmt.Printf("Object \"%s\" is not referencing a valid module !\n", objectName)
			invalidConfig = true
		}
	}

	// Stop on configuration error
	if invalidConfig {
		return
	}

	// Run objects where load on startup is active
	for objectName, configuration := range modulesConfiguration {
		present, valid, start := boolConfig(configuration, MODULE_LOAD_ON_STARTUP)
		if present {
			if valid {
				if start {
					_ = RunModule(objectName, MODULE_LOAD_ON_STARTUP, objectControl)
				}
			} else {
				fmt.Printf("Should \"%s\" be started or not ?\n", objectName)
			}
		}
	}

	time.Sleep(5 * time.Second)

	// Send stop message to child threads
	for _, object := range runningObjects {
		if object.startReason == MODULE_LOAD_ON_STARTUP {
			object.SendOrder(OBJECT_STOP)
		}
	}

	var control	ControlMessage

	// Wait for the child threads
	for ; len(runningObjects) > 0 ; {
		select {
			case control = <-objectControl:
				if control.text == OBJECT_STOP {
					control.from.CleanUp ()
				}
			default:
				time.Sleep(1 * time.Second)
		}
	}
}

