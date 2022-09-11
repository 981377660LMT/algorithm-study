# 大端序转小端序

from typing import List


def convert_little_to_big(numbers: List[int]) -> List[int]:
    return [int.from_bytes(n.to_bytes(4, byteorder="little"), byteorder="big") for n in numbers]


def string_to_numbers(input: str) -> List[int]:
    return [int(n, base=16) for n in input.split()]


def print_big_endian(numbers: List[int]):
    print(" ".join(f"0x{n:08x}" for n in numbers))  # 08表示输出8个字符。x是输出16进制


def main():
    input_str = input()
    numbers = string_to_numbers(input_str)
    big_endian = convert_little_to_big(numbers)
    print(big_endian)
    print_big_endian(big_endian)


if __name__ == "__main__":
    main()
# 0x00112233
