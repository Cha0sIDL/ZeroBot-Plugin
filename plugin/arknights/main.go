package arknights

import (
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/fogleman/gg"
)

func init() {
	Fonts, _ = gg.LoadFontFace(text.FontFile, 18)
	engine := control.Register("arknight", &control.Options{
		DisableOnDefault: false,
		//	PublicDataFolder: "ArkNights",
		Help: "查公招\n" + "方舟今日资源",
	})
	engine.OnRegex("^查公招$").Handle(recruit)
	engine.OnKeyword("方舟资源").Handle(daily)

}
