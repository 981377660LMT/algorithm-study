from typing import List


class Solution:
    def findRepeatDocument(self, documents: List[int]) -> int:
        def f(arr: List[int]) -> int:
            for v in arr:
                pos = abs(v) - 1
                if arr[pos] < 0:
                    return pos + 1
                arr[pos] = -abs(arr[pos])
            raise ValueError("No duplicate found")

        for i, v in enumerate(documents):
            documents[i] = v + 1  # 映射到 [1,n]
        res = f(documents)
        return res - 1
