
go-vaddy: VAddy API command tool
=================================

VAddy API command tool using golang  
http://vaddy.net

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
| Linux(32bit)  | vaddy-linux-32bit  |
| Windows(32bit)| vaddy-win-32bit.exe|



## 利用方法 (脆弱性検査の実行と結果の取得)

### 終了ステータス
Go-vaddyは、エラーや脆弱性が発見されなかった場合は終了コード 0を返します。これは一般的なコマンドの正常終了と同じ終了コードです。
エラーや脆弱性があった場合は、終了コード1を返します。

### 引数タイプ
Go-vaddyでは、コマンドの引数を指定するパターンと、OSの環境変数にセットするパターンが選べます。


###コマンド引数
最後のオプション`crawl_id`は必須ではありません。これを指定しない場合は最新のクロールIDのデータを使って脆弱性検査します。

Usage: `vaddy-linux-64bit auth_key username(LoginID)  hostname crawl_id(optional)`

    ./vaddy-linux-64bit 123455667789  ichikaway  www.examplevaddy.net 30



###コマンド引数（クロールラベル指定）
Usage: `vaddy-linux-64bit auth_key username(LoginID)  hostname crawl_label_keyword`

    vaddy-linux-64bit 123455667789  ichikaway  www.examplevaddy.net useredit

あなたが付けたクロールラベルの文言を指定した例です。例えば、クロールラベルにuseredit1, useredit2のように付けていた場合は、検索でヒットした中の最新のクロールIDを指定して検査します。
検索で指定のラベルのものが見つからない場合は、`crawl_id`を指定しない検査となります（最新のcrawl idが利用されます）。


### ENV
コマンドの引数ではなく、OSの環境変数を利用した実行も可能です。


    export VADDY_TOKEN="123455667789"  
    export VADDY_USER="ichikaway"  
    export VADDY_HOST="www.examplevaddy.com"  
    export VADDY_CRAWL="30"  


`VADDY_CRAWL`はオプション項目で、指定しない場合は最新のクロールデータを使って検査します。  
下記の例のようにクロールラベルの指定も可能です。

    export VADDY_CRAWL="search result pages"  
