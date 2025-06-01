#!/bin/sh

start(){
#    cd ../server
#    go run main.go
	cd ..
	./ninja_game
}

start_nohup(){
#    cd ../server
#    cat /dev/null > nohup.out
#    nohup go run main.go
	cd ..
    cat /dev/null > nohup.out
    nohup ./ninja_game &
}

start_nohup_out(){
#    cd ../server
#    cat /dev/null > nohup.out
#    nohup go run main.go
	cd ..
    nohup ./ninja_game > nohup.out &
}

proto(){
    cd ../proto3/
    ./export.sh
}

config(){
    cd ../../acinconfig/bin
    ./gengoconfig
}

# 如果编译的结果需要gdb调试则使用参数-gcflags “-N -l”,会关闭内联优化(可调试版本,性能低)
# 如果编译的结果需要发布.则使用-ldflags “-w -s”,可以去掉调试信息,减小大约一半的大小,关闭内联优化(不可调试版本,性能高)
# -s: 去掉符号信息 -w: 去掉DWARF调试信息
build(){
	cd ../server
	go build -ldflags '-w -s' main.go
	upx main		# 再加壳压缩
	mv main ../script/ninja_game
}

# mysql(){
# #	mysql -h127.0.0.1:3306 -uroot -phuang@feng@wu! -Dgame < ../sql/game.sql
# }

windows(){
	cd ../server
	GOOS=windows GOARCH=386 go build -ldflags '-w -s' -o main.exe main.go
}

# SIGINT=2
stop(){
	kill -2 `ps -ef|grep -v grep | grep ./ninja_game |awk '{print $2}'`
}

uptconfig(){
	kill -30 `ps -ef|grep -v grep | grep ./ninja_game |awk '{print $2}'`
}
upttable(){
	kill -31 `ps -ef|grep -v grep | grep ./ninja_game |awk '{print $2}'`
}

help()
{
	echo " manager  command:     "
	echo " start    以交互方式启动  "
	echo " start_nohup    以后台方式启动  "
	echo " start_nohup_out    以后台方式启动打印日志 "
	echo " stop     先踢人下线再关闭服务器  "
	echo " remote   远程连接服务器  "
	echo " proto    生成protobuff文件  "
	echo " config   生成配置文件  "
	echo " windows  生成window可执行包exe  "
	echo " build    生成linux包  "
	echo " uptconfig    更新配置表  "
	echo " upttable    更新策划表  "
}

case $1 in
	'start') start ;;
	'start_nohup') start_nohup ;;
	'start_nohup_out') start_nohup_out ;;
	'stop') stop ;;
	'remote') remote ;;
	'proto') proto ;;
	'config') config ;;
	'windows') windows ;;
	'build') build ;;
	'mysql') mysql ;;
	'uptconfig') uptconfig ;;
	'upttable') upttable ;;
	*) help ;;
esac