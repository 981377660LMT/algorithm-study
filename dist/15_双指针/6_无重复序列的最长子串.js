"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 双指针即可
const lengthOfLongestSubstring = (s) => {
    let left = 0;
    let right = 0;
    let max = 0;
    const map = new Map();
    while (right <= s.length - 1) {
        if (!map.has(s[right])) {
            map.set(s[right], true);
            max = Math.max(map.size, max);
            right++;
        }
        else {
            left++;
            map.delete(s[left - 1]);
        }
    }
    return max;
};
console.log(lengthOfLongestSubstring('abcabcbb'));
