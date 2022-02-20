package main

import (
	"auto-NewLine/data"
	"auto-NewLine/env"
	"auto-NewLine/newline"
	"errors"
	"log"
	"math"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/urfave/cli/v2"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	exep, err := env.GetExecDir()
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
			&cli.StringFlag{
				Name:    "encode",
				Aliases: []string{"e"},
				Value:   "shift-jis",
				Usage:   "エンコード",
			},
			&cli.PathFlag{
				Name:    "setting",
				Aliases: []string{"s"},
				Value:   path.Join(exep, "setting.yml"),
				Usage:   "設定ファイルのパス",
			},
		},
		Action: appfunc,
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func appfunc(c *cli.Context) error {
	minlen := c.Int("minlen")
	maxlen := c.Int("maxlen")
	textp := c.Path("text")

	enc, err := GetEncoding(c.String("encode"))
	if err != nil {
		return err
	}

	text, err := LoadText(textp, enc)
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

	nlinfos, err := newline.GetInfo(string(text), stg.BreakPatterns, "user-dict.txt")
	if err != nil {
		return err
	}

	return WriteTextFile(
		textp,
		newline.Break(nlinfos, *stg, minlen, float32(float64(textlen)/math.Ceil(float64(textlen)/float64(maxlen))), maxlen),
		enc)
}

func GetEncoding(enc string) (encoding.Encoding, error) {
	e, ok := map[string]encoding.Encoding{
		"utf-8":     nil,
		"shift-jis": japanese.ShiftJIS,
	}[enc]

	if !ok {
		return nil, errors.New("the encoding was not found")
	}
	return e, nil
}

func LoadText(textp string, enc encoding.Encoding) ([]byte, error) {
	text, err := os.ReadFile(textp)
	if err != nil {
		return nil, err
	}

	if enc != nil {
		text, _, err = transform.Bytes(enc.NewDecoder(), text)
		if err != nil {
			return nil, err
		}
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

func WriteTextFile(path string, text string, enc encoding.Encoding) error {
	bytes := []byte(text)

	if enc != nil {
		b, _, err := transform.Bytes(enc.NewEncoder(), []byte(text))
		bytes = b
		if err != nil {
			return err
		}
	}

	return os.WriteFile(path, bytes, 0777)
}
