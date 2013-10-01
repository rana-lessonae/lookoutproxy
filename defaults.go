/*
lookoutproxy default values
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
import "time"


const MIN_INT32	int64 = - (1 << 31)
const MAX_INT32	int64 = (1 << 31) - 1
const MIN_INT64	int64 = - (1 << 63)
const MAX_INT64	int64 = (1 << 63) - 1

// Max number of concurrent running threads
var MAX_THREADS	int64 = int64(runtime.NumCPU())

// Queue sizes
const CONTROL_QUEUE_SIZE = 10
const MESSAGE_QUEUE_SIZE = 1000

// Control messages
const OBJECT_STOP =	"STOP"
const OBJECT_CONNECT =	"CONNECT!"

// Minimum and maximum wait time (ns)
const MIN_WAIT_TIME =	1 * time.Millisecond
const MAX_WAIT_TIME =	2 * time.Second

// Maximum datagram sizes
const MAX_IP_SIZE =	65535
const MAX_UDP_SIZE =	MAX_IP_SIZE - 20
const MAX_TCP_SIZE =	MAX_IP_SIZE - 20

// Minimum and maximum port numbers
const MIN_UDP_PORT =	1
const MAX_UDP_PORT =	65535
const MIN_TCP_PORT =	1
const MAX_TCP_PORT =	65535

