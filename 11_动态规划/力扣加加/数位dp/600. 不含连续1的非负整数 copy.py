class Solution:
    def findIntegers(self, num: int) -> int:
        # ----num的二级制位数。所有≤n的数字形成的完全二叉树的高度
        n = 0
        while num >> n:
            n += 1
        print(n)

        # ----每个index，是0还是1的情况
        dp = [[0 for _ in range(2)] for _ in range(n + 1)]
        dp[0][0] = 1

        for i in range(1, n + 1):
            dp[i][0] = dp[i - 1][0] + dp[i - 1][1]
            dp[i][1] = dp[i - 1][0]

        # ----从左往右
        res = 0
        pre_bit = 0  # 左一位是0还是1
        for i in range(n - 1, -1, -1):
            # ----当前位是1
            if num & (1 << i):
                res += dp[i + 1][0]

                # ----如果出现连续1的情形。xx11xxx，计算停止
                if pre_bit == 1:
                    break
                pre_bit = 1
            # ----当前位是0
            else:
                pre_bit = 0

            # ----如果能到遍历到num，num本身也是一个合理的数字
            if i == 0:
                res += 1

        return res


print(Solution().findIntegers(5))
