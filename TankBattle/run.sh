#!/bin/sh

startwork()
{
    rm -rf $PWD/log/*

    nohup $PWD/bin/rcenterserver -config=$PWD/bin/config/config1.json &
    sleep 1s
    nohup $PWD/bin/logicserver -config=$PWD/bin/config/config1.json &
    nohup $PWD/bin/roomserver -config=$PWD/bin/config/config1.json &
}

stopwork()
{
    SERVERLIST='rcenterserver logicserver roomserver'

    for serv in $SERVERLIST
    do
        echo "stop $serv"
        ps aux|grep "/$serv"|sed -e "/grep/d"|awk '{print $2}'|xargs kill 2&>/dev/null
    done

    echo "running server:"`ps x|grep "server -c"|sed -e '/grep/d'|wc -l`
}

case $1 in
stop)
    stopwork
;;
start)
    stopwork
    sleep 1s
    startwork
;;
esac