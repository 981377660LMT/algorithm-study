lution:
        def countDistinct(self, nums: List[int], k: int, p: int) -> int:
            r = cnt = res = 0
            right = [0] * (n := len(nums))
            for l, x in enumerate(nums):
                while r < n and cnt + (v := nums[r] % p == 0) <= k:
                    cnt += v
                    r += 1
                res += r - l
                right[l] = r
                cnt -= x % p == 0
            lcp = LCP(nums, sa := SAIS(nums))
            return r