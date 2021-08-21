import itertools


class Solution:
    def isLongPressedName(self, name: str, typed: str) -> bool:
        return all(
            item1[0] == item2[0] and len(list(item1[1])) <= len(list(item2[1]))
            for item1, item2 in itertools.zip_longest(
                itertools.groupby(name), itertools.groupby(typed), fillvalue=("$", 0)
            )
        )


s = Solution()
print(s.isLongPressedName("alex", "aaleex"))

