from typing import List

# 1 <= arr.length <= 16
# 1 <= arr[i].length <= 26

# 字符串 s 是将 arr 某一子序列字符串连接所得的字符串，
# 如果 s 中的每一个字符都只出现过一次，那么它就是一个可行解。


# 暴力解法 复杂度 2^n
class Solution:
    def maxLength(self, arr: List[str]) -> int:
        dp = [set()]
        for word in arr:
            if len(set(word)) < len(word):
                continue
            curSet = set(word)
            # 注意此处:`while iterating the list, we can't modify` 所以要取dp[:]
            for preSet in dp[:]:
                if curSet & preSet:
                    continue
                dp.append(preSet | curSet)
        return max(len(s) for s in dp)


print(Solution().maxLength(arr=["un", "iq", "ue"]))
# 输出：4
# 解释：所有可能的串联组合是 "","un","iq","ue","uniq" 和 "ique"，最大长度为 4。
a = [1, 2]
for i in a[:]:
    print(i)
    a.append(7)

