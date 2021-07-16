"use strict";
// 描述
// 给定非空字符串 s，您最多可以删除一个字符。判断是否可以成为回文。
// 该字符串仅包含小写字符 a-z,字符串的最大长度为 50000。
Object.defineProperty(exports, "__esModule", { value: true });
// 思路：
// 设置两个指针，往中间靠拢，然后当遇到s[L] != s[R]的时候，
// 就判断[L + 1, R]或者[L, R-1]部分可以不可以构成回文串，只要有其中一个可以，就返回true即可。
const isValidPalindromw = (n) => {
    let leftIndex = 0;
    let rightIndex = n.length - 1;
    let isValid = true;
    const isAbsolutePalindromw = (n) => n === n.split('').reverse().join('');
    while (leftIndex <= rightIndex) {
        if (n[leftIndex] !== n[rightIndex]) {
            const retry1 = n.slice(leftIndex, rightIndex);
            const retry2 = n.slice(leftIndex + 1, rightIndex + 1);
            // 救不回了
            if (!isAbsolutePalindromw(retry1) && !isAbsolutePalindromw(retry2)) {
                isValid = false;
                break;
            }
        }
        leftIndex++;
        rightIndex--;
    }
    return isValid;
};
console.log(isValidPalindromw('abcddcba'));
console.log(isValidPalindromw('akbcddcba'));
console.log(isValidPalindromw('akkbcddcba'));
