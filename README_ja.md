
go-vaddy: VAddy API Command-Line Tool
=================================

VAddy API Command-Line Tool using golang  
https://vaddy.net

VAddyの脆弱性検査の実行と結果の取得を自動化するコマンドツールです。

## OSの種類

`go-vaddy/bin`ディレクトリに、OS毎の実行ファイルが置いてあります。
もしlinux(64bit)をお使いの場合は、vaddy-linux-64bitというファイルを実行してください。

例： `./vaddy-linux-64bit  api_key userID FQDN`

| OS            | file               |
| ------------- |:------------------:|
| Linux(64bit)  | vaddy-linux-64bit  |
| MacOS(64bit)  | vaddy-macosx-64bit |
| Windows(64bit)| vaddy-win-64bit.exe|
| FreeBSD(64bit)| vaddy-freebsd-64bit|
| Linux(32bit)  | vaddy-linux-32bit  |
| Windows(32bit)| vaddy-win-32bit.exe|
| FreeBSD(32bit)| vaddy-freebsd-32bit|



## 利用方法 (脆弱性検査の実行と結果の取得)

### 終了ステータス
Go-vaddyは、エラーや脆弱性が発見されなかった場合は終了コード 0を返します。これは一般的なコマンドの正常終了と同じ終了コードです。
エラーや脆弱性があった場合は、終了コード1を返します。


### 環境変数

検査対象のサーバをVAddyに登録した時期によってご利用のVAddyのプロジェクトのバージョン(V1/V2)が異なります。  
ご利用のプロジェクトのバージョンを確認する場合は、ログイン後のDashboard画面にてご確認ください。  

### VADDY_TOKENの取得方法
コマンドに設定する `VADDY_TOKEN` の情報は、下記のAPI設定ページから「Create WebAPI key」ボタンを押してAPIキーを発行してください。
発行した画面の「API Auth Key」の値が、`VADDY_TOKEN`にセットする値になります。  
https://console.vaddy.net/user/webapi

### 設定情報
#### V1プロジェクトの場合

    export VADDY_TOKEN="123455667789"  
    export VADDY_USER="ichikaway"  
    export VADDY_HOST="www.examplevaddy.com"  
    #export VADDY_CRAWL="30"  

#### V2プロジェクトの場合

    export VADDY_TOKEN="123455667789"
    export VADDY_USER="ichikaway"
    export VADDY_PROJECT_ID="your project id"
    #export VADDY_CRAWL="30"

`VADDY_USER`はログイン時のログインIDの値をセットしてください。  
`VADDY_CRAWL`はオプション項目で、指定しない場合は最新のクロールデータを使って検査します。  
下記の例のようにクロールラベルの指定も可能です。  
例えば、クロールラベルにuseredit1, useredit2のように付けていた場合は、検索でヒットした中の最新のクロールIDを指定して検査します。

    export VADDY_CRAWL="search result pages"  

### コマンド実行
環境変数をセットした後は、下記のようにコマンドを実行します。  ご利用の環境に合わせたコマンドを実行してください。今回はLinux環境の場合のコマンド例です。

    vaddy-linux-64bit


### その他の設定
#### Slack連携
この環境変数をセットすると脆弱性を発見した際にSlackにメッセージ通知できます。

    export SLACK_WEBHOOK_URL="webhook url"
    export SLACK_USERNAME="your user (optional)"
    export SLACK_CHANNEL="your channel (optional)"
    export SLACK_ICON_EMOJI=":smile: (optional)"
    export SLACK_ICON_URL="icon url (optional)"
