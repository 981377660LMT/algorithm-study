def isConsecutive(x: str, y: str) -> bool:
    d = abs(ord(x) - ord(y))
    return d == 1 or d == 25
