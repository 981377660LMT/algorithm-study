// #[inline]：给编译器一个提示，建议它内联这个函数，但编译器可以选择忽略这个建议。
// #[inline(always)]：强烈建议编译器总是内联这个函数，除非技术上不可行，编译器通常会遵循这个指示。
// #[inline(never)]：告诉编译器不要内联这个函数。

struct ParkingSystem([i16; 4]);

impl ParkingSystem {
    #[inline(always)]
    fn new(big: i32, medium: i32, small: i32) -> Self {
        Self([0, big as i16, medium as i16, small as i16])
    }

    #[inline(always)]
    fn add_car(&mut self, car_type: i32) -> bool {
        unsafe {
            let x = self.0.get_unchecked_mut(car_type as usize);
            *x -= 1;
            *x >= 0
        }
    }
}
