from typing import List

# 给定一个起始 IP 地址 ip 和一个我们需要包含的 IP 的数量 n，返回用列表（最小可能的长度）表示的 CIDR块的范围。
# CIDR 块是包含 IP 的字符串，后接斜杠和固定长度。例如：“123.45.67.89/20”。固定长度 “20” 表示在特定的范围中公共前缀位的长度。
class Solution:
    def ipToCIDR(self, ip: str, n: int) -> List[str]:
        ...


print(Solution().ipToCIDR(ip="255.0.0.7", n=10))
# 输出：["255.0.0.7/32","255.0.0.8/29","255.0.0.16/32"]
# 没看懂
