from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def deleteText(self, article: str, index: int) -> str:
        if article[index] == ' ':
            return article

        left = index
        while left > 0 and article[left] != ' ':
            left -= 1

        right = index
        while right < len(article) and article[right] != ' ':
            right += 1

        return (article[:left] + article[right:]).lstrip()


print(Solution().deleteText(article="Singing dancing in the rain", index=10))
print(Solution().deleteText(article="Hello World", index=2))
