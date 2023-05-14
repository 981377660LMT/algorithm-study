from functools import lru_cache


@lru_cache(None)
def grundy(state: int) -> int:
    nexts = set()

    # 将转移状态加入nexts
    # nexts.add(grundy(nextState))

    mex = 0
    while mex in nexts:
        mex += 1
    return mex


# !初始的grundy数:母状态的 SG 数等于各个子状态的 SG 数的异或
# sg = 0
# for state in states:
#     sg ^= grundy(state)
