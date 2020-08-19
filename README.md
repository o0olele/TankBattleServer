# TankBattleServer
A tankbattle online game

# Client
[Simple](https://github.com/o0olele/TankBattleH5Client)

[UnityH5](https://github.com/suprecoder/TankBattleClient)

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

