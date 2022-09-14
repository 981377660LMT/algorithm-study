# 向字符串添加空格
from typing import List


class Solution:
    def addSpaces(self, s: str, spaces: List[int]) -> str:
        """
        数组 spaces 描述原字符串中需要添加空格的下标。
        每个空格都应该插入到给定索引处的字符值 之前 。
        """
        arr = list(s)
        n = len(arr)
        si = len(spaces) - 1
        for i in range(n - 1, -1, -1):
            if i == spaces[si]:
                arr[i] = " " + arr[i]
                si -= 1
        return "".join(arr)


print(Solution().addSpaces("hello world", [5, 7]))
