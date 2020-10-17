package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils"

	"github.com/northfun/house/common/utils/ihttp"
)

var (
	rawCookies = "cna=ZeB2FxVKNTICAXW6B4cFwCok; tracknick=%5Cu534Eshow; thw=hk; miid=933691292144237585; v=0; cookie2=13d3d76f5e89eaece247e0fb1cf1e937; t=ff9f8bdeab0e9e1afa70854b15d6091d; _tb_token_=93b77b4e67ee; xlly_s=1; _samesite_flag_=true; _m_h5_tk=3f249b5519cc219e3696be04867f36d2_1602521845147; _m_h5_tk_enc=87f843503e435020187460e27843f36b; sgcookie=E100%2B790j8PZZ%2B8XYLm%2FGcFF9WTH4gZGQsuNZuGYTyQZw2rx%2FqIyYtoVSCbX6tJsr%2FItjR6xv2jnOOifLkCthoaX6Q%3D%3D; unb=114463686; uc3=id2=UoCNemP8iwKs&nk2=2HK3AAcO&vt3=F8dCufHAt4%2F63%2F90q8M%3D&lg2=VFC%2FuZ9ayeYq2g%3D%3D; csg=963f3c41; lgc=%5Cu534Eshow; cookie17=UoCNemP8iwKs; dnk=%5Cu534Eshow; skt=14730f7ab57e0b98; existShop=MTYwMjUxMjI2OA%3D%3D; uc4=nk4=0%402kduifIDCu%2BvxesyKAXDNoA%3D&id4=0%40UOgxLgL0RTf9vKzgc9mwiEuOBxQ%3D; _cc_=UtASsssmfA%3D%3D; _l_g_=Ug%3D%3D; sg=w67; _nk_=%5Cu534Eshow; cookie1=W51Zt2b2hkBQwBrX0Hhfxr3K9TibKSUpCcVeGiiJkLw%3D; isg=BDc32iN9JIqwjqF2IqT4GIOaxi2B_AteHwI9U4nkU4ZtOFd6kcybrvUaHpBmy-PW; mt=ci=10_1; uc1=existShop=false&pas=0&cookie21=VFC%2FuZ9ainBZ&cookie15=U%2BGCWk%2F75gdr5Q%3D%3D&cookie16=VFC%2FuZ9az08KUQ56dCrZDlbNdA%3D%3D&cookie14=Uoe0b0C%2FnuELBg%3D%3D; l=eBSgoFWqQdO0uZ5SBOfanurza77OSIRYYuPzaNbMiOCP_7fB52COWZ5JcjY6C3GNh6cJR35uFQ62BeYBcQAonxvOgqsSHxHmn; tfstk=ceZ1B39u2SEePBQ2bR6EQd6IweiAwaOIKNGUCPLBhpRoLb1cr3l7YoaDOH3-d"
	// rawCookies2 = "cna=ZeB2FxVKNTICAXW6B4cFwCok; tracknick=%5Cu534Eshow; thw=hk; miid=933691292144237585; v=0; cookie2=13d3d76f5e89eaece247e0fb1cf1e937; t=ff9f8bdeab0e9e1afa70854b15d6091d; _tb_token_=93b77b4e67ee; xlly_s=1; _samesite_flag_=true; _m_h5_tk=3f249b5519cc219e3696be04867f36d2_1602521845147; _m_h5_tk_enc=87f843503e435020187460e27843f36b; sgcookie=E100%2B790j8PZZ%2B8XYLm%2FGcFF9WTH4gZGQsuNZuGYTyQZw2rx%2FqIyYtoVSCbX6tJsr%2FItjR6xv2jnOOifLkCthoaX6Q%3D%3D; unb=114463686; uc3=id2=UoCNemP8iwKs&nk2=2HK3AAcO&vt3=F8dCufHAt4%2F63%2F90q8M%3D&lg2=VFC%2FuZ9ayeYq2g%3D%3D; csg=963f3c41; lgc=%5Cu534Eshow; cookie17=UoCNemP8iwKs; dnk=%5Cu534Eshow; skt=14730f7ab57e0b98; existShop=MTYwMjUxMjI2OA%3D%3D; uc4=nk4=0%402kduifIDCu%2BvxesyKAXDNoA%3D&id4=0%40UOgxLgL0RTf9vKzgc9mwiEuOBxQ%3D; _cc_=UtASsssmfA%3D%3D; _l_g_=Ug%3D%3D; sg=w67; _nk_=%5Cu534Eshow; cookie1=W51Zt2b2hkBQwBrX0Hhfxr3K9TibKSUpCcVeGiiJkLw%3D; mt=ci=10_1; uc1=cookie16=WqG3DMC9UpAPBHGz5QBErFxlCA%3D%3D&cookie21=VFC%2FuZ9ainBZ&cookie15=VT5L2FSpMGV7TQ%3D%3D&existShop=false&pas=0&cookie14=Uoe0b0C%2FmZFyjw%3D%3D; isg=BC4udf0Y_Z1H6wgJ09dxD8JZf4LwL_Ip_iEUjFj3mjHsO86VwL9COdQ596_X-OpB; l=eBSgoFWqQdO0uu3QBOfanurza77OSIRYYuPzaNbMiOCP_7CB5D51WZ5JrbT6C3GNh6KWR35uFQ62BeYBcQAonxvtOTIeFLMmn; tfstk=ccsAB0gU98HAvHpRLZUofif14iNhw3evsxOEXJ_1owxFzQ1maMjKAfaXpjK8e"
	rawCookies2 = "cna=ZeB2FxVKNTICAXW6B4cFwCok; tracknick=%5Cu534Eshow; thw=hk; miid=933691292144237585; v=0; cookie2=13d3d76f5e89eaece247e0fb1cf1e937; t=ff9f8bdeab0e9e1afa70854b15d6091d; _tb_token_=93b77b4e67ee; xlly_s=1; _samesite_flag_=true; _m_h5_tk=3f249b5519cc219e3696be04867f36d2_1602521845147; _m_h5_tk_enc=87f843503e435020187460e27843f36b; sgcookie=E100%2B790j8PZZ%2B8XYLm%2FGcFF9WTH4gZGQsuNZuGYTyQZw2rx%2FqIyYtoVSCbX6tJsr%2FItjR6xv2jnOOifLkCthoaX6Q%3D%3D; unb=114463686; uc3=id2=UoCNemP8iwKs&nk2=2HK3AAcO&vt3=F8dCufHAt4%2F63%2F90q8M%3D&lg2=VFC%2FuZ9ayeYq2g%3D%3D; csg=963f3c41; lgc=%5Cu534Eshow; cookie17=UoCNemP8iwKs; dnk=%5Cu534Eshow; skt=14730f7ab57e0b98; existShop=MTYwMjUxMjI2OA%3D%3D; uc4=nk4=0%402kduifIDCu%2BvxesyKAXDNoA%3D&id4=0%40UOgxLgL0RTf9vKzgc9mwiEuOBxQ%3D; _cc_=UtASsssmfA%3D%3D; _l_g_=Ug%3D%3D; sg=w67; _nk_=%5Cu534Eshow; cookie1=W51Zt2b2hkBQwBrX0Hhfxr3K9TibKSUpCcVeGiiJkLw%3D; mt=ci=10_1; uc1=cookie16=VFC%2FuZ9az08KUQ56dCrZDlbNdA%3D%3D&cookie21=VFC%2FuZ9ainBZ&cookie15=U%2BGCWk%2F75gdr5Q%3D%3D&existShop=false&pas=0&cookie14=Uoe0b0C%2FmtE5KA%3D%3D; isg=BG5utVtRvd0A58jJkxexT4IZv8IwbzJpvuFUzJg32nEsew7VAP-CeRR5M--XuCqB; l=eBSgoFWqQdO0u5pBBOfanurza77OSIRYYuPzaNbMiOCP_2fB5mSlWZ5JkCT6C3GNh6KWR35uFQ62BeYBcQAonxvtOTIeFLMmn; tfstk=cl_FBViUnUX_Z5T7DeTrAkPRKVfdw9XlhVRJxgWVJvn1nBfDbYdi5pRwMKUMx"

	// userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"

	urlSlc = []string{
		// "https://zc-paimai.taobao.com/zc_item_list.htm?spm=a219w.7474998.pagination.2.4f893c54tXtWwp&location_code=410105&auction_source=0&front_category=56950002&item_biz_type=6&sorder=2&st_param=-1&auction_start_seg=-1&page=1",
		// "https://zc-paimai.taobao.com/zc_item_list.htm?spm=a219w.7474998.pagination.2.4f893c54tXtWwp&location_code=410105&auction_source=0&front_cat%20%20egory=56950002&item_biz_type=6&sorder=2&st_param=-1&auction_start_seg=-1&page=1",
		// "https://zc-item.taobao.com/auction/623899018936.htm",
		// "https://sf-item.taobao.com/sf_item/623899018936.htm?spm=a219w.7474998.paiList.1.17303c54FDNxNQ",
		// "https://desc.alicdn.com/i7/600/390/602393294291/TB1.O5sfX67gK0jSZPf8quhhFla.desc%7Cvar%5Edesc%3Bsign%5E6308d2e3366668e9e12fb6e59fe019d5%3Blang%5Egbk%3Bt%5E1569291173",
		"https://desc.alicdn.com/i3/550/540/550547990743/TB14SagRXXXXXXgXFXX8qtpFXlX.desc%7Cvar%5Edesc%3Bsign%5E3cd0a8f57b8cc8265b8f57c6db2f551a%3Blang%5Egbk%3Bt%5E1560250762",
	}
)

func main() {
	mainQ, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	mainC := colly.NewCollector(
		colly.UserAgent(userAgent),
	)

	utils.GenCookies(rawCookies2)

	for i := range urlSlc {
		// mainQ.AddURL(urlSlc[i])
		// mainC.SetCookies(urlSlc[i], cc)

		data, err := ihttp.Get(urlSlc[i])
		if err != nil {
			fmt.Println("=======", err)
			return
		}

		retMap, err := ihttp.DealSubjectMatterTable(data)
		if err != nil {
			return
		}
		fmt.Println(retMap)

		var tb tbtype.TableSubjectMatterInfo
		utils.StructByReflect(retMap, &tb)
		fmt.Println(retMap)
		fmt.Println(tb)
	}

	mainC.OnRequest(func(r *colly.Request) {
		r.Headers.Set("sec-fetch-dest", "document")
		r.Headers.Set("sec-fetch-mode", "navigate")
		r.Headers.Set("sec-fetch-site", "none")
		r.Headers.Set("sec-fetch-user", "?1")
		r.Headers.Set("upgrade-insecure-requests", "1")
		r.Headers.Set("authority", "zc-paimai.taobao.com")
		r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	})

	// mainC.OnHTML("#J_NoticeDetail", func(e *colly.HTMLElement) {
	mainC.OnHTML("#J_desc", func(e *colly.HTMLElement) {
		// mainC.OnHTML("#J_ImgBooth", func(e *colly.HTMLElement) {
		fmt.Println("======d==url", e.Request.URL)
		fmt.Println("======d==body",
			strings.Replace(e.Text,
				" ", "", -1))
	})

	// mainC.OnHTML("#pro_details p", func(e *colly.HTMLElement) {
	// 	fmt.Println(e.Text)
	// })

	mainQ.Run(mainC)
}
