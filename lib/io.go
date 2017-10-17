package lib

import (
	"os"
	"io"
	"bufio"
	"fmt"
)

// CopyLines : copy [start, end) lines from src to dest
// @Param : index start from 0
// Warning : it does not deal well with lines longer than 65536 characters
func CopyLines(src io.Reader, dest io.Writer, start, end int) error {
	scanner := bufio.NewScanner(src)

	for i := 0; i < start; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("Index Error: start_page (%v) greater than total pages (%v), no output written",
				start, i);
		}
	}

	writer := bufio.NewWriter(dest)
    for ; start < end; start++ {
		if (!scanner.Scan()) {
			if err := writer.Flush(); err != nil {
				panic(err)
			}
			return fmt.Errorf("Index Error: end_page (%v) greater than total pages (%v), less output than expected",
				end, start)
		}
		fmt.Fprintln(writer, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}
	return nil
}

// CopyPages : copy [start, end) pages from src to dest
// @Param : index start from 0
func CopyPages(src io.Reader, dest io.Writer, start, end int) error {
	const delim byte = '\f'

	reader := bufio.NewReader(src)
	for i := 0; i < start; i++ {
		if _, err := reader.ReadBytes(delim); err != nil {
			if err == io.EOF {
				return fmt.Errorf("Index Error: start_page (%v) greater than total pages (%v), no output written",
					start, i);
			}
			fmt.Fprintf(os.Stderr, "Read file error.\n")
			panic(err)
		}
	}

	writer := bufio.NewWriter(dest)
    for start < end {
		page, err := reader.ReadBytes(delim)
		if err != nil {
			if err == io.EOF {
				writer.Write(page[:])
				if err := writer.Flush(); err != nil {
					panic(err)
				}
				return fmt.Errorf("Index Error: end_page (%v) greater than total pages (%v), less output than expected",
					end, start)
			}
			fmt.Fprintf(os.Stderr, "Read file error.\n")
			panic(err)
		}
		if nn, err := writer.Write(page[:len(page) - 1]); err != nil {
			fmt.Fprintf(os.Stderr, "Expect write %v bytes but %v bytes.\n", len(page) - 1, nn)
			panic(err)
		}
		if start++; start != end {
			if err := writer.WriteByte(delim); err != nil {
				fmt.Fprintf(os.Stderr, "Expect write %q but fail.\n", delim)
				panic(err)
			}
		}
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}
	return nil
}
