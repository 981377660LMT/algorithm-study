impl Solution {
    pub fn number_of_subsequences(nums: Vec<i32>) -> i64 {
        let n = nums.len();

        let max_val = *nums.iter().max().unwrap() as usize;

        let mut freq = vec![vec![0u32; max_val + 1]; n + 1];

        for i in (0..n).rev() {
            let (left, right) = freq.split_at_mut(i + 1);
            left[i].copy_from_slice(&right[0]);
            left[i][nums[i] as usize] += 1;
        }

        let mut res: u64 = 0;

        for q in 0..n {
            for r in (q + 2)..n {
                let x = nums[q] as u64;
                let y = nums[r] as u64;

                for p in 0..=q.saturating_sub(2) {
                    let t = (nums[p] as u64) * y;
                    if x != 0 && t % x == 0 {
                        let needed_val = (t / x) as usize;
                        if needed_val != 0 && needed_val <= max_val {
                            if r + 2 < n {
                                res += freq[r + 2][needed_val] as u64;
                            }
                        }
                    }
                }
            }
        }

        res as i64
    }
}
