"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 请你找出一个最大整数 m ，以满足 str = [str2, m] 可以从 str1 获得(子序列)。
const getMaxRepetion = (s1, n1, s2, n2) => {
    let s1Index = 0;
    let s2Index = 0;
    for (let i = 0; i < n1; i++) {
        for (let j = 0; j < s1.length; j++) {
            s1Index++;
            if (s1[j] === s2[s2Index])
                s2Index++;
            if (s2Index === s2.length)
                break;
        }
    }
    return Math.floor((s1.length * n1) / s1Index / n2);
};
console.log(getMaxRepetion('acb', 4, 'ab', 2));
