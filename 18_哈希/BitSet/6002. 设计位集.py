class Bitset:
    def __init__(self, size: int):
        self.capacity = size
        self.size = 0
        self.bit = 0

    def fix(self, idx: int) -> None:
        idx = self.capacity - 1 - idx
        if not (self.bit & (1 << idx)):
            self.size += 1
            self.bit |= 1 << idx

    def unfix(self, idx: int) -> None:
        idx = self.capacity - 1 - idx
        if self.bit & (1 << idx):
            self.size -= 1
            # 注意消去某位是取反
            self.bit &= ~(1 << idx)

    def flip(self) -> None:
        self.bit ^= (1 << (self.capacity)) - 1
        self.size = self.capacity - self.size

    def all(self) -> bool:
        return self.size == self.capacity

    def one(self) -> bool:
        return self.size != 0

    def count(self) -> int:
        # return self.size
        return self.bit.bit_count()

    def toString(self) -> str:
        cur = bin(self.bit)[2:]
        # padStart
        return cur.zfill(self.capacity)
