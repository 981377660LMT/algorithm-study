#![allow(clippy::new_without_default)]

#[derive(Debug, Clone)]
pub struct XorShift {
    x: u64,
    y: u64,
    z: u64,
    w: u64,
}

impl XorShift {
    pub fn new() -> XorShift {
        XorShift {
            x: 123456789,
            y: 362436069,
            z: 521288629,
            w: 88675123,
        }
    }

    pub fn gen(&mut self) -> u64 {
        let t = self.x ^ (self.x << 11);
        self.x = self.y;
        self.y = self.z;
        self.z = self.w;
        self.w = (self.w ^ (self.w >> 19)) ^ (t ^ (t >> 8));
        self.w
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_xor_shift() {
        let mut xor_shift = XorShift::new();
        for _ in 0..100 {
            let _ = xor_shift.gen();
        }
    }
}
