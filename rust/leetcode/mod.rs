struct Solution;

fn max64(a: i64, b: i64) -> i64 {
    if a > b {
        return a;
    }
    return b;
}

// 100341. 网格图操作后的最大分数
// https://leetcode.cn/contest/biweekly-contest-135/problems/maximum-score-from-grid-operations/
impl Solution {
    fn maximum_score(grid: Vec<Vec<i32>>) -> i64 {
        unsafe {
            let n = grid.len();
            if n == 1 {
                return 0;
            }

            let mut colacc = vec![vec![0i64; n + 1]; n];
            for col in 0..n {
                for i in 0..n {
                    *colacc.get_unchecked_mut(col).get_unchecked_mut(i + 1) =
                        colacc.get_unchecked(col).get_unchecked(i)
                            + (*grid.get_unchecked(i).get_unchecked(col) as i64);
                }
            }

            let mut dp = vec![vec![0i64; n + 1]; n + 1];
            for fst in 0..=n {
                for snd in 0..=n {
                    if fst > snd {
                        *dp.get_unchecked_mut(fst).get_unchecked_mut(snd) =
                            colacc.get_unchecked(1).get_unchecked(fst)
                                - colacc.get_unchecked(1).get_unchecked(snd);
                    } else if fst < snd {
                        *dp.get_unchecked_mut(fst).get_unchecked_mut(snd) =
                            colacc.get_unchecked(0).get_unchecked(snd)
                                - colacc.get_unchecked(0).get_unchecked(fst);
                    }
                }
            }

            for col in 2..n {
                let mut ndp = vec![vec![0i64; n + 1]; n + 1];
                let mut hi = vec![0i64; n + 1];
                for fst in 0..=n {
                    for snd in 0..=n {
                        *hi.get_unchecked_mut(snd) = max64(
                            *hi.get_unchecked(snd),
                            *dp.get_unchecked(fst).get_unchecked(snd),
                        );
                        for thd in snd + 1..=n {
                            let extra = *colacc.get_unchecked(col - 1).get_unchecked(thd)
                                - colacc.get_unchecked(col - 1).get_unchecked(fst.max(snd));
                            *ndp.get_unchecked_mut(snd).get_unchecked_mut(thd) = max64(
                                *ndp.get_unchecked(snd).get_unchecked(thd),
                                dp.get_unchecked(fst).get_unchecked(snd) + extra,
                            );
                        }
                    }
                }

                for snd in 0..=n {
                    for thd in 0..=snd {
                        *ndp.get_unchecked_mut(snd).get_unchecked_mut(thd) = max64(
                            *ndp.get_unchecked(snd).get_unchecked(thd),
                            *hi.get_unchecked(snd) + *colacc.get_unchecked(col).get_unchecked(snd)
                                - colacc.get_unchecked(col).get_unchecked(thd),
                        );
                    }
                }

                dp = ndp;
            }

            dp.into_iter()
                .map(|snd| snd.into_iter().max().unwrap())
                .max()
                .unwrap()
        }
    }
}
