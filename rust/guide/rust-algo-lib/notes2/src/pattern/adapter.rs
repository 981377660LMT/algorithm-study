// 适配器是一种结构型设计模式， 它能使不兼容的对象能够相互合作。
// 设计三个模块：适配器、适配者、目标
//
// !Adapter 通过实现 Target 接口，将 Adaptee 转换为 Target 接口的形式

mod adaptee {
    pub struct OldApi;
    impl OldApi {
        pub fn old_request(&self) -> String {
            ".roivaheb s'eetpadA ehT :gnihtemoS".into()
        }
    }
}

mod target {
    pub trait ITarget {
        fn request(&self) -> String;
    }

    pub struct OrdinaryTarget;
    impl ITarget for OrdinaryTarget {
        fn request(&self) -> String {
            "OrdinaryTarget: The default target's behavior.".into()
        }
    }

    pub fn call(target: impl ITarget) {
        println!("{}", target.request());
    }
}

#[allow(clippy::module_inception)]
mod adapter {
    use super::adaptee::OldApi;
    use super::target::ITarget;

    // 将适配者转换为`Target`接口.
    pub struct TargetAdapter {
        adaptee: OldApi,
    }

    impl TargetAdapter {
        pub fn new(adaptee: OldApi) -> Self {
            Self { adaptee }
        }
    }

    impl ITarget for TargetAdapter {
        fn request(&self) -> String {
            // 适配器将适配者的行为转换为目标接口的行为.
            self.adaptee.old_request().chars().rev().collect()
        }
    }
}

fn main() {
    use adaptee::OldApi;
    use adapter::TargetAdapter;
    use target::{call as run, OrdinaryTarget};

    let target = OrdinaryTarget;
    println!("Client: I can work just fine with the Target objects:");
    run(target);

    let adaptee = OldApi;
    println!("Client: The Adaptee class has a weird interface. See, I don't understand it:");
    let adapter = TargetAdapter::new(adaptee);
    println!("Client: But I can work with it via the Adapter:");
    run(adapter);
}
