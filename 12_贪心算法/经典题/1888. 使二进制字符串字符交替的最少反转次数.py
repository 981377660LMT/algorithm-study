# 你可以按任意顺序执行以下两种操作任意次：

# 类型 1 ：删除 字符串 s 的第一个字符并将它 添加 到字符串结尾。
# 类型 2 ：选择 字符串 s 中任意一个字符并将该字符 反转 ，也就是如果值为 '0' ，则反转得到 '1' ，反之亦然。

# 请你返回使 s 变成 交替 字符串的前提下， `类型 2 的 最少 操作次数` 。
# 我们称一个字符串是 交替 的，需要满足`任意相邻字符都不同`。
class Solution:
    def minFlips(self, s: str) -> int:
        n = len(s)
        if n == 1:
            return 0

        presum0 = [0 for _ in range(n + 1)]  # 目标：开头为0,0101模式需要的反转次数
        presum1 = [0 for _ in range(n + 1)]  # 目标：开头为1,1010模式需要的反转次数
        for i in range(n):
            c = s[i]
            if i % 2 == 0:
                presum0[i + 1] = presum0[i] + (c == '1')
                presum1[i + 1] = presum1[i] + (c == '0')
            else:
                presum0[i + 1] = presum0[i] + (c == '0')
                presum1[i + 1] = presum1[i] + (c == '1')

        res = 0x7FFFFFFF
        # 偶数长度
        if not n & 1:
            res = min(presum0[n], presum1[n])
        # 奇数长度
        else:
            #  枚举删除前多少个数
            for i in range(n):
                # 偶数模式，插到后面的数变奇数模式
                res1 = presum0[n] - presum0[i] + presum1[i]
                # 奇数模式，插到后面的数变偶数模式
                res2 = presum1[n] - presum1[i] + presum0[i]
                res = min(res, res1, res2)
        return res


print(Solution().minFlips(s="111000"))
# 输出：2
# 解释：执行第一种操作两次，得到 s = "100011" 。
# 然后对第三个和第六个字符执行第二种操作，得到 s = "101010" 。
