// !代理模式的优点在于它可以在不修改被代理对象的代码的情况下，对被代理对象的功能进行扩展。
// 代理对象拥有和服务对象相同的接口， 这使得当其被传递给客户端时可与真实对象互换。
//
// 两个模块: 代理对象(Proxy)和服务对象(Subject)。
// 代理对象拥有一个指向服务对象的引用，当客户端调用代理对象的方法时，代理对象会将请求转发给服务对象。
//
// 例子：nginx 代理 server，可以进行 rate limiting、request caching 等操作。

mod subject {
    pub trait IServer {
        fn handle_request(&mut self, url: &str, method: &str) -> (u16, String);
    }

    pub struct Server;
    impl IServer for Server {
        fn handle_request(&mut self, url: &str, method: &str) -> (u16, String) {
            if url == "/app/status" && method == "GET" {
                (200, "OK".into())
            } else {
                (404, "Not Found".into())
            }
        }
    }
}

#[allow(clippy::module_inception)]
mod proxy {
    use std::collections::HashMap;

    use super::subject::{IServer, Server};

    /// NGINX server is a proxy to an application server.
    pub struct NginxServer {
        server: Server,
        max_allowed_requests: u32,
        rate_limiter: HashMap<String, u32>,
    }

    impl NginxServer {
        pub fn new() -> Self {
            Self {
                server: Server,
                max_allowed_requests: 2,
                rate_limiter: HashMap::default(),
            }
        }

        pub fn check_rate_limiting(&mut self, url: &str) -> bool {
            let rate = self.rate_limiter.entry(url.into()).or_insert(1);
            if *rate > self.max_allowed_requests {
                false
            } else {
                *rate += 1;
                true
            }
        }
    }

    impl IServer for NginxServer {
        fn handle_request(&mut self, url: &str, method: &str) -> (u16, String) {
            if !self.check_rate_limiting(url) {
                return (403, "Forbidden".into());
            }
            self.server.handle_request(url, method)
        }
    }
}

fn main() {
    use proxy::NginxServer;
    use subject::IServer;
    let app_status = &"/app/status".to_string();
    let create_user = &"/user/create".to_string();

    let mut nginx = NginxServer::new();
    let (status, message) = nginx.handle_request(app_status, "GET");
    println!("Status: {}, Message: {}", status, message);
    let (status, message) = nginx.handle_request(app_status, "GET");
    println!("Status: {}, Message: {}", status, message);
    let (status, message) = nginx.handle_request(app_status, "GET");
    println!("Status: {}, Message: {}", status, message);
    let (status, message) = nginx.handle_request(create_user, "POST");
    println!("Status: {}, Message: {}", status, message);
    let (status, message) = nginx.handle_request(app_status, "GET");
    println!("Status: {}, Message: {}", status, message);
}
