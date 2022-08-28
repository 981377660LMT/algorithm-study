"""
Geohash是一个哈希函数,用于将位置坐标对转换为base32字符串。
"""

from typing import List


BASE32 = "0123456789bcdefghjkmnpqrstuvwxyz"
MAPPING = {BASE32[i]: i for i in range(32)}


class GeoHash:
    @staticmethod
    def _getBin(value: float, lower: float, upper: float) -> str:
        res = []
        for _ in range(30):
            mid = (lower + upper) / 2
            if value > mid:
                lower = mid
                res.append("1")
            else:
                upper = mid
                res.append("0")
        return "".join(res)

    @staticmethod
    def _getLocation(bin: str, lower: float, upper: float) -> float:
        for char in bin:
            mid = (lower + upper) / 2
            if char == "1":
                lower = mid
            else:
                upper = mid
        return (lower + upper) / 2

    def encode(self, latitude: float, longitude: float, precision: int) -> str:
        """
        将(纬度、经度)对转换为geohash字符串。
        对经度和维度二分,先经后纬经纬交替,01每五个一组转32进制,然后映射到地图上的位置
        """
        bin1, bin2 = self._getBin(longitude, -180, 180), self._getBin(latitude, -90, 90)

        code = []
        for a, b in zip(bin1, bin2):
            code.append(a)
            code.append(b)
        code = "".join(code)

        res = []
        for i in range(0, len(code), 5):
            res.append(BASE32[int(code[i : i + 5], 2)])
        return "".join(res)[:precision]

    def decode(self, geohash: str) -> List[float]:
        """将geohash字符串转换为(纬度、经度)对。"""
        bin_ = []
        for char in geohash:
            index = MAPPING[char]
            bin_.append(bin(index)[2:].zfill(5))
        bin_ = "".join(bin_)

        bin1 = bin_[::2]
        bin2 = bin_[1::2]
        res1 = self._getLocation(bin1, -180, 180)
        res2 = self._getLocation(bin2, -90, 90)

        return [res2, res1]


if __name__ == "__main__":
    print(GeoHash().encode(39.92816697, 116.38954991, 12))  # "wx4g0s8q3jf9"
    print(GeoHash().decode("wx4g0s"))  # lat = 39.92706299 and lng = 116.39465332
    print(GeoHash().decode("w"))  # lat = 22.50000000 and lng = 112.50000000
