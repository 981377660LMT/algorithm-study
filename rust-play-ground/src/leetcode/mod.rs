// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/
// let res = 0
// nums.sort((a, b) => a - b)
// for (let i = 0; i < nums.length; i++) {
//   for (let j = i + 1; j < nums.length && nums[j] <= 2 * nums[i]; j++) {
//     res = Math.max(res, nums[i] ^ nums[j])
//   }
// }

impl Solution {
    pub fn maximum_strong_pair_xor(mut nums: Vec<i32>) -> i32 {
        unsafe { Solution::maximum_strong_pair_xor_2(nums) }
    }

    #[target_feature(enable = "avx2")]
    pub unsafe fn maximum_strong_pair_xor_2(mut nums: Vec<i32>) -> i32 {
        let mut res = 0;
        nums.sort();
        let nums2 = nums.clone();
        for (i, v1) in nums.into_iter().enumerate() {
            for v2 in nums2.iter().take(i + 1) {
                if v2 <= &(v1 << 1) {
                    res = res.max(v1 ^ v2);
                }
            }
        }
        return res;
    }
}
