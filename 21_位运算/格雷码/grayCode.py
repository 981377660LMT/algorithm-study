def grayCode2Binary(grayCode: str) -> str:
    """格雷码转二进制."""
    binary = [grayCode[0]]
    for i in range(1, len(grayCode)):
        if grayCode[i] == "0":
            binary.append(binary[i - 1])
        else:
            binary.append(str(1 - int(binary[i - 1])))
    return "".join(binary)


def binary2GrayCode(binary: str) -> str:
    """二进制转格雷码."""
    grayCode = [binary[0]]
    for i in range(1, len(binary)):
        if binary[i] == grayCode[i - 1]:
            grayCode.append("0")
        else:
            grayCode.append("1")
    return "".join(grayCode)


if __name__ == "__main__":
    for i in range(16):
        print(f"{i:2d} {grayCode2Binary(binary2GrayCode(bin(i)[2:].zfill(4)))}")
    for i in range(16):
        print(f"{i:2d} {binary2GrayCode(bin(i)[2:].zfill(4))}")
    for i in range(16):
        print(f"{i:2d} {grayCode2Binary(bin(i)[2:].zfill(4))}")
