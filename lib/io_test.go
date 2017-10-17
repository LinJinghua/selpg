package lib

import (
	"testing"
	"bytes"
	"bufio"
)

func TestCopyLines(t *testing.T) {
	cases := []struct {
		start, end int
		want []byte
	}{
		{0, 2, []byte("1\n2\n")},
		{5, 6, []byte("")},
		{2, 2000, []byte("\f3\n\f4\n")},
	}

	readerData := []byte("1\n2\n\f3\n\f4\n")
	
	for _, c := range cases {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		reader := bytes.NewReader(readerData)
		CopyLines(reader, writer, c.start, c.end)
		writer.Flush()
		if !bytes.Equal(buf.Bytes(), c.want) {
			t.Errorf("CopyLines(%q, writer, %v, %v) == %q, want %q", 
				string(readerData[:]), c.start, c.end, string(buf.Bytes()[:]), string(c.want[:]))
		}
	}
}

func TestCopyPages(t *testing.T) {
	cases := []struct {
		start, end int
		want []byte
	}{
		{0, 2, []byte("1\n2\n\f3\n")},
		{5, 6, []byte("")},
		{1, 2000, []byte("3\n\f4\n")},
	}

	readerData := []byte("1\n2\n\f3\n\f4\n")
	
	for _, c := range cases {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		reader := bytes.NewReader(readerData)
		CopyPages(reader, writer, c.start, c.end)
		if !bytes.Equal(buf.Bytes(), c.want) {
			t.Errorf("CopyPages(%q, writer, %v, %v) == %q, want %q", 
				string(readerData[:]),  c.start, c.end, string(buf.Bytes()[:]), string(c.want[:]))
		}
	}
}
