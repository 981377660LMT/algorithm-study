from typing import List
from collections import defaultdict


class Solution:
    def wordSquares(self, words: List[str]) -> List[List[str]]:
        """
        Backtracking with prefix-hash acceleration.
        1. Build a map from every possible prefix to the list of words having that prefix.
        2. For each word as the first row, recursively build the word square:
           - At step k, form the prefix by taking the k-th character of each word in the current square.
           - Look up all candidates with that prefix, append each in turn, and recurse.
        Time: O(N · L · A) where N = number of words, L = word length, A = branching factor.
        Space: O(N·L^2) for the prefix map plus O(L) recursion stack.
        """
        if not words:
            return []
        n = len(words[0])

        prefix_map = defaultdict(list)
        for w in words:
            for i in range(n + 1):
                prefix_map[w[:i]].append(w)

        res: List[List[str]] = []
        path: List[str] = []

        def bt(step: int):
            if step == n:
                res.append(path[:])
                return
            prefix = "".join(word[step] for word in path)
            for cand in prefix_map.get(prefix, []):
                path.append(cand)
                bt(step + 1)
                path.pop()

        for w in words:
            path = [w]
            bt(1)

        return res


if __name__ == "__main__":
    sol = Solution()
    words = ["area", "lead", "wall", "lady", "ball"]
    # Possible output (order may vary):
    # [
    #   ["wall","area","lead","lady"],
    #   ["ball","area","lead","lady"]
    # ]
    for square in sol.wordSquares(words):
        for row in square:
            print(row)
        print()
