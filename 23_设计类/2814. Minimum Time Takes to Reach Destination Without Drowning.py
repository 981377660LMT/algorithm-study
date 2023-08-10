from typing import List


class Solution:
    def minimumSeconds(self, land: List[List[str]]) -> int:
        ROW, COL = len(land), len(land[0])
