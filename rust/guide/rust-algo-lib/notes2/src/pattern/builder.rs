// Builder 设计模式与 Fluent Interface 惯用语（链式调用）不同，尽管 Rust 开发人员有时会互换使用这些术语。
// !Fluent Interface 指的是一种编程风格，其中方法返回 self 或者 &mut self，
// ```rust
// let car = Car::default().places(5).gas(30)
// ```
//
// !Builder 模式中有三个角色：Builder、Product 和 Director。
// Builder 是一个 trait，它定义了构建 Product 所需的方法。
// Product 是一个结构体，它是最终构建的对象。
// Director 是一个结构体，它指挥 Builder 如何构建 Product。
//
// !工头（Director）告诉工人（Builder）如何建造房子（Product）。
//
// !Builder 模式的目的是将构建复杂对象的过程与表示分离。
// !这样，您可以使用相同的构建过程构建不同的表示。

#![allow(unused)]

mod car {
    pub mod director {
        // 指挥者知道如何构建产品.
        pub struct Director;
        use super::builder::IBuilder;
        use super::product::CarType;

        impl Director {
            pub fn construct_sports_car(builder: &mut impl IBuilder) {
                builder.set_car_type(CarType::SportsCar);
                builder.set_seats(2);
            }

            pub fn construct_city_car(builder: &mut impl IBuilder) {
                builder.set_car_type(CarType::CityCar);
                builder.set_seats(4);
            }

            pub fn construct_suv(builder: &mut impl IBuilder) {
                builder.set_car_type(CarType::Suv);
                builder.set_seats(4);
            }
        }
    }

    // #[allow(clippy::module_inception)]
    pub mod builder {
        use super::product::{Car, CarManual, CarType};

        pub trait IBuilder {
            type Output;
            fn set_car_type(&mut self, car_type: CarType);
            fn set_seats(&mut self, seats: u16);
            fn build(self) -> Self::Output;
        }

        #[derive(Default)]
        pub struct CarBuilder {
            type_: Option<CarType>,
            seats: Option<u16>,
        }

        impl IBuilder for CarBuilder {
            type Output = Car;

            fn set_car_type(&mut self, car_type: CarType) {
                self.type_ = Some(car_type);
            }

            fn set_seats(&mut self, seats: u16) {
                self.seats = Some(seats);
            }

            fn build(self) -> Self::Output {
                Car::new(
                    self.type_.expect("Please set car type"),
                    self.seats.expect("Please set seats"),
                )
            }
        }

        #[derive(Default)]
        pub struct CarManualBuilder {
            type_: Option<CarType>,
            seats: Option<u16>,
        }

        impl IBuilder for CarManualBuilder {
            type Output = CarManual;

            fn set_car_type(&mut self, car_type: CarType) {
                self.type_ = Some(car_type);
            }

            fn set_seats(&mut self, seats: u16) {
                self.seats = Some(seats);
            }

            fn build(self) -> Self::Output {
                CarManual::new(
                    self.type_.expect("Please set car type"),
                    self.seats.expect("Please set seats"),
                )
            }
        }
    }

    pub mod product {
        pub struct Car {
            type_: CarType,
            seats: u16,
        }

        impl Car {
            pub fn new(car_type: CarType, seats: u16) -> Self {
                Self {
                    seats,
                    type_: car_type,
                }
            }

            pub fn default() -> Self {
                Self {
                    seats: 4,
                    type_: CarType::CityCar,
                }
            }

            pub fn seats(&self) -> u16 {
                self.seats
            }

            pub fn car_type(&self) -> CarType {
                self.type_
            }
        }

        pub struct CarManual {
            car_type: CarType,
            seats: u16,
        }

        impl CarManual {
            pub fn new(car_type: CarType, seats: u16) -> Self {
                Self { car_type, seats }
            }
        }

        impl std::fmt::Display for CarManual {
            fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                writeln!(f, "Car type: {:?}", self.car_type)?;
                writeln!(f, "Car manual: {} seats", self.seats)?;
                Ok(())
            }
        }

        #[derive(Debug, Clone, Copy)]
        pub enum CarType {
            CityCar,
            SportsCar,
            Suv,
        }
    }
}

fn main() {
    // !注意需要导入IBuilder后才能使用build方法
    use car::builder::{CarBuilder, CarManualBuilder, IBuilder};
    use car::director::Director;
    use car::product::{Car, CarManual};

    let mut car_builder = CarBuilder::default();
    // Director 接收客户端builder并告诉它如何构建产品.
    Director::construct_sports_car(&mut car_builder);
    let car: Car = car_builder.build();
    println!("Car type: {:#?}", car.car_type());

    let mut car_manual_builder = CarManualBuilder::default();
    Director::construct_suv(&mut car_manual_builder);
    let car_manual: CarManual = car_manual_builder.build();
    println!("{}", car_manual);
}
