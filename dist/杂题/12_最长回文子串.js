"use strict";
// Manacher's ALGORITHM: O(n)时间求字符串的最长回文子串
// 在每个字符的两边都插入一个特殊的符号。比如 abba 变成 #a#b#b#a#
// 将所有可能的奇数/偶数长度的回文子串都转换成了奇数长度
