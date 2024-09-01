struct Solution;

impl Solution {
    pub fn minimum_total_distance(robot: Vec<i32>, factory: Vec<Vec<i32>>) -> i64 {
        let mut targets = vec![];
        factory.iter().for_each(|v| {
            let (pos, cnt) = (v[0], v[1]);
            for _ in 0..cnt {
                targets.push(pos);
            }
        });

        let mut cost_matrix = vec![vec![0i64; robot.len()]; targets.len()];
        for (i, &r) in robot.iter().enumerate() {
            for (j, &t) in targets.iter().enumerate() {
                unsafe {
                    *cost_matrix.get_unchecked_mut(j).get_unchecked_mut(i) = (r - t).abs() as i64;
                }
            }
        }

        let res = linear_sum_assignment(cost_matrix, false);

        res.0 as i64
    }
}

// TODO: 加泛型.不过速度并不快.
// https://leetcode.cn/problems/minimum-total-distance-traveled/
fn linear_sum_assignment(cost_matrix: Vec<Vec<i64>>, maximize: bool) -> (i64, Vec<i32>, Vec<i32>) {
    let mut n = cost_matrix.len();
    let mut m = cost_matrix[0].len();
    if n > m {
        let mut new_cost_matrix: Vec<Vec<i64>> = vec![vec![0; n]; m];
        for (i, row) in cost_matrix.iter().enumerate() {
            row.iter().enumerate().for_each(|(j, &val)| unsafe {
                *new_cost_matrix.get_unchecked_mut(j).get_unchecked_mut(i) = val;
            });
        }
        let (res, row_index, col_index) = linear_sum_assignment(new_cost_matrix, maximize);
        return (res, col_index, row_index);
    }

    let mut a = vec![vec![0; m + 1]; n + 1];
    let f = if maximize { -1 } else { 1 };
    for (i, row) in cost_matrix.iter().enumerate() {
        row.iter().enumerate().for_each(|(j, &val)| unsafe {
            *a.get_unchecked_mut(i + 1).get_unchecked_mut(j + 1) = f * val;
        });
    }
    n += 1;
    m += 1;

    let mut p = vec![0i32; m];
    let mut way = vec![0i32; m];
    let mut x = vec![0i64; n];
    let mut y = vec![0i64; m];
    let mut min_v = vec![0i64; m];
    let mut used = vec![false; m];

    unsafe {
        for i in 1..n {
            *p.get_unchecked_mut(0) = i as i32;
            for j in 0..m {
                *min_v.get_unchecked_mut(j) = std::i64::MAX;
                *used.get_unchecked_mut(j) = false;
            }
            let mut j0 = 0;
            while *p.get_unchecked(j0) != 0 {
                let i0 = *p.get_unchecked(j0);
                let mut j1 = 0;
                *used.get_unchecked_mut(j0) = true;
                let mut delta = std::i64::MAX;
                for j in 1..m {
                    if *used.get_unchecked(j) {
                        continue;
                    }
                    let curr = *a.get_unchecked(i0 as usize).get_unchecked(j as usize)
                        - x.get_unchecked(i0 as usize)
                        - y.get_unchecked(j as usize);
                    if curr < *min_v.get_unchecked(j) {
                        *min_v.get_unchecked_mut(j) = curr;
                        *way.get_unchecked_mut(j) = j0 as i32;
                    }
                    if *min_v.get_unchecked(j) < delta {
                        delta = *min_v.get_unchecked(j);
                        j1 = j;
                    }
                }
                for j in 0..m {
                    if *used.get_unchecked(j) {
                        x[*p.get_unchecked(j) as usize] += delta;
                        y[j] -= delta;
                    } else {
                        *min_v.get_unchecked_mut(j) -= delta;
                    }
                }
                j0 = j1;
            }

            loop {
                *p.get_unchecked_mut(j0) = *p.get_unchecked(*way.get_unchecked(j0) as usize);
                j0 = *way.get_unchecked(j0) as usize;
                if j0 == 0 {
                    break;
                }
            }
        }
    }

    let mut res = -y[0];
    x.remove(0);
    y.remove(0);
    let row_index = (0..n as i32 - 1).collect::<Vec<i32>>();
    let mut col_index = vec![0i32; n as usize];
    p.iter().enumerate().for_each(|(i, &val)| unsafe {
        *col_index.get_unchecked_mut(val as usize) = i as i32 - 1;
    });
    col_index.remove(0);
    if maximize {
        res = -res;
    }
    return (res, col_index, row_index);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        assert_eq!(
            Solution::minimum_total_distance(vec![0, 4, 6], vec![vec![2, 2], vec![6, 2]]),
            4
        )
    }
}
