// 以字节为单位返回字符串的长度。
const byteSize = (str: BlobPart) => new Blob([str]).size
byteSize('😀') // 4
byteSize('Hello World') // 11
