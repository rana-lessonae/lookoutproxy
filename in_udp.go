/*
lookoutproxy in_udp input UDP module
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
	moduleDefinition["in_udp"] = ModuleDefinition {
		inputs:	[]string{},
		outputs:	[]string{"output"},
		create:	Create_in_udp,
		check:	Check_in_udp,
		run:	Run_in_udp,
	}
}


func Create_in_udp(name string, staticConfiguration ObjectConfiguration, control chan ControlMessage) (object Object) {
	object.Init("in_udp", name, staticConfiguration, control)
	return
}


func Check_in_udp(object Object) (fully_valid bool) {
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

	// Listening address presence and validity are not checked due to potential long name resolution time

	// Check listening port number without checking if the port is available and the current user has enough privilege to open it
	present, valid, _ = int64Config(object.configuration, LISTENING_PORT, MIN_UDP_PORT, MAX_UDP_PORT)
	if ! present {
		fmt.Printf("\"%s\" does not have a \"%s\" parameter !\n", object.name, LISTENING_PORT)
		fully_valid = false
	} else if ! valid {
		fmt.Printf("\"%s\" does not have a valid \"%s\" value: %s\n", object.name, LISTENING_PORT, object.configuration[LISTENING_PORT])
		fully_valid = false
	}

	fmt.Printf("\"%s\" checked.\n", object.name)
	return
}


func Run_in_udp(name string, object Object) () {
	defer func () { object.SendInfo(OBJECT_STOP) } ()

	// Define timeout value on socket reads
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
	destination, present := object.configuration[LISTENING_ADDRESS]
	if ! present {
		destination = ""
	} else if strings.Contains(destination, ":") {
		destination = "[" + destination + "]"
	}
	destination += ":" + object.configuration[LISTENING_PORT]

	// Open port
	connection, problem := net.ListenPacket(protocol, destination)
	if problem != nil {
		fmt.Printf("\"%s\" cannot open \"%s\" port at \"%s\".\n", object.name, protocol, destination)
		return
	}
	defer connection.Close()

	datagram := make([]byte, MAX_UDP_SIZE)
	outputList := object.outputs["output"]

	for {
		select {
			case controlMessage := <-object.control.in:
				if controlMessage.text == OBJECT_STOP {
					return
				}
			default:
				problem := connection.SetReadDeadline(time.Now().Add(delay))
				if problem != nil {
					return
				}

				length, source, problem := connection.ReadFrom(datagram)
				if problem == nil {
					fmt.Printf("\"%s\" received %d bytes from \"%s\".\n", object.name, length, source.String())
					for _, channel := range outputList {
						channel <- Message {
							binaryData:	datagram[:length],
						}
					}
				}
		}
	}
}

