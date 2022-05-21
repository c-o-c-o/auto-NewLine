# auto-NewLine
テキストファイルをいい感じに改行するソフトウェアです。

# 使い方
```
auto-NewLine.exe [options]
options
  -t テキストファイルのパス
  -min 最小文字数[10]
  -max 最大文字数[30]
  -aim どの程度で改行を試みるかの割合[-1]
       0.0 - 1.0 を設定します
       範囲外の場合、テキストに合わせて自動的に設定されます
  -e テキストファイルのエンコード[shift-jis]
  -s 設定ファイルのパス[setting.yml]
```

```
example
  auto-NewLine.exe -t テキストファイル.txt
```

#
 This software is released under the MIT License, see LICENSE.