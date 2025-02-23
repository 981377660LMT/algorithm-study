class Solution:
    def hasSameDigits(self, s: str) -> bool:
        sb = []
        for _ in range(len(s) - 2):
            for i in range(len(s) - 1):
                sb.append((int(s[i]) + int(s[i + 1])) % 10)
            s = "".join(map(str, sb))
            sb.clear()
        return s[0] == s[1]
