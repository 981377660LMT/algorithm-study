# 833. 字符串中的查找与替换
# https://leetcode.cn/problems/find-and-replace-in-string/description/
# 线性做法 https://leetcode.cn/problems/find-and-replace-in-string/solutions/2388853/xian-xing-zuo-fa-pythonjavacgojs-by-endl-uofo/


from typing import List


class Solution:
    def findReplaceString(
        self, s: str, indices: List[int], sources: List[str], targets: List[str]
    ) -> str:
        replace = [(c, 1) for c in s]
        for i, src, tar in zip(indices, sources, targets):
            if s.startswith(src, i):
                replace[i] = (tar, len(src))

        res = []
        ptr = 0
        while ptr < len(s):
            res.append(replace[ptr][0])
            ptr += replace[ptr][1]
        return "".join(res)


print(
    Solution().findReplaceString(
        s="abcd", indices=[0, 2], sources=["a", "cd"], targets=["eee", "ffff"]
    )
)


# 输出："eeebffff"
# 解释：
# "a" 从 S 中的索引 0 开始，所以它被替换为 "eee"。
# "cd" 从 S 中的索引 2 开始，所以它被替换为 "ffff"。
