package main

import (
	"auto-NewLine/data"
	"auto-NewLine/newline"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/urfave/cli/v2"
)

func main() {
	exedp, err := func() (string, error) {
		p, err := os.Executable()
		return filepath.Dir(p), err
	}()
	if err != nil {
		println(err)
		return
	}

	app := &cli.App{
		Name:            "auto-NewLine",
		Usage:           Version,
		Description:     "指定したテキストファイルを最小文字数と最大文字数に合わせていい感じに改行します",
		Version:         Version,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "text",
				Aliases:  []string{"t"},
				Usage:    "改行するテキストファイルのパス",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "minlen",
				Aliases: []string{"min"},
				Value:   10,
				Usage:   "最小文字数",
			},
			&cli.IntFlag{
				Name:    "maxlen",
				Aliases: []string{"max"},
				Value:   30,
				Usage:   "最大文字数",
			},
			&cli.Float64Flag{
				Name:    "aimpos",
				Aliases: []string{"aim"},
				Value:   -1.0,
				Usage:   "どの程度で改行を試みるかの割合。0.0 - 1.0 を設定します。範囲外の場合、テキストに合わせて自動的に設定されます",
				EnvVars: []string{},
			},
			&cli.PathFlag{
				Name:    "setting",
				Aliases: []string{"s"},
				Value:   path.Join(exedp, "setting.yml"),
				Usage:   "設定ファイルのパス",
			},
		},
		Action: appfunc(exedp),
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func appfunc(exedp string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		minlen := c.Int("minlen")
		maxlen := c.Int("maxlen")
		textp := c.Path("text")

		text, err := LoadText(textp)
		if err != nil {
			return err
		}

		// テキストが最大文字数以下なら何もしない
		textlen := len([]rune(string(text)))
		if textlen <= maxlen {
			return nil
		}

		stg, err := LoadSetting(c.Path("setting"))
		if err != nil {
			return err
		}

		nlinfos, err := newline.GetInfo(string(text), stg.BreakPatterns, path.Join(exedp, "user-dict.txt"))
		if err != nil {
			return err
		}

		aimpos := GetStringCountAimPos(
			float64(textlen),
			float64(minlen),
			float64(maxlen),
			c.Float64("aimpos"))

		breaked, err := newline.Break(nlinfos, *stg, minlen, aimpos, maxlen)
		if err != nil {
			return err
		}

		return WriteTextFile(
			textp,
			breaked)
	}
}

func GetStringCountAimPos(textlen, minlen, maxlen, aimpos float64) float64 {
	if aimpos < 0 || aimpos > 1 {
		return textlen / math.Ceil(textlen/maxlen)
	} else {
		return minlen + aimpos*(maxlen-minlen)
	}
}

func LoadText(textp string) ([]byte, error) {
	text, err := os.ReadFile(textp)
	if err != nil {
		return nil, err
	}

	return text, nil
}

func LoadSetting(path string) (*data.Setting, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := data.Setting{}
	err = yaml.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func WriteTextFile(path string, text string) error {
	bytes := []byte(text)
	return os.WriteFile(path, bytes, 0777)
}
