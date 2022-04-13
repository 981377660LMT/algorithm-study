from typing import List


class DisCretizer:
    """离散化"""

    def __init__(self, nums: List[int]) -> None:
        allNums = sorted(set(nums))
        self.mapping = {allNums[i]: i + 1 for i in range(len(allNums))}

    def getDisCretizedValue(self, num: int) -> int:
        if num not in self.mapping:
            raise ValueError(f'{num} not in {self.mapping}')
        return self.mapping[num]

    def __len__(self) -> int:
        return len(self.mapping)


if __name__ == '__main__':
    discretizer = DisCretizer([666, 3, 21])
    print(discretizer.getDisCretizedValue(666))

