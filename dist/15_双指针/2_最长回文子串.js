"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 双指针中心扩展
// 中心扩散法怎么去找回文串？
// 从每一个位置出发，向两边扩散即可。遇到不是回文的时候结束。
const longestPalindrome = (str) => {
    if (str.length <= 1)
        return str;
    let max = '';
    for (let i = 0; i < str.length; i++) {
        //  技术和偶数两种
        for (let j = 0; j <= 1; j++) {
            let left = i;
            let right = i + j;
            // 向两边扩张
            while (left >= 0 && right < str.length && str[left] === str[right]) {
                left--;
                right++;
            }
            if (right - left > max.length) {
                max = str.slice(left + 1, right - 1 + 1);
            }
        }
    }
    return max;
};
console.log(longestPalindrome('babad'));
