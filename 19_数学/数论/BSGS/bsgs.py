"""bsgs与exbsgs"""


def bsgs(base: int, target: int, p: int) -> int:
    """Baby-step Giant-step

    在base和p互质的情况下,求解 base^x ≡ target (mod p) 的最小解x
    """
    ...


def exbsgs(base: int, target: int, p: int) -> int:
    """Extended Baby-step Giant-step

    求解 base^x ≡ target (mod p) 的最小解x
    """
    ...


if __name__ == "__main__":
    ...
