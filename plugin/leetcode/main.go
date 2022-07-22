package leetcode

import (
	"fmt"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"regexp"
	"strings"
)

func init() {
	en := control.Register("leetcode", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "力扣每日一题 \n" + "- (leetcode|力扣)每日一题，每日算法",
	})
	en.OnFullMatchGroup([]string{"每日一题", "leetcode每日一题", "力扣每日一题", "每日算法"}).SetBlock(true).Limit(ctxext.LimitByUser).Handle(
		func(ctx *zero.Ctx) {
			getLeetcodeDaily(ctx)
		})
}

func getLeetcodeDaily(ctx *zero.Ctx) {
	client := resty.New()
	res, err := client.R().SetHeader("origin", "https://leetcode-cn.com").SetHeader("user-agent", web.RandUA()).SetBody(map[string]interface{}{
		"operationName": "questionOfToday",
		"variables":     "",
		"query":         "query questionOfToday { todayRecord {   question {     questionFrontendId     questionTitleSlug     __typename   }   lastSubmission {     id     __typename   }   date   userStatus   __typename }}",
	}).Post("https://leetcode-cn.com/graphql")
	if err != nil {
		return
	}
	titleJson := gjson.ParseBytes(res.Body())
	EnglishTitle := titleJson.Get("data.todayRecord.0.question.questionTitleSlug").String()
	QuestionUrl := "https://leetcode-cn.com/problems/" + EnglishTitle
	res, err = client.R().SetHeader("origin", "https://leetcode-cn.com").SetHeader("user-agent", web.RandUA()).SetBody(map[string]interface{}{
		"operationName": "questionData",
		"query":         "query questionData($titleSlug: String!) {  question(titleSlug: $titleSlug) {    questionId    questionFrontendId    boundTopicId    title    titleSlug    content    translatedTitle    translatedContent    isPaidOnly    difficulty    likes    dislikes    isLiked    similarQuestions    contributors {      username      profileUrl      avatarUrl      __typename    }    langToValidPlayground    topicTags {      name      slug      translatedName      __typename    }    companyTagStats    codeSnippets {      lang      langSlug      code      __typename    }    stats    hints    solution {      id      canSeeDetail      __typename    }    status    sampleTestCase    metaData    judgerAvailable    judgeType    mysqlSchemas    enableRunCode    envInfo    book {      id      bookName      pressName      source      shortDescription      fullDescription      bookImgUrl      pressImgUrl      productUrl      __typename    }    isSubscribed    isDailyQuestion    dailyRecordStatus    editorType    ugcQuestionId    style    __typename  }}",
		"variables": map[string]interface{}{
			"titleSlug": EnglishTitle,
		},
	}).Post("https://leetcode-cn.com/graphql")
	dailyData := gjson.ParseBytes(res.Body())
	Data := dailyData.Get("data.question")
	ID := Data.Get("questionFrontendId").String()
	Difficulty := Data.Get("difficulty").String()
	ChineseTitle := Data.Get("translatedTitle").String()
	Content := Data.Get("translatedContent").String()
	html, _ := htmlquery.Parse(strings.NewReader(Content))
	StringContent := htmlquery.InnerText(html)
	rg := regexp.MustCompile(`(\r\n?|\n){2,}`)
	msg := fmt.Sprintf(`[今日算法]：%s
[链接]：%s
[题目]：%s
[难度]：%s
[题目描述]：%s
`, ChineseTitle, QuestionUrl, ID, Difficulty, rg.ReplaceAllString(StringContent, "$1"))
	ctx.SendChain(message.Text(msg))
}
