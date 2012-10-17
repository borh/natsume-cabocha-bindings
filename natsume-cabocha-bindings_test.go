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
package natsume_cabocha_bindings

import (
	"bytes"
	"testing"
)

func TestParseToLatticeString(t *testing.T) {
	output := ParseToLatticeString(input)
	if output != outputCorrect {
		t.Errorf("Echo: expected %q got %q", outputCorrect, output)
	}
}

func BenchmarkParseToLatticeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output := ParseToLatticeString(input)
		if output != outputCorrect {
			b.Errorf("Echo: expected %q got %q", outputCorrect, output)
		}
	}
}

func TestParseToSentenceJSON(t *testing.T) {
	output := ParseToSentence(input).ToJSON()
	if !bytes.Equal(output, outputCorrectJSON) {
		t.Errorf("Echo: expected %q got %q", outputCorrectJSON, output)
	}
}

func BenchmarkParseToSentenceJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output := ParseToSentence(input).ToJSON()
		if !bytes.Equal(output, outputCorrectJSON) {
			b.Errorf("Echo: expected %q got %q", outputCorrectJSON, output)
		}
	}
}

var input = "hello，未知語"
var outputCorrect = `* 0 -1D 3/3 0.000000
hello	名詞,普通名詞,一般,*,*,*	O
，	補助記号,読点,*,*,*,*,,，,，,,,記号,，,,,,*,*,*,*,*,*,*,*,*	O
未知	名詞,普通名詞,形状詞可能,*,*,*,ミチ,未知,未知,ミチ,ミチ,漢,未知,ミチ,ミチ,ミチ,*,*,*,*,*,*,1,C3,*	O
語	名詞,普通名詞,一般,*,*,*,ゴ,語,語,ゴ,ゴ,漢,語,ゴ,ゴ,ゴ,*,*,*,*,*,*,1,C3,*	O
EOS
`
var outputCorrectJSON = []byte("[\n  {\n    \"id\": 0,\n    \"link\": -1,\n    \"prob\": 0,\n    \"head\": 3,\n    \"tail\": 3,\n    \"tokens\": [\n      {\n        \"begin\": 0,\n        \"end\": 5,\n        \"pos1\": \"名詞\",\n        \"pos2\": \"普通名詞\",\n        \"pos3\": \"一般\",\n        \"pos4\": \"*\",\n        \"cType\": \"*\",\n        \"cForm\": \"*\",\n        \"lForm\": \"\",\n        \"lemma\": \"hello\",\n        \"orth\": \"hello\",\n        \"pron\": \"\",\n        \"kana\": \"\",\n        \"goshu\": \"不明\",\n        \"orthBase\": \"hello\",\n        \"pronBase\": \"\",\n        \"kanaBase\": \"\",\n        \"formBase\": \"\",\n        \"iType\": \"\",\n        \"iForm\": \"\",\n        \"IConType\": \"\",\n        \"fType\": \"\",\n        \"fForm\": \"\",\n        \"fConType\": \"\",\n        \"aType\": \"\",\n        \"aConType\": \"\",\n        \"aModType\": \"\",\n        \"ne\": \"O\"\n      },\n      {\n        \"begin\": 5,\n        \"end\": 6,\n        \"pos1\": \"補助記号\",\n        \"pos2\": \"読点\",\n        \"pos3\": \"*\",\n        \"pos4\": \"*\",\n        \"cType\": \"*\",\n        \"cForm\": \"*\",\n        \"lForm\": \"\",\n        \"lemma\": \"，\",\n        \"orth\": \"，\",\n        \"pron\": \"\",\n        \"kana\": \"\",\n        \"goshu\": \"記号\",\n        \"orthBase\": \"，\",\n        \"pronBase\": \"\",\n        \"kanaBase\": \"\",\n        \"formBase\": \"\",\n        \"iType\": \"*\",\n        \"iForm\": \"*\",\n        \"IConType\": \"*\",\n        \"fType\": \"*\",\n        \"fForm\": \"*\",\n        \"fConType\": \"*\",\n        \"aType\": \"*\",\n        \"aConType\": \"*\",\n        \"aModType\": \"*\",\n        \"ne\": \"O\"\n      },\n      {\n        \"begin\": 6,\n        \"end\": 8,\n        \"pos1\": \"名詞\",\n        \"pos2\": \"普通名詞\",\n        \"pos3\": \"形状詞可能\",\n        \"pos4\": \"*\",\n        \"cType\": \"*\",\n        \"cForm\": \"*\",\n        \"lForm\": \"ミチ\",\n        \"lemma\": \"未知\",\n        \"orth\": \"未知\",\n        \"pron\": \"ミチ\",\n        \"kana\": \"ミチ\",\n        \"goshu\": \"漢\",\n        \"orthBase\": \"未知\",\n        \"pronBase\": \"ミチ\",\n        \"kanaBase\": \"ミチ\",\n        \"formBase\": \"ミチ\",\n        \"iType\": \"*\",\n        \"iForm\": \"*\",\n        \"IConType\": \"*\",\n        \"fType\": \"*\",\n        \"fForm\": \"*\",\n        \"fConType\": \"*\",\n        \"aType\": \"1\",\n        \"aConType\": \"C3\",\n        \"aModType\": \"*\",\n        \"ne\": \"O\"\n      },\n      {\n        \"begin\": 8,\n        \"end\": 9,\n        \"pos1\": \"名詞\",\n        \"pos2\": \"普通名詞\",\n        \"pos3\": \"一般\",\n        \"pos4\": \"*\",\n        \"cType\": \"*\",\n        \"cForm\": \"*\",\n        \"lForm\": \"ゴ\",\n        \"lemma\": \"語\",\n        \"orth\": \"語\",\n        \"pron\": \"ゴ\",\n        \"kana\": \"ゴ\",\n        \"goshu\": \"漢\",\n        \"orthBase\": \"語\",\n        \"pronBase\": \"ゴ\",\n        \"kanaBase\": \"ゴ\",\n        \"formBase\": \"ゴ\",\n        \"iType\": \"*\",\n        \"iForm\": \"*\",\n        \"IConType\": \"*\",\n        \"fType\": \"*\",\n        \"fForm\": \"*\",\n        \"fConType\": \"*\",\n        \"aType\": \"1\",\n        \"aConType\": \"C3\",\n        \"aModType\": \"*\",\n        \"ne\": \"O\"\n      }\n    ]\n  }\n]")
