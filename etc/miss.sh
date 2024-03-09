#!/bin/bash

export ctx_shy=${ctx_shy:=https://shylinux.com}
if [ -f $PWD/.ish/plug.sh ]; then source $PWD/.ish/plug.sh; elif [ -f $HOME/.ish/plug.sh ]; then source $HOME/.ish/plug.sh; else
	temp=$(mktemp); if curl -h &>/dev/null; then curl -o $temp -fsSL $ctx_shy; else wget -O $temp -q $ctx_shy; fi; source $temp intshell
fi; require conf.sh; require miss.sh; ish_sys_cli_prepare

ish_miss_prepare_compile
ish_miss_prepare_develop
ish_miss_prepare_project

ish_miss_prepare_contexts
ish_miss_prepare_intshell
ish_miss_prepare_learning
ish_miss_prepare_volcanos
ish_miss_prepare_toolkits
ish_miss_prepare_icebergs
ish_miss_prepare_release
ish_miss_prepare_modules
ish_miss_prepare icons

ish_miss_prepare go-sql-mysql

ish_miss_make; [ -z "$*" ] || ish_miss_serve "$@"
