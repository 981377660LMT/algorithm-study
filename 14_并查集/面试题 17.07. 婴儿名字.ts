const useUnionFindMap = (size: number) => {
  const parent = new Map<string, string>()

  const add = (key: string) => {
    if (parent.has(key)) return
    parent.set(key, key)
  }

  const find = (key: string) => {
    while (parent.get(key) !== key) {
      key = parent.get(key)!
    }

    return key
  }

  // 在结果列表中，选择 字典序最小 的名字作为真实名字
  const union = (key1: string, key2: string) => {
    const root1 = find(key1)
    const root2 = find(key2)
    const parentRoot = root1 < root2 ? root1 : root2
    const childRoot = root1 < root2 ? root2 : root1
    parent.set(childRoot, parentRoot)
  }

  return { union, find, add }
}

/**
 * @param {string[]} names
 * @param {string[]} synonyms
 * @return {string[]}
 * 设计一个算法打印出每个真实名字的实际频率
 * 在结果列表中，选择 字典序最小 的名字作为真实名字
 */
const trulyMostPopular = function (names: string[], synonyms: string[]): string[] {
  const uf = useUnionFindMap(names.length)

  const detail = names
    .map(info => info.split(/[\(\)]/))
    .map(item => ({
      name: item[0],
      count: item[1],
    }))

  synonyms.forEach(pair => {
    const tmp = pair.replace(/[\(\)]/g, '').split(',')
    const [key1, key2] = tmp

    uf.add(key1)
    uf.add(key2)
    uf.union(key1, key2)
  })

  const res = new Map<string, number>()

  detail.forEach(item => {
    // 找不到root则说明没有另外名字
    const root = uf.find(item.name) || item.name
    res.set(root, (res.get(root) || 0) + +item.count)
  })

  return [...res.entries()].map(entry => `${entry[0]}(${entry[1]})`)
}

console.log(
  trulyMostPopular(
    [
      'Fcclu(70)',
      'Ommjh(63)',
      'Dnsay(60)',
      'Qbmk(45)',
      'Unsb(26)',
      'Gauuk(75)',
      'Wzyyim(34)',
      'Bnea(55)',
      'Kri(71)',
      'Qnaakk(76)',
      'Gnplfi(68)',
      'Hfp(97)',
      'Qoi(70)',
      'Ijveol(46)',
      'Iidh(64)',
      'Qiy(26)',
      'Mcnef(59)',
      'Hvueqc(91)',
      'Obcbxb(54)',
      'Dhe(79)',
      'Jfq(26)',
      'Uwjsu(41)',
      'Wfmspz(39)',
      'Ebov(96)',
      'Ofl(72)',
      'Uvkdpn(71)',
      'Avcp(41)',
      'Msyr(9)',
      'Pgfpma(95)',
      'Vbp(89)',
      'Koaak(53)',
      'Qyqifg(85)',
      'Dwayf(97)',
      'Oltadg(95)',
      'Mwwvj(70)',
      'Uxf(74)',
      'Qvjp(6)',
      'Grqrg(81)',
      'Naf(3)',
      'Xjjol(62)',
      'Ibink(32)',
      'Qxabri(41)',
      'Ucqh(51)',
      'Mtz(72)',
      'Aeax(82)',
      'Kxutz(5)',
      'Qweye(15)',
      'Ard(82)',
      'Chycnm(4)',
      'Hcvcgc(97)',
      'Knpuq(61)',
      'Yeekgc(11)',
      'Ntfr(70)',
      'Lucf(62)',
      'Uhsg(23)',
      'Csh(39)',
      'Txixz(87)',
      'Kgabb(80)',
      'Weusps(79)',
      'Nuq(61)',
      'Drzsnw(87)',
      'Xxmsn(98)',
      'Onnev(77)',
      'Owh(64)',
      'Fpaf(46)',
      'Hvia(6)',
      'Kufa(95)',
      'Chhmx(66)',
      'Avmzs(39)',
      'Okwuq(96)',
      'Hrschk(30)',
      'Ffwni(67)',
      'Wpagta(25)',
      'Npilye(14)',
      'Axwtno(57)',
      'Qxkjt(31)',
      'Dwifi(51)',
      'Kasgmw(95)',
      'Vgxj(11)',
      'Nsgbth(26)',
      'Nzaz(51)',
      'Owk(87)',
      'Yjc(94)',
      'Hljt(21)',
      'Jvqg(47)',
      'Alrksy(69)',
      'Tlv(95)',
      'Acohsf(86)',
      'Qejo(60)',
      'Gbclj(20)',
      'Nekuam(17)',
      'Meutux(64)',
      'Tuvzkd(85)',
      'Fvkhz(98)',
      'Rngl(12)',
      'Gbkq(77)',
      'Uzgx(65)',
      'Ghc(15)',
      'Qsc(48)',
      'Siv(47)',
    ],
    [
      '(Gnplfi,Qxabri)',
      '(Uzgx,Siv)',
      '(Bnea,Lucf)',
      '(Qnaakk,Msyr)',
      '(Grqrg,Gbclj)',
      '(Uhsg,Qejo)',
      '(Csh,Wpagta)',
      '(Xjjol,Lucf)',
      '(Qoi,Obcbxb)',
      '(Npilye,Vgxj)',
      '(Aeax,Ghc)',
      '(Txixz,Ffwni)',
      '(Qweye,Qsc)',
      '(Kri,Tuvzkd)',
      '(Ommjh,Vbp)',
      '(Pgfpma,Xxmsn)',
      '(Uhsg,Csh)',
      '(Qvjp,Kxutz)',
      '(Qxkjt,Tlv)',
      '(Wfmspz,Owk)',
      '(Dwayf,Chycnm)',
      '(Iidh,Qvjp)',
      '(Dnsay,Rngl)',
      '(Qweye,Tlv)',
      '(Wzyyim,Kxutz)',
      '(Hvueqc,Qejo)',
      '(Tlv,Ghc)',
      '(Hvia,Fvkhz)',
      '(Msyr,Owk)',
      '(Hrschk,Hljt)',
      '(Owh,Gbclj)',
      '(Dwifi,Uzgx)',
      '(Iidh,Fpaf)',
      '(Iidh,Meutux)',
      '(Txixz,Ghc)',
      '(Gbclj,Qsc)',
      '(Kgabb,Tuvzkd)',
      '(Uwjsu,Grqrg)',
      '(Vbp,Dwayf)',
      '(Xxmsn,Chhmx)',
      '(Uxf,Uzgx)',
    ]
  )
)

export {}
// ["John(27)","Chris(36)"]
