为什么不能在循环、判断内部使用 Hook
useState 的实现原理
useState 是基于 Array+Cursor 来实现的,state 和 cursor 必须有固定的的对应关系

第一次渲染时候，根据 useState 顺序，逐个声明 state 并且将其放入全局 Array 中。每次声明 state，都要将 cursor 增加 1。
更新 state，触发再次渲染的时候。cursor 被重置为 0。按照 useState 的声明顺序，依次拿出最新的 state 的值，视图更新。

```TS
import React from "react";
import ReactDOM from "react-dom";
const states: any[] = [];
let cursor: number = 0;
function useState<T>(initialState: T): [T, (newState: T) => void] {
    const currenCursor = cursor;
    states[currenCursor] = states[currenCursor] || initialState; // 检查是否渲染过
    function setState(newState: T) {
        states[currenCursor] = newState;
        render();
    }
    ++cursor; // update: cursor
    return [states[currenCursor], setState];
}
function App() {
    const [num, setNum] = useState < number > 0;
    const [num2, setNum2] = useState < number > 1;
    return (
        <div>
            <div>num: {num}</div>
            <div>
                <button onClick={() => setNum(num + 1)}>加 1</button>
                <button onClick={() => setNum(num - 1)}>减 1</button>
            </div>
            <hr />
            <div>num2: {num2}</div>
            <div>
                <button onClick={() => setNum2(num2 * 2)}>扩大一倍</button>
                <button onClick={() => setNum2(num2 / 2)}>缩小一倍</button>
            </div>
        </div>
    );
}
function render() {
    ReactDOM.render(<App />, document.getElementById("root"));
    cursor = 0; // 重置cursor
}
render(); // 首次渲染
```

useEffect 的实现原理

react 路由配置
https://xin-tan.com/2019-09-11-react-router/

```TS
import { Route, Switch, SwitchProps, RouteProps } from "react-router-dom";
function renderRoutes(params: {
    routes: RouteProps[];
    switchProps?: SwitchProps;
}) {
    const { switchProps, routes } = params;
    return (
        <Switch {...switchProps}>
            {routes.map((route, index) => (
                <Route
                    key={index}
                    path={route.path}
                    component={route.component}
                    exact={route.exact || true}
                    strict={route.strict || false}
                ></Route>
            ))}
        </Switch>
    );
}

import { RouteProps } from "react-router-dom";
const config: RouteProps[] = [
    {
        path: "/",
        component: HomePage
    },
    {
        path: "/user",
        component: UserPage
    }
];

import React, { Component } from "react";
import { BrowserRouter } from "react-router-dom";
const routes = renderRoutes({
    routes: config
});
class App extends Component {
    render() {
        return <BrowserRouter>{routes}</BrowserRouter>;
    }
}
export default App;
```

React 合成事件的好处与不足
跨浏览器兼容
统一管理,内存更优
支持事件类型不如浏览器本身多

目前微服务理念非常火，后端架构都像无状态服务转变，这方便基于 k8s 的横向扩容，应对突发流量。
但是在前端开发中，尤其是控制台这种业务逻辑很重、交互细节繁多的场景下，都是基于数据状态来渲染视图。
**Redux 状态管理**
对于一个 ajax 请求，它有 4 种状态：
未发送
发送中
收到结果
成功返回
出错失败
