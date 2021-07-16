"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 给定不同面额的硬币 coins 和一个总金额 amount。
// 编写一个函数来计算可以凑成总金额所需的最少的硬币个数。
// 如果没有任何一种硬币组合能组成总金额，返回 -1。
const coinChange = (coins, amount) => {
    // 因为Infinity不可能和方法数相等，所以采用Infinity
    const dp = Array(amount + 1).fill(Infinity);
    dp[0] = 0;
    for (let index = 1; index <= amount; index++) {
        for (const coin of coins) {
            if (index - coin >= 0) {
                dp[index] = Math.min(dp[index], dp[index - coin] + 1);
            }
        }
    }
    return dp[amount] === Infinity ? -1 : dp[amount];
};
console.log(coinChange([1, 2, 5], 11));
