# TankBattleServer
A tankbattle online game

# Client
[YoRHa-Hacking](https://github.com/o0olele/YoRHa-Hacking)

# Prepare
## plugins
  - `github.com/golang/protobuf`
  - `github.com/golang/tools`
  - `github.com/garyburd/redigo`
  - `github.com/golang/glog`
  - `golang.org/x`
  - `google.golang.org/genproto`
  - `google.golang.org/grpc`
  - `google.golang.org/protobuf`

get these plugins and move them into `./TankBattleBase/src`

## necessities
  - nginx
  - redis
  - [recastnavigation-go](https://github.com/fananchong/recastnavigation-go)

# Start

~~~shell
cd TankBattle
mkdir log
make
./run.sh
~~~

