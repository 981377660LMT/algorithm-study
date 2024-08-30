struct Solution;

impl Solution {
    pub fn buddy_strings(s: String, goal: String) -> bool {
        if s.len() != goal.len() {
            return false;
        }
        let diff = s
            .bytes()
            .zip(goal.bytes())
            .filter(|u| u.0 != u.1)
            .collect::<Vec<_>>();

        match diff.len() {
            2 => diff[0].0 == diff[1].1 && diff[0].1 == diff[1].0,
            0 => s.bytes().collect::<std::collections::HashSet<_>>().len() < s.len(),
            _ => false,
        }
    }
}
