
VAddy for Private Net command tool
======================================

VAddy for Private Net コマンド  
[https://vaddy.net](https://vaddy.net)

イントラネット内のサーバや、ローカルPC、VMなどの環境にもVAddyの脆弱性検査が実施できるツールです。


![screen](../images/screen.png "screen")


## 動作環境

このツールでは、javaとssh, ssh-keygen, psコマンドを利用します。  
そのため、現在ではMacとLinuxのみサポートとしています。


## コマンドの終了ステータス

エラーや脆弱性が発見されなかった場合は終了コード 0を返します。これは一般的なコマンドの正常終了と同じ終了コードです。
エラーや脆弱性があった場合は、終了コード1を返します。

## 設定方法

VAddyのWeb画面からWebAPIキーを発行してください。  
[https://console.vaddy.net/user/webapi](https://console.vaddy.net/user/webapi)  

次に、`conf/vaddy.conf.example` ファイルを`conf/vaddyconf`にコピーして、設定情報を書き込みます。


## 動作説明

このツールは、sshのリモートポートフォワードを使ったsshトンネルを作ります。
ローカルのWebサーバのポートが、sshトンネルによりVAddyのサーバ側に公開されますので、VAddyはそこを通して検査します。  
ローカルのWebサーバのポートが外部に公開される形ですが、その公開されたポートへはVAddyサーバのみアクセスできるように制限されていますのでご安心ください。



## 使い方 

### 引数

    Usage: ./vaddy_privatenet.sh action [-crawl crawl_id or crawl_label]   


| action        |                                                           | 
| ------------- |:---------------------------------------------------------:| 
| connect       | VAddyサーバとsshトンネルを張ります                                       |
| disconnect    | sshトンネルを切断します                                                 |
| scan          | VAddyサーバとsshトンネルを張り、検査を実行します。検査後はトンネルを切断します。 |

connectは、sshトンネルをバックグラウンドで作成した後に、VAddyサーバ経由でのコネクションチェックなど経て、問題なければ正常終了してプロンプトに戻ります。  
例えば、psコマンドでsshプロセスを見れば動作しているか確認できます。  
`ps aux | grep 'ssh -i vaddy/ssh/'`



#### 例1 (connect)

    ./vaddy_privatenet.sh connect

#### 例2 (検査開始)
scanアクションでは、検査実行前に自動的にsshトンネルを作成します。検査後は自動的にトンネルを切断します。

     ./vaddy_privatenet.sh scan


#### 例3 (クロールID 1234を指定した検査)
     ./vaddy_privatenet.sh scan -crawl 1234


オプション`crawl`は必須ではありません。これを指定しない場合は最新のクロールIDのデータを使って脆弱性検査します。


