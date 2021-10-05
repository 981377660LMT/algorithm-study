function uniqueMorseRepresentations(words: string[]): number {}

console.log(uniqueMorseRepresentations(['gin', 'zen', 'gig', 'msg']))
// 输出: 2
// 解释:
// 各单词翻译如下:
// "gin" -> "--...-."
// "zen" -> "--...-."
// "gig" -> "--...--."
// "msg" -> "--...--."
// 共有 2 种不同翻译, "--...-." 和 "--...--.".

// 这道题如果逆过来，根据摩尔斯码翻译出原始字符串，就是一道很好的DFA的题
