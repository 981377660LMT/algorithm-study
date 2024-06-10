use std::{
    marker::PhantomData,
    ops::{Add, Sub},
};

#[derive(Debug, Copy, Clone)]
struct Meter {
    value: f64,
}

impl Meter {
    fn new(value: f64) -> Self {
        Self { value }
    }
}

impl Add for Meter {
    type Output = Self;

    fn add(self, rhs: Self) -> Self::Output {
        Meter {
            value: self.value + rhs.value,
        }
    }
}

impl Sub for Meter {
    type Output = Self;

    fn sub(self, rhs: Self) -> Self::Output {
        Meter {
            value: self.value - rhs.value,
        }
    }
}

// 如果此时，我们还需要为Kilogram、Liter等类型实现相同的逻辑，则需要重复的实现多个Add和Sub Trait；
// 虽然我们可以使用过程宏（Macro）实现.

#[derive(Debug, Copy, Clone)]
struct Unit<T> {
    value: f64,
    unit_type: PhantomData<T>,
}

impl<T> Unit<T> {
    fn new(value: f64) -> Self {
        Self {
            value,
            unit_type: PhantomData,
        }
    }
}

impl<T> Add for Unit<T> {
    type Output = Self;

    fn add(self, rhs: Self) -> Self::Output {
        Unit::new(self.value + rhs.value)
    }
}

impl<T> Sub for Unit<T> {
    type Output = Self;

    fn sub(self, rhs: Self) -> Self::Output {
        Unit::new(self.value - rhs.value)
    }
}

#[derive(Debug, Copy, Clone)]
struct MeterType;
type Meter2 = Unit<MeterType>;

#[derive(Debug, Copy, Clone)]
struct KilogramType;
type Kilogram = Unit<KilogramType>;

// 得益于PhantomData：

// !泛型占位.
// 我们直接使用PhantomData对unit_type进行了初始化，而无需为各种类型实现Default Trait；
// !同时，使用 PhantomData<T>，我们可以将类型参数的用途明确地在代码中传达；

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        let one = Meter::new(1.1);
        let two = Meter::new(2.2);

        let four = one + two;
        dbg!(four);

        let zero = four - two;
        dbg!(zero);

        let one = Meter2::new(1.1);
        let two = Meter2::new(2.2);

        let four = one + two;
        dbg!(four);

        let zero = four - two;
        dbg!(zero);
    }
}
