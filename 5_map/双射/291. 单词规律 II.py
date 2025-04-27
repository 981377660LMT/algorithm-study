from typing import Dict, Set


class Solution:
    def wordPatternMatch(self, pattern: str, s: str) -> bool:
        """
        Backtracking with bijection constraints.
        We try to assign each pattern character to a non-empty substring of s,
        ensuring:
          1. Consistency: once assigned, reuse that mapping.
          2. Injectivity: no two pattern chars map to the same substring.
        Since pattern and s are ≤20, the exponential search is fine.
        """

        def dfs(pi: int, si: int, mapping: Dict[str, str], used: Set[str]) -> bool:
            if pi == len(pattern) and si == len(s):
                return True
            if pi == len(pattern) or si == len(s):
                return False

            pch = pattern[pi]

            if pch in mapping:
                sub = mapping[pch]
                if not s.startswith(sub, si):
                    return False
                return dfs(pi + 1, si + len(sub), mapping, used)

            for end in range(si + 1, len(s) + 1):
                candidate = s[si:end]
                if candidate in used:
                    continue
                mapping[pch] = candidate
                used.add(candidate)
                if dfs(pi + 1, end, mapping, used):
                    return True
                mapping.pop(pch)
                used.remove(candidate)

            return False

        return dfs(0, 0, dict(), set())


if __name__ == "__main__":
    sol = Solution()
    tests = [
        ("abab", "redblueredblue", True),
        ("aaaa", "asdasdasdasd", True),
        ("aabb", "xyzabcxzyabc", False),
    ]
    for pattern, s, expected in tests:
        result = sol.wordPatternMatch(pattern, s)
        print(f"pattern={pattern!r}, s={s!r} -> {result} (expected {expected})")
