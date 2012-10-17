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
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var serverAddr string
var once sync.Once

func startServer() {
	http.Handle("/", websocket.Handler(websocketHandler))
	http.Handle("/json", websocket.Handler(websocketHandlerJSON))

	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()

	log.Print("Test WebSocket server listening on ", serverAddr)
}

func TestWebsocketHandler(t *testing.T) {
	once.Do(startServer)

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s", serverAddr), "", "http://localhost/")
	if err != nil {
		t.Fatal("dialing", err)
	}

	input := []byte("hello")
	if _, err := ws.Write(input); err != nil {
		t.Errorf("Write: %v", err)
	}
	var output = make([]byte, 512)
	n, err := ws.Read(output)
	if err != nil {
		t.Errorf("Read: %v", err)
	}
	output = output[0:n]

	actual_output := []byte(c.ParseToLatticeString(string(input)))

	if !bytes.Equal(output, actual_output) {
		t.Errorf("Echo: expected %q got %q", actual_output, output)
	}
}

func TestWebsocketHandlerJSON(t *testing.T) {
	once.Do(startServer)

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s/json", serverAddr), "", "http://localhost/")
	if err != nil {
		t.Fatal("dialing", err)
	}

	input := []byte("hello")
	if _, err := ws.Write(input); err != nil {
		t.Errorf("Write: %v", err)
	}
	var output = make([]byte, 1024)
	n, err := ws.Read(output)
	if err != nil {
		t.Errorf("Read: %v", err)
	}
	output = output[0:n]

	actual_output := []byte(c.ParseToSentence(string(input)).ToJSON())

	if !bytes.Equal(output, actual_output) {
		t.Errorf("Echo: expected %q got %q", actual_output, output)
	}
}

func BenchmarkWebsocketHandler(b *testing.B) {
	b.StopTimer()
	once.Do(startServer)
	ws, err := websocket.Dial(fmt.Sprintf("ws://%s", serverAddr), "", "http://localhost/")
	if err != nil {
		b.Fatal("dialing", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		var input string = "こんにちは、世界。"
		err := websocket.Message.Send(ws, input)
		if err != nil {
			fmt.Println(err)
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println(err)
			break
		}

		var actual_reply string = `* 0 1D 0/0 0.000000
こんにちは	感動詞,一般,*,*,*,*,コンニチハ,今日は,こんにちは,コンニチワ,コンニチハ,混,こんにちは,コンニチワ,コンニチハ,コンニチハ,*,*,*,*,*,*,5,*,*	O
、	補助記号,読点,*,*,*,*,,、,、,,,記号,、,,,,*,*,*,*,*,*,*,*,*	O
* 1 -1D 0/0 0.000000
世界	名詞,普通名詞,一般,*,*,*,セカイ,世界,世界,セカイ,セカイ,漢,世界,セカイ,セカイ,セカイ,*,*,*,*,*,*,"1,2",C1,*	O
。	補助記号,句点,*,*,*,*,,。,。,,,記号,。,,,,*,*,*,*,*,*,*,*,*	O
EOS
`

		if reply != actual_reply { // FIXME somehow these strings are different?!
			b.Errorf("Echo: expected\n%q\ngot\n%q", actual_reply, reply)
		}
	}
}
