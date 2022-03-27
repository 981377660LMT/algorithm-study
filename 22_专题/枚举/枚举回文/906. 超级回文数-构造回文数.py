# 如果一个正整数自身是回文数，而且它也是一个回文数的平方，那么我们称这个数为超级回文数。
# 返回包含在范围 [L, R] 中的超级回文数的数目。
# L 和 R 是表示 [1, 10^18) 范围的整数的字符串。

# https://leetcode.com/problems/super-palindromes/discuss/174835/tell-you-how-to-get-all-super-palindrome(detailed-explanation)
# 1.获取所有小于1e9的回文数(因为input<10^18)
# 2.遍历这些回文数 一一检验

# 对每个x 我们可以构造出11个回文数  xx,x0x,x1x,..,x9x

palindrome = [1, 2, 3, 4, 5, 6, 7, 8, 9]
# 构造九位数的回文只需要左右部分最多四位 一共有11*10000+9个回文数
for side in range(1, 10000):
    s1 = str(side) + str(side)[::-1]
    palindrome.append(int(s1))
    for mid in range(10):
        s2 = str(side) + str(mid) + str(side)[::-1]
        palindrome.append(int(s2))
palindrome.sort()


class Solution:
    def superpalindromesInRange(self, left: str, right: str) -> int:
        res = []
        for p in palindrome:
            cur = p ** 2
            if cur < int(left):
                continue
            if cur > int(right):
                break
            cand = str(p ** 2)
            if cand == cand[::-1]:
                res.append(cand)
        return len(res)


print(Solution().superpalindromesInRange(left="4", right="1000"))
# 输出：4
# 解释：
# 4，9，121，以及 484 是超级回文数。
# 注意 676 不是一个超级回文数： 26 * 26 = 676，但是 26 不是回文数。
