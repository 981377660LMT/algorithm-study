def lcp(s1: str, s2: str) -> int:
    for i, (a, b) in enumerate(zip(s1, s2)):
        if a != b:
            return i
    return min(len(s1), len(s2))


if __name__ == "__main__":
    print(lcp("ab", "ac"))  # expect: 1
    print(lcp("ab", "ab"))  # expect: 2
    print(lcp("ab", "abc"))  # expect: 2
    print(lcp("ab", "a"))  # expect: 1
    print(lcp("ab", ""))  # expect: 0
