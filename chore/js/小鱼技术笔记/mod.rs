#[derive(Debug)]
struct User {
    user_id: i32,
    sex: i32,
}

fn main() {
    //'a作用域

    let user_ref: &User;
    {
        //'b作用域
        let user = User {
            user_id: 10,
            sex: 2,
        };
        user_ref = &user;
        //正确，可以运行，解引用的地方
        //所有者是user，作用域在'b，解引用的地方是'b，所有者的作用域'b大于等于解引用的作用域'b。
        println!("{:?}", user_ref)
    }

    //错误，引用所指向的所有权已经丢失了，无法读取
    //所有者是user，作用域在'b，解引用的地方是'a，所有者的作用域少于解引用的作用域'a
    println!("{:?}", user_ref)
}
