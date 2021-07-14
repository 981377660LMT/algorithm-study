"use strict";
// 不要用forEach!
// 除了抛出异常以外，没有办法中止或跳出 forEach() 循环。
// 如果你需要中止或跳出循环，forEach() 方法不是应当使用的工具。
Object.defineProperty(exports, "__esModule", { value: true });
// 婚姻介绍所找对象登记到Map
const twoSum = (arr, target) => {
    const map = new Map();
    for (let index = 0; index < arr.length; index++) {
        const element = arr[index];
        if (map.has(element)) {
            return [map.get(element), index];
        }
        else {
            const matchValue = target - element;
            map.set(matchValue, index);
        }
    }
};
console.log(twoSum([2, 7, 11, 15], 9));
