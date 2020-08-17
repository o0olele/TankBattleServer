APPS = logicserver rcenterserver roomserver
PROTOS = server.proto

BASE = $(PWD)/../TankBattleBase
PROTOC = /usr/local/bin
LIAOCC= /home/liaocc/go

BUILDVER = '0.1'
BUILDDATE = 'date +%F'

all: install

install: 
	export GOPATH=$(PWD):$(BASE):$(LIAOCC)\
	&& for ser in $(APPS);\
	do \
		go install -x $$ser;\
		if [ "$$?" != "0" ]; then\
			exit 1;\
		fi;\
	done

proto: 
	
	for ser in $(PROTOS);\
	do \
		protoc --go_out=plugins=grpc:$(PWD)/src proto/$$ser;\
		if [ "$$?" != "0" ]; then\
			exit 1;\
		fi;\
	done

 .PHONY : proto