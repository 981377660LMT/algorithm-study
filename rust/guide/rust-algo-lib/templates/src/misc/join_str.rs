pub trait JoinStr {
    fn join_str(self, _: &str) -> String;
}

impl<I: Iterator<Item = T>, T: ToString> JoinStr for I {
    fn join_str(self, s: &str) -> String {
        self.map(|x| x.to_string()).collect::<Vec<_>>().join(s)
    }
}

#[cfg(test)]
mod tests {
    use super::JoinStr;

    #[test]
    fn test_join_str() {
        let v = vec![1, 2, 3];
        assert_eq!(v.iter().join_str(","), "1,2,3");
    }
}
