use std::fmt;

const DAY: i64 = 24 * 60;
const HOUR: i64 = 60;

#[derive(Debug, PartialEq)]
pub struct Clock {
    minutes: i64,
}

impl Clock {
    pub fn new(hours: i64, minutes: i64) -> Self {
        Clock {
            minutes: (hours * HOUR + minutes).rem_euclid(DAY),
        }
    }

    pub fn add_minutes(&self, minutes: i64) -> Self {
        Clock {
            minutes: (self.minutes + minutes).rem_euclid(DAY),
        }
    }
}

// to_string() is a trait method of Display.
// cannot be derived.
impl std::fmt::Display for Clock {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:02}:{:02}", self.minutes / HOUR, self.minutes % HOUR)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        println!("{}", -100 % 24);
    }
}
