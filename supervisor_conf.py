#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os.path as pp
import sys

conf_model = '''
[unix_http_server]
file=/tmp/supervisor.sock   ; (the path to the socket file)

[inet_http_server]         ; inet (TCP) server disabled by default
port=127.0.0.1:9001        ; (ip_address:port specifier, *:port for all iface)
username=nada              ; (default is no username (open server))
password=nada               ; (default is no password (open server))

[supervisord]
logfile=/tmp/supervisord.log ; (main log file;default $CWD/supervisord.log)
logfile_maxbytes=50MB        ; (max main logfile bytes b4 rotation;default 50MB)
logfile_backups=10           ; (num of main logfile rotation backups;default 10)
loglevel=info                ; (log level;default info; others: debug,warn,trace)
pidfile=/tmp/supervisord.pid ; (supervisord pidfile;default supervisord.pid)
nodaemon=false               ; (start in foreground if true;default false)
minfds=1024                  ; (min. avail startup file descriptors;default 1024)
minprocs=200                 ; (min. avail process descriptors;default 200)

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
serverurl=unix:///tmp/supervisor.sock ; use a unix:// URL  for a unix socket
'''

m = """[program:%s_server]
command=%s
autostart=true
autorestart=false
directory=%s
stderr_logfile=%s/log/%s_err.log
environment=GOPATH="%s"

"""

def gen_program_conf(base, name):
    command = "%s/bin/%s_main" % (base, name)
    return m % (name, command, base, base, name, base)

def gen_gs_program_conf(base, name, server_id):
    gscmd = "gs%d" % (server_id)
    command = "%s/bin/%s_main -s %d" % (base, name, server_id)
    return m % (gscmd, command, base, base, gscmd, base)

def main():
    base = pp.dirname(pp.realpath(sys.argv[0]))
    conf = conf_model
    conf += gen_program_conf(base, 'center')
    conf += gen_gs_program_conf(base, 'gs', 1)
    conf += gen_gs_program_conf(base, 'gs', 2)
    conf += gen_program_conf(base, 'gate')
    open("supervisord.conf", 'w').write(conf)
    print conf

if __name__ == "__main__":
    main()
