// https://leetcode.cn/problems/differences-between-two-objects/

// 请你编写一个函数，它接收两个深度嵌套的对象或数组 obj1 和 obj2 ，并返回一个新对象表示它们之间差异。
// 该函数应该比较这两个对象的属性，并识别任何变化。
// 返回的对象应仅包含从 obj1 到 obj2 的值不同的键。
// 对于每个变化的键，值应表示为一个数组 [obj1 value, obj2 value] 。
// 不存在于一个对象中但存在于另一个对象中的键不应包含在返回的对象中。
// 在比较两个数组时，数组的索引被视为它们的键。
// 最终结果应是一个深度嵌套的对象，其中每个叶子的值都是一个差异数组。
// 你可以假设这两个对象都是 JSON.parse 的输出结果。

// !1. obj1 和 obj2 的类型不同
// !2. obj1 和 obj2 的类型相同 => 基本类型直接返回/非基本类型递归处理

function objDiff(obj1: any, obj2: any): any {}

export {}

// const isPrimitive = o => o == null || typeof o != "object";
// const { toString } = Object.prototype;
// /**
//  * @param {object} obj1
//  * @param {object} obj2
//  * @return {object}
//  */
// function objDiff(obj1, obj2) {
//     if (toString.call(obj1) != toString.call(obj2)) return [obj1, obj2];
//     if (isPrimitive(obj1)) return obj1 === obj2 ? {} : [obj1, obj2];
//     const res = {};
//     for (let k of Object.keys(obj1).filter(k => k in obj2)) {
//         let t = objDiff(obj1[k], obj2[k]);
//         if (Object.keys(t).length) res[k] = t;
//     }
//     return res;
// };

console.log(Object.prototype.toString.call([]))

function getType(obj: unknown): string {
  return Object.prototype.toString.call(obj).slice(8, -1)
}

function isObject(obj: unknown): obj is object {
  return typeof obj === 'object' && obj !== null
}
