package test

import (
	"testing"
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"strconv"
)

var testFileName = "testInput.txt"
var testResultFileName = "testResult.txt"
var stdResultFileName = "stdResult.txt"
var stdCFileName = "selpg.c"
var selpgName = "selpg"
var gccName = "gcc"
var stdName = "std"
var gitName = "git"
var diffName = "diff"
var pipeExecFileName = "test"

func createTestFile() {
	file, err := os.Create(testFileName)
	if err != nil {
		fmt.Printf("Write file %q error:", testFileName)
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := 1; i <= 1000; i++  {
		if nn, err := writer.WriteString(strconv.Itoa(i)); err != nil {
			fmt.Fprintf(os.Stderr, "Expect write %v bytes but %v bytes.\n", len(strconv.Itoa(i)), nn)
            panic(err)
        }
		if err := writer.WriteByte('\n'); err != nil {
			fmt.Fprintf(os.Stderr, "Expect write \\n but 0 bytes.\n")
            panic(err)
        }
		if i % 10 == 0 {
			if err := writer.WriteByte('\f'); err != nil {
				fmt.Fprintf(os.Stderr, "Expect write \\f but 0 bytes.\n")
				panic(err)
			}
		}
	}
	if err = writer.Flush(); err != nil {
        panic(err)
    }
}

func writeToFile(fileName string, output []byte) {
	file, _ := os.Create(fileName)
	if _, err := file.Write(output); err != nil {
		panic(err)
	}
	defer file.Close()
}

func TestMain(t *testing.T) {

	if _, err := exec.LookPath(selpgName); err != nil {
		panic(err)
	}
	if _, err := exec.LookPath(gccName); err != nil {
		panic(err)
	}
	if _, err := exec.LookPath(gitName); err != nil {
		panic(err)
	}
	if _, err := exec.LookPath(pipeExecFileName); err != nil {
		panic(err)
	}
	gcc := exec.Command(gccName, "-o", stdName, stdCFileName)
	gcc.Start()

	createTestFile()

	if err := gcc.Wait(); err != nil {
		panic(err)
	}
	
	{
		selpgOut, _ := usePipe(selpgName, "-s", "1", "-e", "2", "-l", "10", testFileName)
		stdOut, _ := usePipe(stdName, "-s1", "-e2", "-l10", testFileName)

		writeToFile(testResultFileName, selpgOut)
		writeToFile(stdResultFileName, stdOut)

		diff := exec.Command(gitName, diffName, testResultFileName, stdResultFileName)
		if _, err := diff.Output(); err != nil {
			t.Errorf("(%v -s 1 -e 2 -l 10 %v) != (%v -s1 -e2 -l10 %v)", selpgName, testFileName, stdName, testFileName)
		}
	}

	{
		selpgOut, _ := usePipe(selpgName, "-s", "1", "-e", "2", "-f", testFileName)
		stdOut, _ := usePipe(stdName, "-s1", "-e2", "-f", testFileName)

		writeToFile(testResultFileName, selpgOut)
		writeToFile(stdResultFileName, stdOut)

		diff := exec.Command(gitName, diffName, testResultFileName, stdResultFileName)
		if _, err := diff.Output(); err != nil {
			t.Errorf("(%v -s 1 -e 2 -f %v) != (%v -s1 -e2 -f %v)", selpgName, testFileName, stdName, testFileName)
		}
	}

	{
		var pipeOutFileName = "pipeOutFile.txt"
		usePipe(selpgName, "-s", "1", "-e", "10000", "-d", pipeExecFileName, testFileName)

		diff := exec.Command(gitName, diffName, pipeOutFileName, testFileName)
		if _, err := diff.Output(); err != nil {
			t.Errorf("(%v -s 1 -e 10000 -d %v %v) Fail", selpgName, pipeExecFileName, testFileName)
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	// os.Remove(testFileName)
}

func usePipe(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	var out = make([]byte, 0, 1024)
	for {
		tmp := make([]byte, 128)
		n, err := stdout.Read(tmp)
		out = append(out, tmp[:n]...)
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return nil, err
	}

	return out, nil
}
