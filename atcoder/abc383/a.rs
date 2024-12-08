use proconio::{fastout, input};
use std::cmp::max;
use std::collections::HashMap;

#[fastout]
fn main() {
    input! {
        N: usize,
        X: usize,
        K: i64,
        items: [(usize, i64, usize); N],
    }

    let mut color_groups: HashMap<usize, Vec<(usize, i64)>> = HashMap::new();
    for (P, U, C) in items {
        color_groups.entry(C).or_insert_with(Vec::new).push((P, U));
    }

    let INF_NEG: i64 = -1_000_000_000_000_000;

    let mut color_data = Vec::new();

    for (_, group_items) in color_groups {
        let mut dp_color = vec![INF_NEG; X + 1];
        dp_color[0] = 0;
        for &(p, u) in &group_items {
            for cost in (0..=X - p).rev() {
                if dp_color[cost] != INF_NEG {
                    let new_cost = cost + p;
                    if new_cost <= X {
                        dp_color[new_cost] = max(dp_color[new_cost], dp_color[cost] + u);
                    }
                }
            }
        }
        color_data.push(dp_color);
    }

    let mut dp = vec![INF_NEG; X + 1];
    dp[0] = 0;

    for dp_color in color_data {
        let mut new_dp = vec![INF_NEG; X + 1];
        for cost in 0..=X {
            if dp[cost] != INF_NEG {
                new_dp[cost] = max(new_dp[cost], dp[cost]);
            }
        }
        for cost in 0..=X {
            if dp[cost] == INF_NEG {
                continue;
            }
            for add_cost in 1..=X - cost {
                let val = dp_color[add_cost];
                if val != INF_NEG {
                    let nc = cost + add_cost;
                    if nc <= X {
                        new_dp[nc] = max(new_dp[nc], dp[cost] + val + K);
                    }
                }
            }
        }
        dp = new_dp;
    }

    let ans = dp.into_iter().max().unwrap();
    println!("{}", ans);
}
