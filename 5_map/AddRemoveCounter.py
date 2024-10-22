from collections import defaultdict


counter = defaultdict(int)


def add(x: int) -> None:
    preFreq = counter[x]
    counter[x] += 1
    if preFreq:
        # remove (preFreq, x)
        ...
    # add (preFreq+1, x)


def remove(x: int) -> bool:
    preFreq = counter[x]
    if preFreq == 0:
        return False
    counter[x] -= 1
    # remove (preFreq, x)
    if preFreq > 1:
        # add (preFreq-1, x)
        ...
    else:
        counter.pop(x)
    return True
