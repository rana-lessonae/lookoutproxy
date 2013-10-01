/*
lookoutproxy modules management
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

import "fmt"
import "strings"


type ObjectCreation func(name string, staticConfiguration ObjectConfiguration, control chan ControlMessage) (object Object)
type ObjectCheck func(object Object) (fully_valid bool)
type ObjectRun func(name string, object Object) ()

type ControlMessage struct {
	from	*Object
	text	string
}

type Control struct {
	in, out	chan ControlMessage
}

type Message struct {
	textData	string
	binaryData	[]byte
}

type Channel chan Message

// Module definition record and table
type ModuleDefinition struct {
	inputs	[]string
	outputs	[]string
	create	ObjectCreation
	check	ObjectCheck
	run	ObjectRun
}

var moduleDefinition = map[string]ModuleDefinition {
}


type Object struct {
	module	ModuleDefinition
	name	string
	startReason	string
	referenced	uint32
	referencing	[]*Object
	configuration	ObjectConfiguration
	control	Control
	inputs	map[string]Channel
	outputs	map[string]map[string]Channel
}


// Running object list
var runningObjects = make(map[string]*Object)


func (configuration ObjectConfiguration) moduleDef(index string) (present, valid bool, moduleDef ModuleDefinition) {
	var text string

	text, present = configuration[index]

	if present {
		moduleDef, valid = moduleDefinition[text]
	}

	return
}


func (object *Object) Init(module, name string, staticConfiguration ObjectConfiguration, control chan ControlMessage) () {
	object.module = moduleDefinition[module]
	object.name = name
	object.configuration = staticConfiguration

	object.control.in = make(chan ControlMessage, CONTROL_QUEUE_SIZE)
	object.control.out = control

	object.inputs = make(map[string]Channel)
	for _, input := range object.module.inputs {
		object.inputs[input] = make(Channel, MESSAGE_QUEUE_SIZE)
	}

	object.outputs = make(map[string]map[string]Channel)
	for _, output := range object.module.outputs {
		object.outputs[output] = make(map[string]Channel)
	}
}


func (object *Object) Run(control chan ControlMessage) () {
	runningObjects[object.name] = object

	outputList, present := object.configuration[MODULE_OUTPUTS]
	if present {
		for _, outputDef := range (strings.Split(outputList, ",")) {
			var problem	bool

			problem = false
			fields := strings.Split(outputDef, ":")
			if len(fields) != 3 {
				problem = true
			} else {
				var outputObject	*Object

				output, module, input := fields[0], fields[1], fields[2] 
				fmt.Printf("\"%s\" --> \"%s:%s\"\n", output, module, input)

				if ! isRunning(module) {
					outputObject = RunModule(module, MODULE_LOAD_ON_DEMAND, control)
				} else {
					for _, runningObject := range runningObjects {
						if runningObject.name == module {
							outputObject = runningObject
							break
						}
					}
				}

				outputList, validOutput := object.outputs[output]
				inputChannel, validInput := outputObject.inputs[input]

				if validOutput && validInput {
					outputList[module + ":" + input] = inputChannel
					outputObject.referenced += 1
					object.referencing = append(object.referencing, outputObject)
				} else {
					problem = true
				}
			}

			if problem {
				fmt.Printf("\"%s\" has an invalid output definition: \"%s\" !\n", object.name, outputDef)
			}
		}
	}

	go object.module.run (object.name, *object)
	return
}


func RunModule(module, startReason string, control chan ControlMessage) (object *Object) {
	fmt.Printf("\"%s\" should be started now.\n", module)
	configuration, present := modulesConfiguration[module]
	if ! present {
		fmt.Printf("\"%s\" has no configuration !\n", module)
		return
	}

	present, valid, moduleDef := configuration.moduleDef(MODULE_NAME)
	if present && valid {
		fmt.Printf("creating \"%s\" ...\n", module)
		newObject := moduleDef.create(module, configuration, control)
		object = &newObject
		object.startReason = startReason
		fmt.Printf("starting \"%s\" ...\n", module)
		object.Run(control)
	} else {
		fmt.Printf("\"%s\" has no valid associated function !\n", module)
	}

	return
}


func (object *Object) CleanUp () () {
	if object.referenced != 0 {
		fmt.Printf("\"%s\" is being cleaned up with %d connections to its input(s) !\n", object.name, object.referenced)
	}

	for _, reference := range object.referencing {
		reference.referenced -= 1

		if reference.referenced == 0 && reference.startReason != MODULE_LOAD_ON_STARTUP {
			reference.SendOrder(OBJECT_STOP)
		}
	}

	delete (runningObjects, object.name)
	fmt.Printf("\"%s\" has stopped.\n", object.name)
	fmt.Printf("There are %d running objects.\n", len(runningObjects))
}


func isRunning (name string) (running bool) {
	for _, object := range runningObjects {
		if object.name == name {
			running = true
		}
	}

	return
}


func (object *Object) SendOrder (message string) () {
	object.control.in <- ControlMessage {
		from:	nil,
		text:	message,
	}
}


func (object *Object) SendInfo (message string) () {
	object.control.out <- ControlMessage {
		from:	object,
		text:	message,
	}
}

