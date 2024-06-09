use std::fmt;

struct Password(String);

impl fmt::Display for Password {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Password(****)")
    }
}

fn main() {
    let password = Password("hunter2".to_string());
}
