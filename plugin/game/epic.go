package game

//抄的https://github.com/monsterxcn/nonebot_plugin_epicfree
import (
	"bytes"
	"errors"
	"fmt"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io"
	"net/http"
	"strings"
)

const (
	epicServiceName = "epic"
	epicUrl         = "https://www.epicgames.com/graphql"
)

func init() {
	engine := control.Register(epicServiceName, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- epic喜加1 xxx\n",
	})
	engine.OnRegex(`^(E|e)(P|p)(I|i)(C|c)?喜(加一|\+1)$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(
			func(ctx *zero.Ctx) {
				getEpicFree(ctx)
			},
		)
}

func getEpicFree(ctx *zero.Ctx) {
	gameInfo, err := getEpicGame()
	var msg []message.MessageSegment
	var gameDev, gamePub string
	var gameThumbnail string
	if err != nil {
		ctx.SendChain(message.Text(err))
		return
	}
	for _, game := range gjson.Get(gameInfo, "data.Catalog.searchStore.elements").Array() {
		gameName := game.Get("title").String()
		gameCorp := game.Get("seller.name").String()
		gameDev = gameCorp
		gamePub = gameCorp
		gamePrice := game.Get("price.totalPrice.fmtPrice.originalPrice").String()
		gamePromotions := game.Get("promotions.promotionalOffers").Array()
		upcomingPromotions := game.Get("promotions.upcomingPromotionalOffers").Array()
		if len(gamePromotions) == 0 || len(upcomingPromotions) > 0 {
			continue
		}
		for _, image := range game.Get("keyImages").Array() {
			if image.Get("type").String() == "Thumbnail" {
				gameThumbnail = image.Get("url").String()
			}
		}
		for _, pair := range game.Get("customAttributes").Array() {
			if pair.Get("key").String() == "developerName" {
				gameDev = pair.Get("value").String()
			}
			if pair.Get("key").String() == "publisherName" {
				gamePub = pair.Get("value").String()
			}
		}
		gameDesp := game.Get("description").String()
		endDate := game.Get("promotions.promotionalOffers.0.promotionalOffers.0.endDate").String()
		gameUrl := fmt.Sprintf("https://www.epicgames.com/store/zh-CN/p/%s", strings.Replace(game.Get("productSlug").String(), "/home", "", -1))
		if len(gameThumbnail) != 0 {
			msg = append(msg, message.Image(gameThumbnail))
		}
		msg = append(msg, message.Text(fmt.Sprintf("FREE now :: %s (%s)\n\n%s\n\n", gameName, gamePrice, gameDesp)))
		release := fmt.Sprintf("游戏由 %s 开发、%s 出版，", gameDev, gamePub)
		if gameDev == gamePub {
			release = fmt.Sprintf("游戏由 %s 发售，", gameDev)
		}
		msg = append(msg, message.Text(release))
		msg = append(msg, message.Text(fmt.Sprintf("将在 UTC 时间 %s 结束免费游玩，戳链接领取吧~\n%s\n", endDate, gameUrl)))
	}
	ctx.SendChain(msg...)
}

func getEpicGame() (gameInfo string, err error) {
	client := web.NewDefaultClient()
	body := `{
		"query":
		"query searchStoreQuery($allowCountries: String, $category: String, $count: Int, $country: String!, $keywords: String, $locale: String, $namespace: String, $sortBy: String, $sortDir: String, $start: Int, $tag: String, $withPrice: Boolean = false, $withPromotions: Boolean = false) {\n Catalog {\n searchStore(allowCountries: $allowCountries, category: $category, count: $count, country: $country, keywords: $keywords, locale: $locale, namespace: $namespace, sortBy: $sortBy, sortDir: $sortDir, start: $start, tag: $tag) {\n elements {\n title\n id\n namespace\n description\n effectiveDate\n keyImages {\n type\n url\n }\n seller {\n id\n name\n }\n productSlug\n urlSlug\n url\n items {\n id\n namespace\n }\n customAttributes {\n key\n value\n }\n categories {\n path\n }\n price(country: $country) @include(if: $withPrice) {\n totalPrice {\n discountPrice\n originalPrice\n voucherDiscount\n discount\n currencyCode\n currencyInfo {\n decimals\n }\n fmtPrice(locale: $locale) {\n originalPrice\n discountPrice\n intermediatePrice\n }\n }\n lineOffers {\n appliedRules {\n id\n endDate\n discountSetting {\n discountType\n }\n }\n }\n }\n promotions(category: $category) @include(if: $withPromotions) {\n promotionalOffers {\n promotionalOffers {\n startDate\n endDate\n discountSetting {\n discountType\n discountPercentage\n }\n }\n }\n upcomingPromotionalOffers {\n promotionalOffers {\n startDate\n endDate\n discountSetting {\n discountType\n discountPercentage\n }\n }\n }\n }\n }\n paging {\n count\n total\n }\n }\n }\n}\n",
			"variables": {
			"allowCountries": "CN",
				"category": "freegames",
				"count": 1000,
				"country": "CN",
				"locale": "zh-CN",
				"sortBy": "effectiveDate",
				"sortDir": "asc",
				"withPrice": true,
				"withPromotions": true
		}
	}`
	request, err := http.NewRequest("POST", epicUrl, bytes.NewBuffer([]byte(body)))
	if err == nil {
		// 增加header选项
		var response *http.Response
		request.Header.Add("Referer", "https://www.epicgames.com/store/zh-CN/")
		request.Header.Add("Content-Type", "application/json; charset=utf-8")
		response, err = client.Do(request)
		if err == nil {
			if response.StatusCode != http.StatusOK {
				return "", errors.New("epic 可能又抽风啦，请稍后再试")
			}
			data, _ := io.ReadAll(response.Body)
			response.Body.Close()
			gameInfo = binary.BytesToString(data)
			return
		}
	}
	return
}
