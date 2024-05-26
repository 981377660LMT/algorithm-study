macro_rules! implement {
  ($tr:ty; $($t:ty),*) => {
      $(
          impl $tr for $t { }
      )*
  }
}

/// 非負の数値型
pub trait Unsigned {}
implement!(Unsigned; u8, u16, u32, u64, u128, usize);

/// 符号付きの数値型
pub trait Signed {}
implement!(Signed; i8, i16, i32, i64, i128, isize, f32, f64);

/// 整数型
pub trait Int {}
implement!(Int; u8, u16, u32, u64, u128, usize, i8, i16, i32, i64, i128, isize);

/// 浮動小数点型
pub trait Float {}
implement!(Float; f32, f64);

#[cfg(test)]
mod tests {

    #[test]
    fn test_unsigned() {
        let _: u8 = 0;
        let _: u16 = 0;
        let _: u32 = 0;
        let _: u64 = 0;
        let _: u128 = 0;
        let _: usize = 0;
    }

    #[test]
    fn test_signed() {
        let _: i8 = 0;
        let _: i16 = 0;
        let _: i32 = 0;
        let _: i64 = 0;
        let _: i128 = 0;
        let _: isize = 0;
        let _: f32 = 0.0;
        let _: f64 = 0.0;
    }

    #[test]
    fn test_int() {
        let _: u8 = 0;
        let _: u16 = 0;
        let _: u32 = 0;
        let _: u64 = 0;
        let _: u128 = 0;
        let _: usize = 0;
        let _: i8 = 0;
        let _: i16 = 0;
        let _: i32 = 0;
        let _: i64 = 0;
        let _: i128 = 0;
        let _: isize = 0;
    }

    #[test]
    fn test_float() {
        let _: f32 = 0.0;
        let _: f64 = 0.0;
    }
}
