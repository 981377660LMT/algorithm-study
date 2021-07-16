"use strict";
// 字符串 aabcccccaaa 可压缩为 a2b1c5a3
// 可以假设字符串仅包括 a-z 的字母。
Object.defineProperty(exports, "__esModule", { value: true });
// 思路:使用reduce找分界线
const compress = (str) => str
    .split('')
    .reduce((pre, cur, index, arr) => pre + (arr[index + 1] === arr[index] ? cur : cur + '-'), '')
    .slice(0, -1)
    .split('-')
    .reduce((pre, cur) => pre + cur[0] + cur.length, '');
console.log(compress('aabcccccaaa'));
