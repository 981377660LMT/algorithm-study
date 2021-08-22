```JS
// Setup
var str = Array(10000).fill('a').join('');
// while
var str2ar1 = function(s) {
    var r = Array(s.length), i = -1
    while(++i < s.length) r[i] = s[i]
};
// Array.from
var str2ar2 = function(s) {
    return Array.from(s)
};
// split('')
var str2ar3= function(s) {
    return s.split('')
};
// 扩展运算符
var str2ar4= function(s) {
    return [...s]
};
// 正则
var str2ar5= function(s) {
    return s.split(/(?=[\s\S])/u)
};

```

split('')最快
Splitting a string into an array is about 70 times faster with 'a string'.split('') than Array.from('a string')
**Never Use Array.from() to Convert Strings to Arrays**
