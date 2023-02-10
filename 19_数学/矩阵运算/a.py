import numpy as np


def matqpow2(base: np.ndarray, exp: int, mod: int) -> np.ndarray:
    def inner(base: np.ndarray, exp: int, mod: int) -> np.ndarray:
        res = np.eye(*base.shape, dtype=np.uint64)
        bit = 0
        while exp:
            if exp & 1:
                res = (res @ pow2[bit]) % mod
            exp //= 2
            bit += 1
            if bit == len(pow2):  # 预处理
                pow2.append((pow2[-1] @ pow2[-1]) % mod)
        return res

    pow2 = [base]
    return inner(base, exp, mod)
