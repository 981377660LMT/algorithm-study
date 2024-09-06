mod clock;

pub fn is_leap_year(year: i32) -> bool {
    let has_factor = |n| year % n == 0;
    has_factor(4) && (!has_factor(100) || has_factor(400))
}

pub fn raindrops(n: u32) -> String {
    let mut res = String::new();
    if n % 3 == 0 {
        res.push_str("Pling");
    }
    if n % 5 == 0 {
        res.push_str("Plang");
    }
    if n % 7 == 0 {
        res.push_str("Plong");
    }
    if res.is_empty() {
        res = n.to_string();
    }
    res
}

// kth_prime
pub fn nth(n: u32) -> u32 {
    if n == 0 {
        panic!("There is no zeroth prime.");
    }
    let mut primes = vec![2];
    let mut candidate = 3;
    while primes.len() < n as usize {
        if primes.iter().all(|&prime| candidate % prime != 0) {
            primes.push(candidate);
        }
        candidate += 2;
    }
    primes[n as usize - 1]
}

pub fn reverse(input: &str) -> String {
    input.chars().rev().collect()
}

use time::PrimitiveDateTime as DateTime;

// Returns a DateTime one billion seconds after start.
pub fn after(start: DateTime) -> DateTime {
    start + time::Duration::seconds(1_000_000_000)
}
