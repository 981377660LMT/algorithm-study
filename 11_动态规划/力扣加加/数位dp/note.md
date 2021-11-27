https://leetcode-cn.com/problems/non-negative-integers-without-consecutive-ones/solution/shu-wei-dpmo-ban-ji-jie-fa-by-initness-let3/
对于「数位 DP」题，都存在「`询问 [a, b]（a 和 b 均为正整数，且 a < b）区间内符合条件的数值个数为多少`」的一般形式，通常我们需要实现一个查询 **[0, x] 有多少**合法数值的函数 int dp(int x)，然后应用「容斥原理」求解出 [a, b] 的个数：dp(b)−dp(a−1)。

这类题的特征是 `N在10^9量级` 要根据数位而不是遍历数值来做

常见的 dp 定义为 dp[i][j] 其中 i 为数字的长度，
j 为最后一位的数。比如 dp[3][2] 表示这个数一共三位，最后一位是 2 的情况

本质上还是数学题

数位 dp 的一种通用解法就是从数字的最高位开始枚举
`600. 不含连续1的非负整数.ts`

模板

1. 预处理出在枚举过程中需要用到的值
   数组的第一维是数字的长度，第二维是数字的最高位，dp[i][j]表示长为 i 的一个二进制数字，最高位为 j 的情况下包含的不含连续的 1 的数字数量
2. 按照上面所述的方法开始枚举计算：
   (1) 首先将要求的数字转换位一个二进制数组 nums 的形式
   (2) 从最高位开始枚举( `let i=n-1;~i;i--`),枚举从 0 到当前位减一`(j=0;i<nums[i];j++)`，并在枚举过程中计算结果(详见代码注释)
   (3) 返回结果

```C++


class Solution {
private:
    enum {N = 32};
    int dp[N][2];
public:
    int findIntegers(int n) {
        init();
        return dp_func(n);
    }
private:
    void init()
    {
        dp[1][0] = 1, dp[1][1] = 1; //只有1位时，最高位取0或1都是一个合法的数字
        for(int i = 2; i < N; ++i)
        {
            dp[i][0] = dp[i - 1][1] + dp[i - 1][0]; //最高位为0，低一位取0或1都可以
            dp[i][1] = dp[i - 1][0]; //最高位为1，低一位只能取0，来确保没有连续的一
        }
    }

    int dp_func(int n)
    {
        vector<int> nums;
        while(n) nums.push_back(n & 1), n >>= 1; //将n转换为二进制并存在vector里面
        int last(0), res(0); //res是当前已统计的合法数字数量，last是当前枚举的高一位的值
        for(int i = nums.size() - 1; i >= 0; --i) //从最高位开始枚举
        {
            int x = nums[i];
            for(int j = 0; j < x; ++j) //枚举从0到当前位减一
                res += dp[i + 1][j]; //在当前位枚举到自己本身的值之前，它所包含的合法数字数量就对应了之前预处理出的dp值
            if(last == 1 && x == 1) break; //当前位的值和高一位的值都是1时，说明数字已经不合法，直接break
            last = x; //last记录为当前位
            if(!i) ++res; //如果已经遍历到最低一位的本身的值，需要加上这个合法数字，对应的是图中右下方的方案
        }
        return res; //返回结果
    }
};


```

分两类讨论
`902. 最大为 N 的数字组合.ts`
`1012. 至少有 1 位重复的数字.py`

1. 若给定数字是 n 位
   那么先考虑[1,n-1]位长度
   再对 n 位长度考虑
2. 处理 n 位的情况 一般是遍历原 n 的每位 digit
3. 注意临界条件:每个数字都看过，没有 break:当候选条件就等于 n 时 要加 1
