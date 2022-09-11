"""
京东 解码消息

给你一个映射。a=1,b=2,c=3,...z=26。
请你构造一个字符串，
恰好有K个只包含小写字母的字符串能映射成这个字符串。
1<=k<=10^18

输入一个字符串。
要求构造的字符串长度不能超过1e5

示例：
输入：
3
输出：
111
即“aaa”,"ak","ka"
"""


# 考虑全1的情形，
# dp[i]代表连续i个1会有多少种方案，有dp[1] = 1, dp[2] = 2,
# !i大于2时，考虑最后一个1,它可能是和前面一个1合并组成一个k,
# !也可能没有合并保留一个a,所以有转移方程 dp[i] = dp[i-1] + dp[i-2]
# 然后考虑其它情况，如果一个子串中数字只包括1，2，本质上和全1没有区别(21,12,22,11都合法)
# 如果出现大于2的数字，如3，那么3不能和它后面的数字组合，
# !所以可以取一个大于2的数字类似作为一个分割符， 总的方案数相当于分隔符之前子串的方案数乘分隔符之后的子串的方案数
# 如果存在方案数为K的串，当且仅当K = dp[x1]*dp[x2]*...*dp[xn],
#  xi为被分隔符分割开的1或者2序列长度

fib = [1, 1]
while fib[-1] <= int(1e18):
    fib.append(fib[-1] + fib[-2])


def solve(k: int) -> str:
    """恰好有K个只包含小写字母的字符串能映射成这个字符串。"""
    if k == 1:
        return "1"

    counter = [0] * len(fib)
    for i in range(len(fib) - 1, 1, -1):  # !因子从大到小检查到2
        while k % fib[i] == 0:
            k //= fib[i]
            counter[i] += 1

    if k != 1:
        return ""  # !不存在方案

    res = []
    for i in range(len(counter)):
        if counter[i] > 0:
            res.append(("1" * (i - 1) + "3") * counter[i])
    return "".join(res)


print(solve(10))
print(solve(4))
print(solve(int(1e9)))
print(solve(1))
print(solve(7))
