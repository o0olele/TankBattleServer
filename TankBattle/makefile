APPS = logicserver

BASE = $(PWD)

BUILDVER = '0.1'
BUILDDATE = 'date +%F'

all: install

install: 
	export GOPATH=$(PWD)\
	&& for ser in $(APPS);\
	do \
		go install -x $$ser;\
		if [ "$$?" != "0" ]; then\
			exit 1;\
		fi;\
	done


logicserver:
	export GOPATH=$(PWD)\
	go install -x logicserver