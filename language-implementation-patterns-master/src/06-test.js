var expect = require('expect')
var parse = require('./06')

var str = `
\`\`\`zk
        title : 有这些礼物搞定宅男并不难
        desc: 如果你喜欢的男生是个宅男/GEEK，总之就是闷得不行，让你一筹莫展？别着急，调调来帮你！趁着圣诞/新年之际，挑选一个戳中他的礼物表达你的心意，一定会有所进展！
        image : ![](//content.image.alimmdn.com/cms/sites/default/files/20151215/zk/4234cover.jpg)
        haha : 测试新字段
        hehe : 测试另一个新字段
\`\`\`
\`\`\`card
        id: 2419
        title: 游戏、电影的好伴侣：Sony MDR-HW700DS无线耳机
        desc: 投其所好最重要！可是游戏动漫都不懂，咋办？此时请默念，索尼大法好。它是和PS3/PS4等配合最好的多声道耳机，当然xbox系支持也不在话下，游戏迷神器，看电影也是效果拔群，是这个价位上效果最好的无线耳机。
        image: ![](//content.image.alimmdn.com/cms/sites/default/files/20150701/goodthing/c_old.jpg)
\`\`\``


describe('test parse a special markdown block', () => {
	it('should parse correctly', () => {
		let result = parse(str)
		expect(result).toEqual([{
			type: 'zk',
			data: {
				title: '有这些礼物搞定宅男并不难',
				desc: '如果你喜欢的男生是个宅男/GEEK，总之就是闷得不行，让你一筹莫展？别着急，调调来帮你！趁着圣诞/新年之际，挑选一个戳中他的礼物表达你的心意，一定会有所进展！',
				image: '![](//content.image.alimmdn.com/cms/sites/default/files/20151215/zk/4234cover.jpg)',
				haha: '测试新字段',
				hehe: '测试另一个新字段',
			}
		}, {
			type: 'card',
			data: {
				id: '2419',
				title: '游戏、电影的好伴侣：SonyMDR-HW700DS无线耳机',
				desc: '投其所好最重要！可是游戏动漫都不懂，咋办？此时请默念，索尼大法好。它是和PS3/PS4等配合最好的多声道耳机，当然xbox系支持也不在话下，游戏迷神器，看电影也是效果拔群，是这个价位上效果最好的无线耳机。',
				image: '![](//content.image.alimmdn.com/cms/sites/default/files/20150701/goodthing/c_old.jpg)'
			}
		}])
	})
})