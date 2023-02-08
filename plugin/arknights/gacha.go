package arknights

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/FloatTech/imgfactory"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/fogleman/gg"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

func getProfessionImage(profession string) (*image.Image, error) {
	professionImagePath := fmt.Sprintf(arkNightDataPath+"static/profession_img/%s.png", profession)
	professionImage, err := gg.LoadImage(professionImagePath)
	if err != nil {
		return nil, err
	}
	return &professionImage, nil
}

func getRarityImage(rarity int8) (*image.Image, error) {
	rarityImagePath := fmt.Sprintf(arkNightDataPath+"static/gacha_rarity_img/%d.png", rarity)
	rarityImage, err := gg.LoadImage(rarityImagePath)
	if err != nil {
		return nil, err
	}
	return &rarityImage, nil
}

func getRarityBackImage(rarity int8, index int) (*image.Image, error) {
	rarityImage, err := getRarityImage(rarity)
	if err != nil {
		return nil, err
	}
	rarityBackRGBA := imageToRGBA(rarityImage)
	rarityBackImageRGBA := rarityBackRGBA.SubImage(image.Rect(27+index*123, 0, 149+index*123, 720))
	rarityBackImageCrop := gg.NewContextForImage(rarityBackImageRGBA)
	rarityBackImage := rarityBackImageCrop.Image()
	return &rarityBackImage, nil
}

func imageToRGBA(src *image.Image) *image.RGBA {
	if dst, ok := (*src).(*image.RGBA); ok {
		return dst
	}
	b := (*src).Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), *src, b.Min, draw.Src)
	return dst
}

func rollGacha(store storage) (charNames []string) {
	var rarity int
	for i := 0; i < 10; i++ {
		rarityRand := rand.Float64()
		switch {
		case rarityRand < 0.02:
			rarity = 5
		case rarityRand < 0.10:
			rarity = 4
		case rarityRand < 0.58:
			rarity = 3
		default:
			rarity = 2
		}
		if store.is6starsmode() {
			rarity = 5
		}
		indexRand := rand.Intn(len(rarity2CharName[rarity]))
		charNames = append(charNames, rarity2CharName[rarity][indexRand])
	}
	return
}
func gachaTextBuild(charNames []string) (gachaText string) {
	var gachaResult = make([][]string, 6)
	for _, charName := range charNames {
		charData := charTable[charName]
		gachaResult[charData.Rarity] = append(gachaResult[charData.Rarity], charData.Name)
	}
	for rarity := 5; rarity > 1; rarity-- {
		if len(gachaResult[rarity]) > 0 {
			switch {
			case rarity == 5:
				gachaText += "六星干员:"
			case rarity == 4:
				gachaText += "五星干员:"
			case rarity == 3:
				gachaText += "四星干员:"
			case rarity == 2:
				gachaText += "三星干员:"
			}
		}
		for index := 0; index < len(gachaResult[rarity]); index++ {
			if index != len(gachaResult[rarity])-1 {
				gachaText += gachaResult[rarity][index] + ", "
			} else {
				gachaText += gachaResult[rarity][index] + "\n"
			}
		}
	}
	return strings.TrimRight(gachaText, "\n")
}
func drawGachaImage(charNames []string) ([]byte, error) {
	backgroundImage, err := gg.LoadImage(arkNightDataPath + "static/gacha_background_img/2.png")
	background := gg.NewContextForImage(backgroundImage)
	if err != nil {
		return nil, err
	}
	for index, charName := range charNames {
		charData := charTable[charName]
		charImagePath := fmt.Sprintf(arkNightDataPath+"char_img/%s_1.png", charName)
		charImage, err := gg.LoadImage(charImagePath)
		if err != nil {
			return nil, err
		}
		rarityImage, err := getRarityImage(charData.Rarity)
		if err != nil {
			return nil, err
		}
		rarityBackImage, err := getRarityBackImage(charData.Rarity, index)
		if err != nil {
			return nil, err
		}
		professionImage, err := getProfessionImage(charData.Profession)
		if err != nil {
			return nil, err
		}

		background.DrawImage(*rarityBackImage, 0, 0)
		background.DrawImage(*rarityImage, 27+index*123, 0)
		background.DrawImage(charImage, 27+index*123, 175)
		background.DrawImage(*professionImage, 34+int(math.Round(float64(index)*122.5)), 490)
	}
	return imgfactory.ToBase64(background.Image())
}

func gacha(ctx *zero.Ctx) {
	if charTable == nil {
		readFile, err := os.ReadFile(arkNightDataPath + "character_table.json")
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
		}
		rarity2CharName = make([][]string, 6)
		json.Unmarshal(readFile, &charTable) //nolint:errcheck
		for charID, chardata := range charTable {
			if len(chardata.ItemObtainApproach) > 0 {
				rarity2CharName[chardata.Rarity] = append(rarity2CharName[chardata.Rarity], charID)
			}
		}
	}
	c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
	if !ok {
		ctx.SendChain(message.Text("找不到服务!"))
		return
	}
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	store := (storage)(c.GetData(gid))
	rollResult := rollGacha(store)
	i, err := drawGachaImage(rollResult)
	if err != nil {
		return
	}
	sendBase64 := "base64://" + helper.BytesToString(i)
	ctx.SendChain(
		message.Image(sendBase64),
		message.Text(gachaTextBuild(rollResult)),
		message.At(ctx.Event.UserID),
	)
}
