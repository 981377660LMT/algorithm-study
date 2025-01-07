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
package main

func main() {

}
