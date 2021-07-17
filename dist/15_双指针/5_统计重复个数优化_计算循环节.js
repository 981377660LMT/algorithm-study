"use strict";
// 请你找出一个最大整数 m ，以满足 str = [str2, m] 可以从 str1 获得(子序列)。
Object.defineProperty(exports, "__esModule", { value: true });
// 加速循环
const getMaxRepetion = (s1, n1, s2, n2) => {
    let s1Count = 0;
    let s2Count = 0;
    let s2p = 0;
    while (s1Count < n1) {
        for (let index = 0; index < s1.length; index++) {
            const s1letter = s1[index];
            if (s1letter === s2[s2p])
                s2p++;
            if (s2p === s2.length) {
                s2Count++;
                s2p = 0;
            }
        }
        s1Count++;
        if (s2p === 0) {
            // 一共需要循环多少次
            const times = Math.floor(n1 / s1Count);
            s1Count *= times;
            s2Count *= times;
            //这里计数乘循环的次数，继续循环 因为counts1还可能是小于n1的，循环节点不能整除
        }
    }
    return Math.floor(s2Count / n2);
};
console.log(getMaxRepetion('acb', 4, 'ab', 2));
console.log(getMaxRepetion('baba', 11, 'baab', 1));
