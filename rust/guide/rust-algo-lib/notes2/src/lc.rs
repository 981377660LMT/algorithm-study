// def rob(nums: List[int]) -> int:
//     if not nums:
//         return 0

//     dp0, dp1 = 0, nums[0]
//     for i in range(1, len(nums)):
//         dp0, dp1 = max(dp0, dp1), dp0 + nums[i]
//     return max(dp0, dp1)

// 给你一个整数数组 nums 和一个二维数组 queries，其中 queries[i] = [posi, xi]。

// 对于每个查询 i，首先将 nums[posi] 设置为 xi，然后计算查询 i 的答案，该答案为 nums 中 不包含相邻元素 的子序列的 最大 和。

// 返回所有查询的答案之和。

// 由于最终答案可能非常大，返回其对 109 + 7 取余 的结果。

// 子序列 是指从另一个数组中删除一些或不删除元素而不改变剩余元素顺序得到的数组。
struct Solution {}

// impl Solution {
//     pub fn maximum_sum_subsequence(nums: Vec<i32>, queries: Vec<Vec<i32>>) -> i32 {
//         unsafe { Solution::maximum_sum_subsequence_unsafe(nums, queries) }
//     }

//     #[target_feature(enable = "avx2")]
//     pub unsafe fn maximum_sum_subsequence_unsafe(nums: Vec<i32>, queries: Vec<Vec<i32>>) -> i32 {
//         let mut res: i64 = 0;
//         let mut nums64 = nums.into_iter().map(|x| x as i64).collect::<Vec<i64>>();

//         let mut pre_res: Option<i64> = None;
//         for query in queries {
//             let pos = query[0] as usize;
//             let x = query[1] as i64;
//             if nums64[pos] <= 0 && x <= 0 && pre_res.is_some() {
//                 res += pre_res.unwrap();
//                 nums64[pos] = x;
//                 continue;
//             }

//             nums64[pos] = x;
//             let mut dp0: i64 = 0;
//             let mut dp1: i64 = nums64[0];
//             let mut tmp: i64 = 0;
//             for v in nums64.iter().skip(1) {
//                 tmp = dp0.max(dp1);
//                 dp1 = dp0 + v;
//                 dp0 = tmp;
//             }

//             let cur_res = dp0.max(dp1);
//             pre_res = Some(cur_res);
//             res += cur_res;
//         }

//         (res % 1_000_000_007) as i32
//     }
// }

impl Solution {
    pub fn maximum_sum_subsequence(nums: Vec<i32>, queries: Vec<Vec<i32>>) -> i32 {
        unsafe { Solution::maximum_sum_subsequence_unsafe(nums, queries) }
    }

    #[target_feature(enable = "avx2")]
    pub unsafe fn maximum_sum_subsequence_unsafe(nums: Vec<i32>, queries: Vec<Vec<i32>>) -> i32 {
        let mut res: i64 = 0;
        let mut nums64 = nums.into_iter().map(|x| x as i64).collect::<Vec<i64>>();

        for query in queries {
            let pos = query[0] as usize;
            let x = query[1] as i64;

            nums64[pos] = x;
            let mut dp0: i64 = 0;
            let mut dp1: i64 = nums64[0];
            let mut tmp: i64 = 0;
            for v in nums64.iter().skip(1) {
                tmp = dp0.max(dp1);
                dp1 = dp0 + v;
                dp0 = tmp;
            }

            res += dp0.max(dp1);
        }

        (res % 1_000_000_007) as i32
    }

    fn bar(args: &[String]) {}
}
