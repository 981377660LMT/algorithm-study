import DiffMatchPatch from 'diff-match-patch'

// 1. 创建实例
const dmp = new DiffMatchPatch()

{
  // ==================================================
  // 案例 1: Diff (差异比较)
  // ==================================================
  console.log('--- 案例 1: Diff ---')
  const text1 = 'The quick brown fox jumps over the lazy dog.'
  const text2 = 'That quick brown fox jumped over a lazy dog.'

  // 计算差异
  const diffs = dmp.diff_main(text1, text2)

  // 语义化清理，让结果更易读
  dmp.diff_cleanupSemantic(diffs)

  console.log('语义优化后 Diff 结果:', JSON.stringify(diffs))
  // 输出: [[-1,"The"],[1,"That"],[0," quick brown fox jump"],[-1,"s"],[1,"ed"],[0," over "],[-1,"the"],[1,"a"],[0," lazy dog."]]

  // 生成漂亮的 HTML
  const html = dmp.diff_prettyHtml(diffs)
  console.log('HTML 可视化结果:')
  console.log(html)
  // 输出: <del style="background:#ffe6e6;">The</del><ins style="background:#e6ffe6;">That</ins><span> quick brown fox jump</span><del style="background:#ffe6e6;">s</del><ins style="background:#e6ffe6;">ed</ins><span> over </span><del style="background:#ffe6e6;">the</del><ins style="background:#e6ffe6;">a</ins><span> lazy dog.</span>
}

{
  // ==================================================
  // 案例 2: Match (模糊匹配)
  // ==================================================
  console.log('\n--- 案例 2: Match ---')
  const longText = 'The quick brown fox jumps over the lazy dog.'
  const pattern = 'jumps'
  const expectedLoc = 20 // "jumps" 实际在 20

  // 精确匹配
  const matchLoc1 = dmp.match_main(longText, pattern, expectedLoc)
  console.log(`在位置 ${expectedLoc} 附近查找 "${pattern}": 找到于 ${matchLoc1}`) // 输出: 20

  // 模糊匹配
  const fuzzyPattern = 'jumped'
  const matchLoc2 = dmp.match_main(longText, fuzzyPattern, expectedLoc)
  console.log(`在位置 ${expectedLoc} 附近模糊查找 "${fuzzyPattern}": 找到于 ${matchLoc2}`) // 输出: 20 (因为 "jumps" 是最佳匹配)
}

{
  // ==================================================
  // 案例 3: Patch (制作与应用补丁)
  // ==================================================
  console.log('\n--- 案例 3: Patch ---')
  const originalText = 'The rain in Spain stays mainly in the plain.'
  const modifiedText = 'The rain in Spain falls mainly on the plain.'

  // 步骤 1: 基于两个文本制作补丁
  const patches = dmp.patch_make(originalText, modifiedText)

  // (可选) 步骤 2: 将补丁序列化为文本，便于存储或传输
  const patchText = dmp.patch_toText(patches)
  console.log('文本格式的补丁:\n', patchText)
  /*
输出类似:
@@ -15,15 +15,14 @@
 ain in Spain 
-stays
+falls
  mainly on
*/

  // 步骤 3: 将补丁应用到原始文本上
  const [newText, appliedStatus] = dmp.patch_apply(patches, originalText)

  console.log('补丁应用是否成功:', appliedStatus[0]) // true
  console.log('应用补丁后的新文本:', newText)
  console.log('新文本是否与目标文本一致:', newText === modifiedText) // true
}
