/*
lookoutproxy out_udp output UDP module
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
import "time"
import "strings"
import "net"


func init() {
	moduleDefinition["out_udp"] = ModuleDefinition {
		inputs:	[]string{"input"},
		outputs:	[]string{},
		create:	Create_out_udp,
		check:	Check_out_udp,
		run:	Run_out_udp,
	}
}


func Create_out_udp(name string, staticConfiguration ObjectConfiguration, control chan ControlMessage) (object Object) {
	object.Init("out_udp", name, staticConfiguration, control)
	return
}


func Check_out_udp(object Object) (fully_valid bool) {
	fully_valid = true

	// Check wait time
	present, valid, _ := durationConfig(object.configuration, PACKET_TIMEOUT, MIN_WAIT_TIME, MAX_WAIT_TIME)
	if present && ! valid {
		fmt.Printf("\"%s\" does not have a valid \"%s\" value: %s\n", object.name, PACKET_TIMEOUT, object.configuration[PACKET_TIMEOUT])
		fully_valid = false
	}

	// Check IP version
	version, present := object.configuration[IP_VERSION]
	if present && (version != "") && (version != "4") && (version != "6") {
		fmt.Printf("\"%s\" does not have a valid \"%s\" value: %s\n", object.name, IP_VERSION, object.configuration[IP_VERSION])
		fully_valid = false
	}

	// Check destination address presence but validity is not checked due to potential long name resolution time
	_, present = object.configuration[DESTINATION_ADDRESS]
	if ! present {
		fmt.Printf("\"%s\" does not have a \"%s\" parameter !\n", object.name, DESTINATION_ADDRESS)
		fully_valid = false
	}

	// Check destination port number
	present, valid, _ = int64Config(object.configuration, DESTINATION_PORT, MIN_UDP_PORT, MAX_UDP_PORT)
	if ! present {
		fmt.Printf("\"%s\" does not have a \"%s\" parameter !\n", object.name, DESTINATION_PORT)
		fully_valid = false
	} else if ! valid {
		fmt.Printf("\"%s\" does not have a valid \"%s\" value: %s\n", object.name, DESTINATION_PORT, object.configuration[DESTINATION_PORT])
		fully_valid = false
	}

	fmt.Printf("\"%s\" checked.\n", object.name)
	return
}


func Run_out_udp(name string, object Object) () {
	defer func () { object.SendInfo(OBJECT_STOP) } ()

	// Define timeout value on socket writes
	present, valid, delay := durationConfig(object.configuration, PACKET_TIMEOUT, MIN_WAIT_TIME, MAX_WAIT_TIME)
	if ! (present && valid) {
		delay = MAX_WAIT_TIME
	}

	// Encode protocol and IP version
	protocol, present := object.configuration[IP_VERSION]
	if ! present {
		protocol = ""
	}
	protocol = "udp" + protocol

	// Encode address and port
	destination := object.configuration[DESTINATION_ADDRESS]
	if strings.Contains(destination, ":") {
		destination = "[" + destination + "]"
	}
	destination += ":" + object.configuration[DESTINATION_PORT]
	destinationAddress, problem := net.ResolveUDPAddr(protocol, destination)
	if problem != nil {
		fmt.Printf("\"%s\" cannot resolve destination \"%s\" !\n", object.name, destination)
		return
	}

	// Open port
	connection, problem := net.ListenPacket(protocol, ":0")
	if problem != nil {
		fmt.Printf("\"%s\" cannot open \"%s\" channel !\n", object.name, protocol)
		return
	}
	defer connection.Close()

	input := object.inputs["input"]

	for {
		select {
			case controlMessage := <-object.control.in :
				if controlMessage.text == OBJECT_STOP {
					return
				}
			case dataMessage := <-input :
				problem := connection.SetWriteDeadline(time.Now().Add(delay))
				if problem != nil {
					return
				}

				length, problem := connection.WriteTo(dataMessage.binaryData, destinationAddress)
				if problem == nil {
					fmt.Printf("\"%s\" sent %d bytes to \"%s\".\n", object.name, length, destination)
				} else {
					fmt.Printf("\"%s\" received error %v while sending to \"%s\".\n", object.name, problem, destination)
				}
		}
	}
}

