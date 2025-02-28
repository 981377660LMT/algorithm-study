// 158. 用 Read4 读取 N 个字符 II - 多次调用
// https://leetcode.cn/problems/read-n-characters-given-read4-ii-call-multiple-times/solutions/2477575/yong-read4-du-qu-n-ge-zi-fu-ii-by-leetco-kmwd/
//
// 给你一个文件 file ，并且该文件只能通过给定的 read4 方法来读取，
// 请实现一个方法使其能够使 read 读取 n 个字符。注意：你的 read 方法可能会被调用多次。
//
// 实现思路：
// !read 方法可能会被调用多次，因此我们需要保存上一次读取后的状态
// 如果方法 read 的参数 n 不是 4 的倍数，则调用方法 read4 读取的字符数会多于 n 个字符，
// 在这种情况下，必须将多余字符存储在自定义缓冲区内，下次调用方法 read 时首先读取自定义缓冲区内的多余字符，
// 如果自定义缓冲区内的多余字符个数少于需要读取的字符个数，再调用方法 read4 读取后面的字符。

package main

/**
 * The read4 API is already defined for you.
 *
 *     read4 := func(buf4 []byte) int
 *
 * // Below is an example of how the read4 API can be called.
 * file := File("abcdefghijk") // File is "abcdefghijk", initially file pointer (fp) points to 'a'
 * buf4 := make([]byte, 4) // Create buffer with enough space to store characters
 * read4(buf4) // read4 returns 4. Now buf = ['a','b','c','d'], fp points to 'e'
 * read4(buf4) // read4 returns 4. Now buf = ['e','f','g','h'], fp points to 'i'
 * read4(buf4) // read4 returns 3. Now buf = ['i','j','k',...], fp points to end of file
 */
var solution = func(read4 func([]byte) int) func([]byte, int) int {
	internalBuffer := [4]byte{} // 内部缓冲区，用于存放 read4 返回的字符
	pos := 0                    // 当前指针
	bufSize := 0                // 上次调用 read4 读取的有效字节数

	return func(buf []byte, n int) (numRead int /* 返回读取的字符个数 */) {
		for numRead < n {
			// 如果内部缓冲区中的数据已经全部被消费完毕，则调用 read4 读取字符
			if pos == bufSize {
				bufSize = read4(internalBuffer[:])
				pos = 0
				// 如果读取的字符数为 0，则说明文件已经读取完毕，直接返回
				if bufSize == 0 {
					break
				}
			}

			// 将内部缓冲区中的字符拷贝到 buf 中
			for numRead < n && pos < bufSize {
				buf[numRead] = internalBuffer[pos]
				numRead++
				pos++
			}
		}

		return numRead
	}
}
