"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const findMaxConsecutiveOnes = (nums) => nums
    .join('')
    .split('0')
    .reduce((pre, cur) => Math.max(pre, cur.length), 0);
console.log(findMaxConsecutiveOnes([1, 1, 0, 1, 1, 1]));
