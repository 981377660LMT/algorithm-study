"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 返回字符串 s 中长度为 k 的单个子字符串中可能包含的最大元音字母数。
// 用集合比数组的includes快
const maxVowels = (s, k) => {
    const vowels = new Set(['a', 'e', 'i', 'o', 'u']);
    let l = 0;
    let r = k - 1;
    let max = 0;
    let cur = 0;
    for (let i = 0; i < k; i++) {
        if (vowels.has(s[i])) {
            max++;
            cur++;
        }
    }
    while (r < s.length - 1) {
        l++;
        r++;
        if (vowels.has(s[l - 1]))
            cur--;
        if (vowels.has(s[r]))
            cur++;
        max = Math.max(max, cur);
        if (max === k)
            break;
    }
    return max;
};
console.log(maxVowels('leetcode', 3));
