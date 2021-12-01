from typing import List

# 用新的字母组替换原有的字母组
# 如果 x 从原始字符串 S 中的位置 i 开始，那么就用 y 替换出现的 x。如果没有，则什么都不做。
# 保证在替换时不会有任何重叠

# 总结：先`倒着排序`再替换
# 这样的好处是替换后长度变化不会再干扰我们


class Solution:
    def findReplaceString(
        self, str: str, indices: List[int], sources: List[str], targets: List[str]
    ) -> str:
        for i, s, t in sorted(zip(indices, sources, targets), reverse=True):
            slice = str[i : i + len(s)]
            if slice == s:
                str = str[:i] + t + str[i + len(s) :]
        return str


print(
    Solution().findReplaceString(
        str="abcd", indices=[0, 2], sources=["a", "cd"], targets=["eee", "ffff"]
    )
)
# 输出："eeebffff"
# 解释：
# "a" 从 S 中的索引 0 开始，所以它被替换为 "eee"。
# "cd" 从 S 中的索引 2 开始，所以它被替换为 "ffff"。
