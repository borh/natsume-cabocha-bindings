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
package natsume_cabocha_binding

// #cgo LDFLAGS: -lcabocha
// #include <stdio.h>
// #include <stdlib.h>
// #include <cabocha.h>
// struct cabocha_t {};
import "C"

import (
	"encoding/csv"
	"encoding/xml"
	"log"
	re "regexp"
	"strconv"
	"strings"
	"encoding/json"
)

// Morpheme struct containing named strings for UniDic features.
type Token struct {
	Id       int               `xml:"id,attr" json:"id"`
	Features map[string]string `xml:"features" json:"features"`
	Ne       string            `xml:"ne,attr" json:"ne"`
}

type Chunk struct {
	Id     int64    `xml:"id,attr" json:"id"`
	Link   int64    `xml:"link,attr" json:"link"`
	Prob   float64  `xml:"score,attr" json:"prob"`
	Head   int64    `xml:"head,attr" json:"head"`
	Tail   int64    `xml:"func,attr" json:"tail"`
	Tokens []*Token `json:"tokens"`
}

// Sentence struct type wrapper for slice of Morpheme structs.
type Sentence struct {
	Chunks []*Chunk `json:"chunks"`
}

type TokenXML struct {
	XMLName  xml.Name `xml:"tok"`
	Id       int      `xml:"id,attr"`
	Features string   `xml:"feature,attr"`
	Ne       string   `xml:"ne,attr"`
	Orth     string   `xml:",chardata"`
	// Orth is included in Features.
}

type ChunkXML struct {
	XMLName xml.Name    `xml:"chunk"`
	Id      int64       `xml:"id,attr"`
	Link    int64       `xml:"link,attr"`
	Rel     string      `xml:"rel,attr"`
	Prob    float64     `xml:"score,attr"`
	Head    int64       `xml:"head,attr"`
	Tail    int64       `xml:"func,attr"`
	Tokens  []*TokenXML // `xml:"tok"`
}

// Sentence struct type wrapper for slice of Morpheme structs.
type SentenceXML struct {
	XMLName xml.Name `xml:"sentence"`
	//Chunks  []*Chunk// `xml:"chunk"`
	Data string `xml:",innerxml"`
}

func NewChunk(s string) *Chunk {
	fields := strings.Split(s, " ")
	head_tail := strings.Split(fields[3], "/")
	c := new(Chunk)
	c.Id, _ = strconv.ParseInt(fields[1], 10, 64)
	c.Link, _ = strconv.ParseInt(strings.Replace(fields[2], "D", "", 1), 10, 64)
	c.Head, _ = strconv.ParseInt(head_tail[0], 10, 64)
	c.Tail, _ = strconv.ParseInt(head_tail[1], 10, 64)
	c.Prob, _ = strconv.ParseFloat(fields[4], 64)
	return c
}

var chunkHeaderRe = re.MustCompile(`^\*[^\t\*]+$`)

var featureMap = map[int]string{
	0:  "pos1",
	1:  "pos2",
	2:  "pos3",
	3:  "pos4",
	4:  "cType",
	5:  "cForm",
	6:  "lForm",
	7:  "lemma",
	8:  "orth",
	9:  "pron",
	10: "kana",
	11: "goshu",
	12: "orthBase",
	13: "pronBase",
	14: "kanaBase",
	15: "formBase",
	16: "iType",
	17: "iForm",
	18: "iConType",
	19: "fType",
	20: "fForm",
	21: "fConType",
	22: "aType",
	23: "aConType",
	24: "aModType",
}

// Takes the CaboCha output of one sentence as a string and returns a pointer to the corresponding Sentence struct.
// CaboCha output should comprise one (un-split) sentence only.
func NewSentence(cabocha_out string) *Sentence {
	mecab_lines := strings.Split(cabocha_out, "\n")
	s := new(Sentence)
	c := new(Chunk)
	i := 0
	for _, line := range mecab_lines {
		if chunkHeaderRe.MatchString(line) { // New chunk
			if len(c.Tokens) != 0 {
				s.Chunks = append(s.Chunks, c)
			}
			c = NewChunk(line)
			continue
		} else if line == "EOS" || line == "" { // End
			s.Chunks = append(s.Chunks, c)
			break // TODO NewSentence must be called on only one sentence
		}

		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			log.Println("Error decoding token features:", fields)
			continue
		}

		csvReader := csv.NewReader(strings.NewReader(fields[1]))
		featuresSlice, err := csvReader.Read()
		if err != nil {
			log.Println("Error decoding feature csv field:", err)
		}
		features := make(map[string]string)
		for j, el := range featuresSlice {
			features[featureMap[j]] = el
		}

		t := &Token{
			Id:       i,
			Features: features,
			Ne:       fields[2],
		}
		c.Tokens = append(c.Tokens, t)

		i += 1
	}
	return s
}

func NewParser(opt string) *C.cabocha_t {
	return C.cabocha_new2(C.CString(opt))
}

const (
	FormatTree = iota
	FormatLattice
	FormatTreeLatice
	FormatXml
	FormatNone
)

func ParseToFormat(cabo *C.cabocha_t, s string, format _Ctype_int) string {
	tree := C.cabocha_sparse_totree(cabo, C.CString(s))
	return C.GoString(C.cabocha_tree_tostr(tree, format))
}

// Convenience function that returns the CaboCha output as a Sentence
// struct.
func ParseToSentence(s string) *Sentence {
	return NewSentence(ParseToFormat(Parser, s, FormatLattice))
}

// Convenience function that returns the CaboCha output as a lattice
// formatted string.
func ParseToLatticeString(s string) string {
	return ParseToFormat(Parser, s, FormatLattice)
}

func (s Sentence) ToJSON() []byte {
	jsonSentence, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println(err)
	}
	return []byte(jsonSentence)
}

func (s Sentence) ToXML() []byte {
	xmlSentence, err := xml.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println(err)
	}
	return []byte(xmlSentence)
}

// better way to instantiate only once?
var Parser = NewParser("-d /usr/lib64/mecab/dic/unidic -b /usr/lib64/mecab/dic/unidic/dicrc -r /etc/cabocharc-unidic -P UNIDIC")