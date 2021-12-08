from math import comb

# 1 <= n <= 50
# 给你一个整数 n，请返回长度为 n 、仅由元音 (a, e, i, o, u) 组成且按 字典序排列 的字符串数量。

# 将 n 个小球放到 5 个盒子里，盒子可以为空。
# x1+x2+...+xn=5 => (x1+1)+(x2+1)+..+(xn+1)=n+5
# 答案是 C(n + 5 - 1, 5 - 1) = C(n + 4, 4)


class Solution:
    def countVowelStrings(self, n: int) -> int:
        return comb(n + 4, 4)


print(Solution().countVowelStrings(n=2))
# 输出：15
# 解释：仅由元音组成的 15 个字典序字符串为
# ["aa","ae","ai","ao","au","ee","ei","eo","eu","ii","io","iu","oo","ou","uu"]
# 注意，"ea" 不是符合题意的字符串，因为 'e' 在字母表中的位置比 'a' 靠后

