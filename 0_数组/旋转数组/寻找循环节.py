def findRepetend(word: str) -> str:
    period = (word * 2).find(word, 1, -1)
    return word[:period] if period != -1 else ''


if __name__ == '__main__':
    word1 = 'niconico'
    word2 = 'niconiconi'
    assert findRepetend(word1) == 'nico'
    assert findRepetend(word2) == ''
