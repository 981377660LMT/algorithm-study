from collections import Counter

# 每次操作会使排列变成它的前一个排列，
# 题目等价于求从小到大数第几个排列-1，
# 或者相当于求有多少个排列比它小
# 求在所有的`排列组合中`，当前这个是字典序第k小 (k从0开始)
# 有点类似康托展开的思想


class Solution:
    def makeStringSorted(self, s: str) -> int:
        ################ 思路：在所有排列中，这是第几小的
        MOD = 1000000007

        def quick_pow(x, n) -> int:  # 快速幂，偷懒了，直接调函数
            return pow(x, n, MOD)

        n = len(s)

        fac = [0 for _ in range(n + 1)]  # 阶乘 % MOD
        fac[0] = 1
        reverse_fac = [0 for _ in range(n + 1)]  # 乘法逆元
        reverse_fac[0] = [1]

        for x in range(1, n):
            fac[x] = fac[x - 1] * x % MOD
            reverse_fac[x] = quick_pow(fac[x], MOD - 2)

        ## 统计每个字母出现的次数
        ch_freq = Counter(s)

        res = 0
        ######## 挨着求出，比s[i]小的字符数量
        for i in range(n - 1):
            ######## 比s[i]小的字母的，总出现频率
            less_cnt = sum(freq for ch, freq in ch_freq.items() if ch < s[i])
            #### 当前位的选择*剩下的全排列
            numerator = less_cnt * fac[n - 1 - i] % MOD
            #### 排列组合 借助乘法逆元  其实python3直接计算就行。这样写就是为了和c++代码同步
            cur = numerator
            for _, freq in ch_freq.items():
                # 除以freq!等于乘以freq!的逆元
                cur = cur * reverse_fac[freq] % MOD

            res = (res + cur) % MOD

            # 删除
            ch_freq[s[i]] -= 1
            if ch_freq[s[i]] == 0:
                del ch_freq[s[i]]

        return res


print(Solution().makeStringSorted(s="cba"))
# 输出：5
# 解释：模拟过程如下所示：
# 操作 1：i=2，j=2。交换 s[1] 和 s[2] 得到 s="cab" ，然后反转下标从 2 开始的后缀字符串，得到 s="cab" 。
# 操作 2：i=1，j=2。交换 s[0] 和 s[2] 得到 s="bac" ，然后反转下标从 1 开始的后缀字符串，得到 s="bca" 。
# 操作 3：i=2，j=2。交换 s[1] 和 s[2] 得到 s="bac" ，然后反转下标从 2 开始的后缀字符串，得到 s="bac" 。
# 操作 4：i=1，j=1。交换 s[0] 和 s[1] 得到 s="abc" ，然后反转下标从 1 开始的后缀字符串，得到 s="acb" 。
# 操作 5：i=2，j=2。交换 s[1] 和 s[2] 得到 s="abc" ，然后反转下标从 2 开始的后缀字符串，得到 s="abc" 。
