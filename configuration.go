/*
lookoutproxy configuration management
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

//import "time"


// Configuration types
type ObjectConfiguration map[string]string
type ModulesConfiguration map[string]ObjectConfiguration

// Configuration field IDs
const GLOBAL_MAX_THREADS =	"max_threads"
const MODULE_NAME =	"module_name"
const MODULE_LOAD_ON_STARTUP =	"load_on_startup"
const MODULE_LOAD_ON_DEMAND =	"load_on_demand"
const MODULE_OUTPUTS =	"outputs"
const MODULE_INPUTS =	"inputs"
const IP_VERSION =	"ip_version"
const LISTENING_ADDRESS =	"listening_address"
const LISTENING_PORT =	"listening_port"
const PACKET_TIMEOUT =	"packet_timeout"
const DESTINATION_ADDRESS =	"destination_address"
const DESTINATION_PORT =	"destination_port"

// Global configuration
var globalConfiguration = ObjectConfiguration {
	GLOBAL_MAX_THREADS:	"-1",
}

// Modules configuration
var modulesConfiguration = ModulesConfiguration {
	"UDP server v4":	{
		MODULE_NAME:	"in_udp",
		MODULE_LOAD_ON_STARTUP:	"yes",
		MODULE_OUTPUTS:	"output:UDP client v4:input,output:UDP client v6:input",
		IP_VERSION:	"4",
		LISTENING_ADDRESS:	"",
		LISTENING_PORT:	"10162",
		PACKET_TIMEOUT:	"100ms",
	},
	"UDP server v6":	{
		MODULE_NAME:	"in_udp",
		MODULE_LOAD_ON_STARTUP:	"yes",
		MODULE_OUTPUTS:	"output:UDP client v4:input,output:UDP client v6:input",
		IP_VERSION:	"6",
		LISTENING_ADDRESS:	"",
		LISTENING_PORT:	"10162",
		PACKET_TIMEOUT:	"100ms",
	},
	"UDP client v4":	{
		MODULE_NAME:	"out_udp",
		MODULE_LOAD_ON_STARTUP:	"no",
		IP_VERSION:	"4",
		DESTINATION_ADDRESS:	"192.168.0.1",
		DESTINATION_PORT:	"162",
		PACKET_TIMEOUT:	"1s",
	},
	"UDP client v6":	{
		MODULE_NAME:	"out_udp",
		MODULE_LOAD_ON_STARTUP:	"no",
		IP_VERSION:	"6",
		DESTINATION_ADDRESS:	"::1",
		DESTINATION_PORT:	"162",
		PACKET_TIMEOUT:	"1s",
	},
}

