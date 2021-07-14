"use strict";
// 简述instanceof 的原理并用代码实现
// 如果A沿着原型链能找到B.prototype,那么 A instanceof B 为 true
// 解法:遍历A的原型链，找到B.prototype则返回true
const _isinstanceof = (input, target) => {
    let n = input;
    while (n) {
        if (n === target.prototype) {
            return true;
        }
        n = n.__proto__;
    }
    return false;
};
console.log(_isinstanceof(1, Object));
