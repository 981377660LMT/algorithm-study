class Solution:
    def digitsCount(self, d: int, low: int, high: int) -> int:
        @lru_cache(None)
        # 0~tar之间整数x出现的次数
        def find(x, tar):
            if tar < x:
                return 0
            if tar < 10:
                return 1
            # 取十位数以上部分q和个位数部分res
            q, res = divmod(tar, 10)
            # tar个位数上x出现的次数：
            # 举例：x = 3, tar = 124
            # 1、
            # 0 ~ 120范围内个位数上出现3的次数：##3 当个位数是3，十位和百位的取值范围是0~11总共
            # 12个数，即等于 124 // 10 = 12
            # 2、
            # 121 ~ 124范围内个位数上出现3的次数：find(3, 4)
            result = q + find(x, res)
            # tar十位数及以上部分x出现的次数：
            # 1、
            # 对于0 ~ 119，当十位数以上部分出现3时（find(3, 12-1)），个位数的取值是0~9共10种
            # 情况，因为是统计十位数以上部分出现3的次数，整数不能存在前导0，所以要排除掉find(3,0)
            # (当d=0时该部分等于1)，所以0 ~ 119十位数以上部分出现3的次数应该是
            # (find(3, 11) - find(3, 0)) * 10
            # 2、
            # 对于120 ~ 124，当十位数以上部分出现3时（find(3, 12) - find(3, 11)），个位数的取
            # 值是0~4共5种情况
            result += (find(x, q - 1) - find(x, 0)) * 10 + (find(x, q) - find(x, q - 1)) * (res + 1)
            return result

        return find(d, high) - find(d, low - 1)

