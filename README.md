# selpg

标签： go golang

---

## 1. 概览
使用 golang 开发 [开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html) 中的 selpg

#### 程序逻辑：
> selpg 是从文本输入选择页范围的实用程序。该输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入。
selpg 首先处理所有的命令行参数。在扫描了所有的选项参数（也就是那些以连字符为前缀的参数）后，如果 selpg 发现还有一个参数，则它会接受该参数为输入文件的名称并尝试打开它以进行读取。如果没有其它参数，则 selpg 假定输入来自标准输入。

#### 参数处理
 - `"-s Number"`和`"-e Number"`强制选项：
 > selpg 要求用户用两个命令行参数`"-s Number"`（例如，`"-s 10"`表示从第 10 页开始）和`"-e Number"`（例如，`"-e 20"`表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。

 - `"-l Number"`和`"-f"`可选选项：
 > `"-l Number"`：该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-lNumber”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。
`"-f"`：该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用`"\f"`表示）定界。该格式与“每页行数固定”格式相比的好处在于，当每页的行数有很大不同而且文件有很多页时，该格式可以节省磁盘空间。在含有文本的行后面，类型 2 的页只需要一个字符 ― 换页 ― 就可以表示该页的结束。打印机会识别换页符并自动根据在新的页开始新行所需的行数移动打印头。

 - `"-d Destination"`可选选项：
 > selpg 还允许用户使用`"-d Destination"`选项将选定的页直接发送至打印机。
（**由于没有打印机测试，为了进行测试，所以实现有些差别，`Destination`是另一个可执行程序，将`selpg`的输出通过管道输给`Destination`**)

## 2. 使用
键入`selpg -h`或在参数错误时可见到以下使用方法。
```
selpg version 0.0.1
Usage: selpg -s start_page -e end_page [ -f | -l lines_per_page ] [ -d dest ] [ in_filename ]

Options:
  -d printer destination
        requires a printer destination
  -e end page
        requires a end page number
  -f    specify '\p' per page
  -h    get help
  -l n
        specify n lines per page (default 72)
  -s start page
        requires a start page number
```
例如可输入`selpg -s 1 -e 1 input_file`或`selpg -s 10 -e 20 -l 66 input_file`或`selpg -s 10 -e 20 -f input_file`或`selpg -s10 -e20 -d lp input_file`

## 3. 设计
程序主要逻辑有两个：首先解析、检查参数得到需求，而后按照需求进行读取并输出。所以代码分成两部分来分别实现这两部分逻辑。并打包成同一个包`lib`方便主函数调用。
 1. 解析参数并检查
 此处通过golang中自带的包[flag](https://golang.org/pkg/flag/)进行参数解析。并在此处进行参数合法性的检查，若不合法，则给出相应提示。譬如起始页码应为正数，指定行数定页时每页行数应为正数，输入的文件应存在并可读等。
此处功能集成在`usage.go`文件里。
 2. 读写文件
 由于输入源可能为磁盘文件，也可能为标准输入。输出源可能为标准输出，也可能是管道。但这些都可以看成是读写接口。在golang中有[io.Reader](https://golang.org/pkg/io/#Reader)和[io.Writer](https://golang.org/pkg/io/#Writer)来统一。
而后读取行或页都可以通过[Reader.ReadBytes()](https://golang.org/pkg/bufio/#Reader.ReadBytes)方法来进行读取。如行的读取可通过`ReadBytes('\n')`，页的读取可通过`ReadBytes('\f')`。而无论是按行读取，还是按页读取都是读取一个文件的某些范围。按行读取的范围是`[start * pageLen, end * pageLen]`，按页读取的范围是`[start, end]`。所以可以统一为一个函数
```go
func CopyRange(src io.Reader, dest io.Writer, start, end int, delim byte) error
```
来进行读取。
但实际上有[Scanner](https://golang.org/pkg/bufio/#Scanner)实现了按行读取，所以为了练手，也实现了`Scanner`方式的读取行范围并应用在按行读取中。
此处功能集成在`io.go`文件里。

## 4. 测试
每个文件都写了相应的test文件。可用`go test`来进行测试。
如对`io.go`中的``函数的测试：
```go
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
```

由于缺少打印机，所以写了一个`test/test.go`来生成`test`来代替打印机命令`lp`进行测试。
而对主程序的测试选择了与[原文](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)中给出的[样例代码](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html#artrelatedtopics)进行比对。测试的逻辑是，即编译样例代码成一个程序，并在测试里启动样例程序与本程序，设置相同参数，将结果交给`git diff`来进比对。这些在`test/selpg_test.go`文件中可看到。执行测试时需要找到`gcc,selpg,git,test`等程序。

至于`shell`中的重定向等功能，由于以上`selpg_test.go`与源程序进行了比对，所以不再详细测试。
