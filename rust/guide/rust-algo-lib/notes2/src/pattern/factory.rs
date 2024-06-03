// 工厂方法定义了一个方法， 且必须使用该方法代替通过直接调用构造函数来创建对象 （ new操作符） 的方式。
// 子类可重写该方法来更改将被创建的对象所属类。
// !在父类中提供一个创建对象的接口以允许子类决定实例化对象的类型

mod maze_game {
    trait IRoom {
        fn render(&self);
    }

    // 迷宫拥有工厂方法创建房间.
    // static dispatch (generics).
    trait IMazeGame {
        type RoomImpl: IRoom;

        // !factory method.
        // It must be overridden with a concrete implementation.
        fn create_rooms(&self) -> Vec<Self::RoomImpl>;

        fn play(&self) {
            for room in self.create_rooms() {
                room.render();
            }
        }
    }

    #[derive(Clone)]
    struct MagicRoom {
        title: String,
    }

    impl MagicRoom {
        pub fn new(title: String) -> Self {
            Self { title }
        }
    }

    impl IRoom for MagicRoom {
        fn render(&self) {
            println!("MagicRoom: {}", self.title);
        }
    }

    pub struct MagicMaze {
        rooms: Vec<MagicRoom>,
    }

    impl MagicMaze {
        pub fn new() -> Self {
            Self {
                rooms: vec![
                    MagicRoom::new("Infinite Room".into()),
                    MagicRoom::new("Invisible Room".into()),
                ],
            }
        }
    }

    impl IMazeGame for MagicMaze {
        type RoomImpl = MagicRoom;

        fn create_rooms(&self) -> Vec<Self::RoomImpl> {
            self.rooms.clone()
        }
    }

    pub fn run(maze_game: impl IMazeGame) {
        println!("Loading resources...");
        println!("Starting the game...");
        maze_game.play();
    }
}

fn main() {
    let magic_maze = maze_game::MagicMaze::new();
    maze_game::run(magic_maze);
}
