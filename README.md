# grpc-gateway-sample

[golang-grpc-sample](https://github.com/pepese/golang-grpc-sample) の gRPC に対して REST からアクセス可能とする。

## 環境構築

```
$ GO111MODULE=on
$ go mod init
$ go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
$ go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
$ go get -u -v github.com/golang/protobuf/protoc-gen-go
```

## 作成手順

1. 普通に `.proto` を作成する
2. proto に `google.api.http` を追加する（ `rpc` へ `option` を追加）
    - REST の HTTP Method + URI へのマッピングになる
      ```
      option (google.api.http) = {
        get: "/v1/helloworld/sayhello/{name}"
      };
    ```
3. `.proto` ファイルをビルドして、 `.pb.go` ファイルを作成する
    ```
    $ mkdir -p proto/dest/helloworld/v1
    $ mkdir -p proto/dest/helloworld/v2
    ```
    ```
    $ protoc proto/helloworldV1.proto \
        -I ./proto \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis \
        --go_out=plugins=grpc:./proto/dest/helloworld/v1
    ```
    ```
    $ protoc proto/helloworldV2.proto \
        -I ./proto \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis \
        --go_out=plugins=grpc:./proto/dest/helloworld/v2
    ```
    - [公式](https://github.com/grpc-ecosystem/grpc-gateway) などとは googleapis に対するパスが異なるが、これは GO MODULES を利用しているため、 `src` 配下ではなく `pkg` 配下にモジュールが配置されるため
        - パス内でバージョンを指定していることにも注意したい
4. `.proto` ファイルをビルドして、 `.pb.gw.go` ファイルを作成する
    ```
    $ protoc proto/helloworldV1.proto \
        -I ./proto \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis \
        --grpc-gateway_out=logtostderr=true:./proto/dest/helloworld/v1
    ```
    ```
    $ protoc proto/helloworldV2.proto \
        -I ./proto \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis \
        --grpc-gateway_out=logtostderr=true:./proto/dest/helloworld/v2
    ```
5. `.proto` ファイルをビルドして、 Swagger ファイルを作成する
    ```
    $ mkdir doc
    ```
    ```
    $ protoc proto/helloworldV1.proto \
        -I ./proto \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis \
        --swagger_out=logtostderr=true:./doc
    ```
    ```
    $ protoc proto/helloworldV2.proto \
        -I ./proto \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis \
        --swagger_out=logtostderr=true:./doc
    ```

## 動確

```
$ go run main.go
```

[golang-grpc-sample](https://github.com/pepese/golang-grpc-sample) を起動した上で以下を実行することで確認できる。

```
$ curl localhost:8081/v1/helloworld/sayhello/hoge
{"message":"Hello hoge v1"}

$ curl localhost:8081/v2/helloworld/sayhello/hoge
{"message":"Hello hoge v2"}
```