from itertools import zip_longest


class Solution:
    def compareVersion(self, version1: str, version2: str) -> int:
        for v1, v2 in zip_longest(version1.split("."), version2.split("."), fillvalue="0"):
            n1, n2 = int(v1), int(v2)
            if n1 != n2:
                return 1 if n1 > n2 else -1
        return 0
