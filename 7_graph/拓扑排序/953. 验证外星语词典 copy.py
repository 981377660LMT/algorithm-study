from typing import List


class Solution:
    def isAlienSorted(self, words: List[str], order: str) -> bool:
        dic = {v: i for i, v in enumerate(order)}

        for i in range(0, len(words) - 1):
            pre, cur = words[i], words[i + 1]
            for j in range(0, min(len(pre), len(cur))):
                if pre[j] != cur[j]:
                    if dic[pre[j]] > dic[cur[j]]:
                        return False
                    break
            else:
                if len(pre) > len(cur):
                    return False

        return True


print(Solution().isAlienSorted(["word", "world", "row"], "worldabcefghijkmnpqstuvxyz"))

