// !桥接模式的本质：将变化的部分抽象成接口，将不变的部分抽象成类，通过组合的方式将两者组合在一起。
//
// 桥接是一种结构型设计模式，
// !可将业务逻辑或一个大类拆分为不同的层次结构， 从而能独立地进行开发。

// 场景：
// 代码有两层，第一层是引用者(Abstaction)，第二层是被引用者(Implementor)
// 引用者通过被引用者实现自己的方法（将能将一些 （有时是绝大部分） 对自己的调用委派给被引用者）。
// !所有的被引用者都实现了相同的接口（非常关键!!!)。
//
// 桥接模式：
// !抽象部分要与它的实现部分分离：把通用能力抽出来，作为"桥"，用这个桥的能力实现不同的 Abstraction.
// !不要把方法写在某个固定的class上，而是把方法抽出来，更容易作为桥维护.
// 即：Implementor 提供更底层的方法，供引用者调用，引用者可以根据自己的需要实现自己的方法。
// Implementor 与 Abstraction 分离，使它们都可以独立变化。
//
// 例子：遥控器（remotes）操纵设备(device)。

mod remote {
    use super::device::IDevice;

    pub trait IDeviceHolder<D: IDevice> {
        fn device(&mut self) -> &mut D;
    }

    pub trait IRemote<D: IDevice>: IDeviceHolder<D> {
        fn power(&mut self) {
            println!("Remote: power toggle");
            if self.device().is_enabled() {
                self.device().disable();
            } else {
                self.device().enable();
            }
        }
    }

    pub struct BasicRemote<D: IDevice> {
        device: D,
    }

    impl<D: IDevice> BasicRemote<D> {
        pub fn new(device: D) -> Self {
            Self { device }
        }
    }

    impl<D: IDevice> IDeviceHolder<D> for BasicRemote<D> {
        fn device(&mut self) -> &mut D {
            &mut self.device
        }
    }

    impl<D: IDevice> IRemote<D> for BasicRemote<D> {}

    pub struct AdvancedRemote<D: IDevice> {
        device: D,
    }

    impl<D: IDevice> AdvancedRemote<D> {
        pub fn new(device: D) -> Self {
            Self { device }
        }

        pub fn mute(&mut self) {
            println!("Remote: mute");
            self.device.disable();
        }
    }

    impl<D: IDevice> IDeviceHolder<D> for AdvancedRemote<D> {
        fn device(&mut self) -> &mut D {
            &mut self.device
        }
    }

    impl<D: IDevice> IRemote<D> for AdvancedRemote<D> {}
}

mod device {
    pub trait IDevice {
        fn is_enabled(&self) -> bool;
        fn enable(&mut self);
        fn disable(&mut self);
        fn run(&self);
    }

    #[derive(Clone, Debug, Default)]
    pub struct Radio {
        on: bool,
    }

    impl IDevice for Radio {
        fn is_enabled(&self) -> bool {
            self.on
        }

        fn enable(&mut self) {
            self.on = true;
        }

        fn disable(&mut self) {
            self.on = false;
        }

        fn run(&self) {
            println!("Radio is running.");
        }
    }

    #[derive(Clone, Debug, Default)]
    pub struct TV {
        on: bool,
    }

    impl IDevice for TV {
        fn is_enabled(&self) -> bool {
            self.on
        }

        fn enable(&mut self) {
            self.on = true;
        }

        fn disable(&mut self) {
            self.on = false;
        }

        fn run(&self) {
            println!("TV is running.");
        }
    }
}

fn main() {
    use self::device::{IDevice, Radio, TV};
    use self::remote::{AdvancedRemote, BasicRemote, IDeviceHolder, IRemote};

    test_device(Radio::default());
    test_device(TV::default());

    fn test_device(device: impl IDevice + Clone) {
        let mut basic_remote = BasicRemote::new(device.clone());
        basic_remote.power();
        basic_remote.device().run();

        let mut advanced_remote = AdvancedRemote::new(device.clone());
        advanced_remote.power();
        advanced_remote.mute();
        advanced_remote.device().run();
    }
}
