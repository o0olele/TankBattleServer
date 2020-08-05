
default:	build

clean:
	rm -rf Makefile objs

build:
	$(MAKE) -f objs/Makefile

install:
	$(MAKE) -f objs/Makefile install

modules:
	$(MAKE) -f objs/Makefile modules

upgrade:
	/data/nginx-1.10.1/sbin/nginx -t

	kill -USR2 `cat /data/nginx-1.10.1/logs/nginx.pid`
	sleep 1
	test -f /data/nginx-1.10.1/logs/nginx.pid.oldbin

	kill -QUIT `cat /data/nginx-1.10.1/logs/nginx.pid.oldbin`
