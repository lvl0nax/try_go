package main

import (
	"bufio"
	"bytes"
	json "encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_   *json.RawMessage
	_   *jlexer.Lexer
	_   *jwriter.Writer
	_   easyjson.Marshaler
	buf bytes.Buffer
)

//easyjson:json
type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

var userPool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]bool, 100)
	fileScanner := bufio.NewScanner(file)
	i := -1
	fmt.Fprintln(out, "found users:")
	isAnd := false
	isMS := false

	for fileScanner.Scan() {
		i++

		user := userPool.Get().(*User)
		user.UnmarshalJSON(fileScanner.Bytes())
		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			isAnd = strings.Contains(browser, "Android")
			isMS = strings.Contains(browser, "MSIE")

			if isAnd || isMS {
				isAndroid = isAndroid || isAnd
				isMSIE = isMSIE || isMS

				if !seenBrowsers[browser] {
					seenBrowsers[browser] = true
				}
			}
		}

		if isAndroid && isMSIE {
			buf.WriteString("[")
			buf.WriteString(strconv.Itoa(i))
			buf.WriteString("] ")
			buf.WriteString(user.Name)
			buf.WriteString(" <")
			buf.WriteString(strings.Replace(user.Email, "@", " [at] ", -1))
			buf.WriteString(">\n")
			out.Write(buf.Bytes())
			buf.Reset()
		}
		userPool.Put(user)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

func easyjson3486653aDecodeHw3Bench(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3486653aEncodeHw3Bench(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"browsers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3486653aEncodeHw3Bench(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3486653aEncodeHw3Bench(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3486653aDecodeHw3Bench(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3486653aDecodeHw3Bench(l, v)
}
