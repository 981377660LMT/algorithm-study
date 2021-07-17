"use strict";
// 设计一个算法，找出只含素因子 2，3，5 的第 n 小的数
// 1 通常被视为丑数。
// 1 <= n <= 1690
Object.defineProperty(exports, "__esModule", { value: true });
// 这个while循环终止条件值得借鉴
const nthUglyNumber = (n) => {
    const res = [1];
    let i = 0;
    let j = 0;
    let k = 0;
    while (!res[n - 1]) {
        const nextI = res[i] * 2;
        const nextJ = res[j] * 3;
        const nextK = res[k] * 5;
        const nextUglyNumber = Math.min.call(null, nextI, nextJ, nextK);
        if (nextUglyNumber === nextI)
            i++;
        if (nextUglyNumber === nextJ)
            j++;
        if (nextUglyNumber === nextK)
            k++;
        res.push(nextUglyNumber);
    }
    return res[res.length - 1];
};
console.log(nthUglyNumber(10));
