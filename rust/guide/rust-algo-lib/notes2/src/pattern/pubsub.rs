// !发布订阅模式：允许一个对象将其状态的改变通知其他对象。
//
// 观察者模式和发布订阅模式的区别：
//
// !1.是否互相感知(第三者)：
//   在观察者模式中，观察者是知道Subject的，Subject一直保持对观察者进行记录。
//   然而，在发布订阅模式中，发布者和订阅者不知道对方的存在。它们只有通过消息代理进行通信。
// !2.耦合度：
//   在发布订阅模式中，组件是松散耦合的，正好和观察者模式相反。
// !3.同步/异步：
//   观察者模式大多数时候是同步的，比如当事件触发，Subject就会去调用观察者的方法。
//   而发布-订阅模式大多数时候是异步的（使用消息队列）

/// pub-sub pattern.
mod event_emitter {
    use std::collections::HashMap;

    #[derive(PartialEq, Eq, Hash, Clone)]
    pub enum Event {
        Load,
        Save,
    }

    pub type Callback = fn(file_path: String);

    #[derive(Default)]
    pub struct EventEmitter {
        events: HashMap<Event, Vec<Callback>>,
    }

    impl EventEmitter {
        pub fn on(&mut self, event: Event, callback: Callback) {
            self.events.entry(event).or_default().push(callback);
        }

        pub fn off(&mut self, event: Event, callback: Callback) {
            if let Some(callbacks) = self.events.get_mut(&event) {
                callbacks.retain(|&c| c != callback);
            }
        }

        pub fn emit(&self, event: Event, file_path: String) {
            if let Some(callbacks) = self.events.get(&event) {
                for f in callbacks {
                    f(file_path.clone());
                }
            }
        }
    }
}

mod editor {
    use super::event_emitter::{Event, EventEmitter};

    #[derive(Default)]
    pub struct Editor {
        event_emitter: EventEmitter,
        file_path: String,
    }

    impl Editor {
        pub fn events(&mut self) -> &mut EventEmitter {
            &mut self.event_emitter
        }

        pub fn load(&mut self, file_path: String) {
            self.file_path.clone_from(&file_path);
            self.event_emitter.emit(Event::Load, file_path);
        }

        pub fn save(&self) {
            self.event_emitter.emit(Event::Save, self.file_path.clone());
        }
    }
}

fn main() {
    use editor::Editor;
    use event_emitter::Event;

    let mut editor = Editor::default();

    editor.events().on(Event::Load, |file_path| {
        println!("File loaded: {}", file_path);
    });
    editor.events().on(Event::Save, save_listener);

    editor.load("example.txt".into());
    editor.load("example2.txt".into());
    editor.save();

    editor.events().off(Event::Save, save_listener);
    editor.save();

    fn save_listener(file_path: String) {
        let email = "foo".to_string();
        println!("File saved: {} and email sent to {}", file_path, email);
    }
}
