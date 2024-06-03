// !抽象工厂模式总结：抽象工厂是工厂的接口.
// 将实际的创建工作留给了具体工厂类
// 在创建产品时， 客户端代码调用的是工厂对象的构建方法， 而不是直接调用构造函数.

// GUI Elements Factory
mod ui {
    pub trait IButton {
        fn press(&self);
    }

    pub trait ICheckbox {
        fn check(&self);
    }

    // Static Dispatch
    #[allow(dead_code)]
    trait IFactoryStatic {
        type B: IButton;
        type C: ICheckbox;
        fn create_button(&self) -> Self::B;
        fn create_checkbox(&self) -> Self::C;
    }

    // !抽象工厂(接口)
    // Dynamic Dispatch
    // 如果在编译时不知道抽象工厂的具体类型，则应使用 Box 指针实现。
    //
    // 当您在静态调度和动态调度之间做出选择时，很少有明确的正确答案。
    // !您需要在库中使用静态调度(因为你知道具体类型)，在二进制文件中使用动态调度(因为你不知道具体类型)。
    pub trait IFactoryDynamic {
        fn create_button(&self) -> Box<dyn IButton>;
        fn create_checkbox(&self) -> Box<dyn ICheckbox>;
    }
}

mod dom {
    use super::ui::{IButton, ICheckbox, IFactoryDynamic};

    pub struct DomButton;
    impl IButton for DomButton {
        fn press(&self) {
            println!("DomButton pressed");
        }
    }

    pub struct DomCheckbox;
    impl ICheckbox for DomCheckbox {
        fn check(&self) {
            println!("DomCheckbox checked");
        }
    }

    pub struct DomFactory;
    impl IFactoryDynamic for DomFactory {
        fn create_button(&self) -> Box<dyn IButton> {
            Box::new(DomButton)
        }

        fn create_checkbox(&self) -> Box<dyn ICheckbox> {
            Box::new(DomCheckbox)
        }
    }
}

mod canvas {
    use super::ui::{IButton, ICheckbox, IFactoryDynamic};

    pub struct CanvasButton;
    impl IButton for CanvasButton {
        fn press(&self) {
            println!("CanvasButton pressed");
        }
    }

    pub struct CanvasCheckbox;
    impl ICheckbox for CanvasCheckbox {
        fn check(&self) {
            println!("CanvasCheckbox checked");
        }
    }

    pub struct CanvasFactory;
    impl IFactoryDynamic for CanvasFactory {
        fn create_button(&self) -> Box<dyn IButton> {
            Box::new(CanvasButton)
        }

        fn create_checkbox(&self) -> Box<dyn ICheckbox> {
            Box::new(CanvasCheckbox)
        }
    }
}

use self::{canvas::CanvasFactory, dom::DomFactory, ui::IFactoryDynamic};

// Client code with dynamic dispatch.
fn main() {
    // Allocate a factory object in runtime depending on unpredictable input.

    let factory = init();
    render(factory);

    fn init() -> &'static dyn IFactoryDynamic {
        if cfg!(use_dome) {
            &DomFactory
        } else {
            &CanvasFactory
        }
    }

    // Renders a GUI by the given factory.
    // The code demonstrates that it doesn't depend on a concrete factory implementation.
    fn render(factory: &dyn IFactoryDynamic) {
        let button1 = factory.create_button();
        let button2 = factory.create_button();
        let checkbox1 = factory.create_checkbox();
        let checkbox2 = factory.create_checkbox();

        button1.press();
        button2.press();
        checkbox1.check();
        checkbox2.check();
    }
}
