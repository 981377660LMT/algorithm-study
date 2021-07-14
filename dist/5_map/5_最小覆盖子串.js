"use strict";
// 在字符串s里找出包含t所有字符的最小字串
Object.defineProperty(exports, "__esModule", { value: true });
// 思路:
// 用滑动窗口找出所有包含t的字串，返回长度最小的
// 当找到子串时，移动左指针继续
const minWindow = (s, t) => {
    let leftPoint = 0;
    const needMap = new Map();
    for (const alphabet of t) {
        needMap.set(alphabet, needMap.has(alphabet) ? needMap.get(alphabet) + 1 : 1);
    }
    let needType = needMap.size;
    const matchedStringMemo = [];
    for (let rightPoint = 0; rightPoint < s.length; rightPoint++) {
        const element = s[rightPoint];
        if (needMap.has(element)) {
            needMap.set(element, needMap.get(element) - 1);
            if (needMap.get(element) === 0) {
                needType--;
            }
        }
        // 已经找到了一个
        // 移动左指针
        while (needType === 0) {
            // console.log(s.slice(leftPoint, rightPoint + 1))
            const matchedString = s.slice(leftPoint, rightPoint + 1);
            // 记录匹配
            matchedStringMemo.push({ data: matchedString, length: matchedString.length });
            const element = s[leftPoint];
            if (needMap.has(element)) {
                needMap.set(element, needMap.get(element) + 1);
                if (needMap.get(element) === 1) {
                    needType++;
                }
            }
            leftPoint++;
        }
    }
    // 注意sort是原地排序
    return matchedStringMemo.sort((str1, str2) => str1.length - str2.length)[0].data;
};
console.log(minWindow('ADOBECODEBANC', 'ABC'));
