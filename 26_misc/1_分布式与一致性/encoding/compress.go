// 压缩数据处理的艺术：Go语言compress库完全指南
// https://blog.csdn.net/walkskyer/article/details/135671977
//
// https://blog.csdn.net/qq_42835440/article/details/126687582
// 对数据进行压缩，通常有两个思路：
//
// 1. 字典转换 ( dictionary transforms )
// 减少数据中不同符号的数量（即让“字母表”尽可能小）；【目前所有的主流压缩算法，比如GZIP或者7-Zip，都会在核心转换步骤中使用字典转换】
// 2. 变长编码 （variable-length codes）
// 用更少的位数对更常见的符号进行编码（即最常见的“字母”所用的位数最少）。【bzip2 基于该点对数据进行压缩】
//
// 各种压缩算法对比
// https://vps.fmvz.usp.br/CRAN/web/packages/brotli/vignettes/brotli-2015-09-22.pdf

package main

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	{
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

	{
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

	{
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

	{
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

	{
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
}
