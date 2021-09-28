from collections import Counter


class Solution:
    def frequencySort(self, s: str) -> str:
        return ''.join(w * c for (w, c) in Counter(s).most_common())
