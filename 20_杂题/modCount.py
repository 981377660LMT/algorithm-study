# 闭区间 [0,r] 内模mod与k同余的数的个数
def modCount(right: int, k: int, mod: int) -> int:
    """区间 [0,right] 内模mod与k同余的数的个数"""
    assert 0 <= k < mod
    return (right - k + mod) // mod


if __name__ == "__main__":
    for i in range(3):
        assert sum(1 for j in range(100) if j % 3 == i) == modCount(99, i, 3)
