class DigitSwap:
    base = [10**i for i in range(19)]

    def __init__(self, num: int) -> None:
        if num < 0:
            raise ValueError("num must be non-negative")
        self.num = num
        self.digits = []
        x = num
        while x > 0:
            self.digits.append(x % 10)
            x //= 10

    def swap(self, i: int, j: int) -> int:
        if i == j or self.digits[i] == self.digits[j]:
            return self.num
        digits, base = self.digits, self.base
        return self.num + (digits[j] - digits[i]) * (base[i] - base[j])

    def at(self, i: int) -> int:
        return self.digits[i]

    def __len__(self) -> int:
        return len(self.digits)


if __name__ == "__main__":
    ds = DigitSwap(12345)
    print(ds.swap(0, 1))  # 12354
