/*
Copyright (c) 2012 Bor Hodošček. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	c "../"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func bodyReadHelper(w http.ResponseWriter, r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("HTTP POST Body error!", err)
	}
	return string(body)
}

// TODO most likely needs some more work and testing
func websocketHandler(ws *websocket.Conn) {
	for {
		var input string
		err := websocket.Message.Receive(ws, &input)
		if err != nil {
			fmt.Println(err)
			break
		}
		//fmt.Println("Received from client: " + input)

		reply := c.ParseToFormat(c.Parser, input, c.FormatLattice)
		err = websocket.Message.Send(ws, reply)
		if err != nil {
			fmt.Println(err)
			break
		}
		//fmt.Println("Sent: " + reply)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", c.ParseToLatticeString(bodyReadHelper(w, r)))
	})
	http.HandleFunc("/xml", func(w http.ResponseWriter, r *http.Request) {
		w.Write(c.ParseToSentence(bodyReadHelper(w, r)).ToXML())
	})
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Write(c.ParseToSentence(bodyReadHelper(w, r)).ToJSON())
	})
	http.Handle("/ws", websocket.Handler(websocketHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
