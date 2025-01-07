https://blog.csdn.net/walkskyer/category_12544794.html

Go（Golang）标准库中提供了若干“压缩/解压”相关的包，这些包统一放在 `compress/` 命名空间下，涵盖了常用的无损压缩格式及相关工具。常见的有：

- **compress/flate**：实现了 DEFLATE 算法（如 gzip、zlib 等常用的核心压缩算法）。
- **compress/gzip**：基于 flate，实现了常见的 .gz 文件格式压缩与解压。
- **compress/zlib**：基于 flate，实现了 zlib 压缩格式的压缩与解压。
- **compress/lzw**：实现了 LZW（Lempel–Ziv–Welch）算法，常见于 GIF 图像等场景。
- **compress/bzip2**：仅提供 bzip2 的解压功能（**不**支持压缩）。

本回答将依次介绍这些包所包含的功能、用法要点，以及在实际开发中的常见使用模式和注意事项。

---

# 1. compress/flate

## 1.1 功能概述

`compress/flate` 提供了 **DEFLATE** 算法的实现。DEFLATE 是一种将 LZ77 和哈夫曼编码相结合的无损压缩算法，广泛应用于 ZIP、gzip、PNG 等文件格式中。Go 标准库将它作为底层实现，其他包（如 gzip、zlib）也是基于 `flate` 构建的。

- **压缩**：`flate.NewWriter(w io.Writer, level int)`
- **解压**：`flate.NewReader(r io.Reader)`

其中，**`level`** 表示压缩等级，取值通常在 **`flate.NoCompression`** 到 **`flate.BestCompression`**（或 0~9）之间。越高的等级压缩比越好，但压缩速度越慢。

### 1.2 使用示例

```go
package main

import (
    "bytes"
    "compress/flate"
    "fmt"
    "io"
    "log"
)

func main() {
    // 要压缩的数据
    data := []byte("Hello, Go flate! This is some data that will be compressed.")

    // -----------------------
    // 压缩
    // -----------------------
    var buf bytes.Buffer
    // 创建 flate.Writer，level = flate.BestCompression (9)
    fw, err := flate.NewWriter(&buf, flate.BestCompression)
    if err != nil {
        log.Fatal(err)
    }
    // 写入数据
    _, err = fw.Write(data)
    if err != nil {
        log.Fatal(err)
    }
    // 记得Close，刷出缓冲
    err = fw.Close()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Compressed size:", buf.Len())

    // -----------------------
    // 解压
    // -----------------------
    fr := flate.NewReader(&buf)
    decompressed, err := io.ReadAll(fr)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Decompressed data:", string(decompressed))
}
```

#### 注意

- `flate.Writer` 和 `flate.Reader` 都是流式的。多次写入可以按顺序被读取，但也要注意数据边界的管理。
- 在结束写入后，**务必**调用 `Close()` 刷新缓冲，否则可能丢失部分数据。

---

# 2. compress/gzip

## 2.1 功能概述

`compress/gzip` 在 `flate` 基础上封装了 **gzip** 文件格式的读写功能。Gzip 在传输和存储方面应用广泛，尤其是在 Linux/Unix 系统下常见的 `.tar.gz` 格式（tar 归档 + gzip 压缩），以及 HTTP 协议中的 `Content-Encoding: gzip`。

- **压缩**：`gzip.NewWriter(w io.Writer)` / `gzip.NewWriterLevel(w io.Writer, level int)`
- **解压**：`gzip.NewReader(r io.Reader)`

除了压缩和解压的功能，`gzip.Writer` 还允许设置**文件头部**（如 `Name`, `Comment`, `ModTime` 等），用于在生成的 `.gz` 文件中保存原始文件的元信息。

### 2.2 使用示例

```go
package main

import (
    "bytes"
    "compress/gzip"
    "fmt"
    "io"
    "log"
)

func main() {
    // 要压缩的数据
    data := []byte("Hello, Gzip! This is some data that will be compressed with gzip.")

    // -----------------------
    // 压缩
    // -----------------------
    var buf bytes.Buffer
    gw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression) // level: 1~9
    if err != nil {
        log.Fatal(err)
    }
    // 可选：设置一些元信息
    gw.Name = "example.txt"
    gw.Comment = "An example GZIP data"

    // 写入数据
    _, err = gw.Write(data)
    if err != nil {
        log.Fatal(err)
    }

    // 关闭Writer刷出缓冲
    err = gw.Close()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("gzip compressed size:", buf.Len())

    // -----------------------
    // 解压
    // -----------------------
    gr, err := gzip.NewReader(&buf)
    if err != nil {
        log.Fatal(err)
    }
    defer gr.Close()

    // 可以读取到元信息
    fmt.Println("File Name:", gr.Name)
    fmt.Println("Comment:", gr.Comment)

    // 读取解压后内容
    decompressed, err := io.ReadAll(gr)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Decompressed Data:", string(decompressed))
}
```

#### 注意

- 调用 `gzip.NewWriter` 或 `NewWriterLevel` 创建的 `Writer` 也需要在写完数据后 `Close()`。
- `gzip.Reader` 解压时，会在读取流结束或调用 `Close()` 时完成所有资源释放。

---

# 3. compress/zlib

## 3.1 功能概述

`compress/zlib` 同样基于 `flate` 封装，提供了 **zlib** 格式的读写功能。zlib 与 gzip 都使用 DEFLATE 算法，但**文件头**和**校验**（checksum）部分格式不同，典型场景如 PNG 图像的压缩部分、网络传输等。

- **压缩**：`zlib.NewWriter(w io.Writer)` / `zlib.NewWriterLevel(w io.Writer, level int)`
- **解压**：`zlib.NewReader(r io.Reader)` / `zlib.NewReaderDict(r io.Reader, dict []byte)`

可选的 “字典”（dict）在特定场景（比如某些协议或自定义字典）下提高压缩比。

### 3.2 使用示例

```go
package main

import (
    "bytes"
    "compress/zlib"
    "fmt"
    "io"
    "log"
)

func main() {
    data := []byte("Hello, zlib! This is a test for zlib compression and decompression.")

    // -----------------------
    // 压缩
    // -----------------------
    var buf bytes.Buffer
    zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
    if err != nil {
        log.Fatal(err)
    }
    _, err = zw.Write(data)
    if err != nil {
        log.Fatal(err)
    }
    // 关闭以完成压缩
    err = zw.Close()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("zlib compressed size:", buf.Len())

    // -----------------------
    // 解压
    // -----------------------
    zr, err := zlib.NewReader(&buf)
    if err != nil {
        log.Fatal(err)
    }
    defer zr.Close()

    decompressed, err := io.ReadAll(zr)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Decompressed Data:", string(decompressed))
}
```

---

# 4. compress/lzw

## 4.1 功能概述

`compress/lzw` 实现了 **LZW（Lempel-Ziv-Welch）** 算法。LZW 广泛用于老版 GIF 图像压缩、一些打印机协议等场景。但与 DEFLATE 相比，LZW 的压缩比通常更低，也不如 DEFLATE 常见。

- **压缩**：`lzw.NewWriter(w io.Writer, order lzw.Order, litWidth int)`
- **解压**：`lzw.NewReader(r io.Reader, order lzw.Order, litWidth int)`

`order` 参数指定压缩码流的比特顺序，可以是 `lzw.LSB` 或 `lzw.MSB`；`litWidth` 通常为 8，用于指定初始的字节大小。

### 4.2 使用示例

```go
package main

import (
    "bytes"
    "compress/lzw"
    "fmt"
    "io"
    "log"
)

func main() {
    data := []byte("Hello, LZW compression in Go!")

    // -----------------------
    // 压缩
    // -----------------------
    var buf bytes.Buffer
    // LSBOrder, litWidth=8 用于 GIF 等场景
    lw := lzw.NewWriter(&buf, lzw.LSB, 8)
    _, err := lw.Write(data)
    if err != nil {
        log.Fatal(err)
    }
    lw.Close()
    fmt.Println("lzw compressed size:", buf.Len())

    // -----------------------
    // 解压
    // -----------------------
    lr := lzw.NewReader(&buf, lzw.LSB, 8)
    decompressed, err := io.ReadAll(lr)
    if err != nil {
        log.Fatal(err)
    }
    lr.Close()

    fmt.Println("Decompressed Data:", string(decompressed))
}
```

#### 注意

- LZW 并没有像 gzip 那样的自动校验；也不含文件头的概念，所以需要对字节流格式自行管理。
- GIF 的 LZW 编码比这里还多了“图像帧”和“调色板”等格式信息，标准库仅负责 LZW 的压缩/解压核心。

---

# 5. compress/bzip2

## 5.1 功能概述

`compress/bzip2` 提供了对 **bzip2** 格式数据的解压功能（**仅解压，无压缩**）。bzip2 采用了 Burrows–Wheeler Transform（BWT） 加上 RLE 与哈夫曼编码，压缩率较高，但压缩/解压速度相对 gzip 会慢一些。

- **解压**：`bzip2.NewReader(r io.Reader) io.Reader`

官方并未提供 “bzip2.Writer”，因此无法直接使用标准库进行 bzip2 压缩。如果需要，可以使用第三方库（例如 `github.com/dsnet/compress/bzip2`）来实现 bzip2 的压缩功能。

### 5.2 使用示例

```go
package main

import (
    "compress/bzip2"
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
    // 例：从一个 .bz2 文件中解压
    f, err := os.Open("example.bz2")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // 创建bzip2 Reader
    bzReader := bzip2.NewReader(f)

    // 读取解压后的内容（这里直接打印）
    _, err = io.Copy(os.Stdout, bzReader)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("\nDone")
}
```

---

# 6. 注意事项与最佳实践

1. **流式处理**

   - Go 中的压缩包往往是**流式**的，通过实现 `io.Writer` / `io.Reader` 的接口处理数据。
   - 当数据量很大时，应边读边写，避免一次性读入内存。

2. **Close() 及时调用**

   - 不管是压缩还是解压，通常需要调用 `Close()`（或 `Flush()`）来结束流和释放资源；尤其对写入（Writer）非常重要，否则可能会丢失未写完的缓冲数据。

3. **压缩级别**

   - 级别越高，压缩比越好，但 CPU 开销越高，耗时更久。实际场景中需结合**压缩比**和**CPU消耗**做权衡。
   - 一些高并发服务可能倾向使用 `flate.DefaultCompression` 或 `flate.BestSpeed` 以减轻 CPU 负载。

4. **错误处理**

   - 处理压缩/解压时可能遇到数据损坏、格式不正确等异常，需要做好错误处理。

5. **第三方库**

   - 如果标准库功能无法满足需求（如需 bzip2 压缩、lz4、snappy 或 zstd 等），可使用社区提供的优秀第三方库，例如：
     - [github.com/pierrec/lz4](https://github.com/pierrec/lz4)
     - [github.com/golang/snappy](https://github.com/golang/snappy)
     - [github.com/DataDog/zstd](https://github.com/DataDog/zstd)

6. **并发读写**

   - 压缩写入器或解压读取器一般**不**是线程安全的，不能多个 goroutine 同时对同一个 `Writer` 或 `Reader` 进行写/读操作。
   - 如果要并行处理多段数据，可以**并行创建多个**压缩/解压对象，或者在上层做并行切分。

7. **边压缩边传输**
   - 许多服务端/客户端场景中，会将 `gzip.Writer` 或 `zlib.Writer` 包装在网络连接里进行实时压缩传输。
   - 注意在发送完数据后，要先 `Close()` 压缩流再关闭网络连接，否则对端可能无法解压完整的末尾数据。

---

# 7. 总结

Go 标准库在 `compress/` 包下为开发者提供了**常用的无损压缩/解压**能力，涵盖了 gzip、zlib、lzw、bzip2（解压）等格式，每一种的核心都基于 **io.Reader / io.Writer** 的流式处理理念。它们的主要差异在于**格式头部、校验方式**、**压缩/解压**的效率和压缩比，以及是否提供写入（压缩）支持。

| 包名      | 功能                             | 是否支持压缩 | 是否支持解压 | 常见用途                  |
| --------- | -------------------------------- | ------------ | ------------ | ------------------------- |
| **flate** | DEFLATE 算法底层实现             | 是           | 是           | 底层基础，不常直接使用    |
| **gzip**  | gzip 文件格式（DEFLATE+头/校验） | 是           | 是           | .tar.gz 文件、HTTP 压缩等 |
| **zlib**  | zlib 格式（DEFLATE+头/校验）     | 是           | 是           | PNG、网络传输等           |
| **lzw**   | LZW 算法                         | 是           | 是           | GIF 等少数场景            |
| **bzip2** | bzip2 解压                       | 否           | 是           | bzip2 文件解压（.bz2）    |

在实际项目中，**gzip** 与 **zlib** 是最常用的压缩方式，**bzip2** 只提供解压（需要的话可以使用第三方库来支持压缩），而 **lzw** 场景比较专用。  
如果追求更高性能或更新的算法（如 lz4、snappy、zstd 等），则可以借助第三方开源库。结合 Go 的并发特性，可以轻松地构建高吞吐量、可伸缩的压缩/解压服务。
