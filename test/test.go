package main

// import (
// 	"os"
// 	"bufio"
// 	"fmt"
// 	"strconv"
// )

// func main()  {
// 	if len(os.Args) < 2 {
// 		return
// 	}
// 	file, err := os.Create(os.Args[1])
// 	if err != nil {
// 		fmt.Printf("Write file %q error:", os.Args[1])
// 		panic(err)
// 	}
// 	defer file.Close()

// 	writer := bufio.NewWriter(file)
// 	for i := 1; i <= 1000; i++  {
// 		// fmt.Printf("%v ", i)
// 		if nn, err := writer.WriteString(strconv.Itoa(i)); err != nil {
// 			fmt.Fprintf(os.Stderr, "Expect write %v bytes but %v bytes.\n", len(strconv.Itoa(i)), nn)
//             panic(err)
//         }
// 		if err := writer.WriteByte('\n'); err != nil {
// 			fmt.Fprintf(os.Stderr, "Expect write \\n but 0 bytes.\n")
//             panic(err)
//         }
// 		if i % 10 == 0 {
// 			if err := writer.WriteByte('\f'); err != nil {
// 				fmt.Fprintf(os.Stderr, "Expect write \\f but 0 bytes.\n")
// 				panic(err)
// 			}
// 		}
// 	}
// 	if err = writer.Flush(); err != nil {
//         panic(err)
//     }
// }

import (
	"fmt"
	"io"
	"bufio"
	"io/ioutil"
	"os"
)

func main()  {
	file, err := os.Create("pipeOutFile.txt")
	if err != nil {
		fmt.Printf("Write file %q error:", os.Args[1])
		panic(err)
	}
	// defer file.Close()
	var rc io.ReadCloser
	rc = file
	defer rc.Close()

	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(os.Stdin)
	data, err := ioutil.ReadAll(reader)
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Fprintf(writer, string(data))
	if err = writer.Flush(); err != nil {
		panic(err)
	}
}
