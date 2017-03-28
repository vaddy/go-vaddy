#!/usr/bin/env bash

echo '
################################################################
# VAddy Private Net Tools (Version: 0.0.1)
# This software is released under the MIT License,
#
# This tool needs Mac or Linux, Java, ssh command
#
################################################################
'


############## Config/Hook file path ######################
CONFIG_FILE_PATH="./conf/vaddy.conf"
BEFORE_SCAN_HOOK_PATH="./conf/before_scan_hook.sh"
AFTER_SCAN_HOOK_PATH="./conf/after_scan_hook.sh"
VADDY_SCAN_LOG_PATH="./vaddy/scan_result.txt"


#Directory path of go-vaddy command line tool.
GOVADDY_BIN_DIR="../bin/"
###########################################################





############### Init and OS check #########################
load_config () {
	if [ -e ${CONFIG_FILE_PATH} ]; then
		source ${CONFIG_FILE_PATH}
	else
		echo -e "Error: Config file not found in ${CONFIG_FILE_PATH}"
		exit 1
	fi

	export VADDY_TOKEN=${VADDY_AUTH_KEY}
	export VADDY_HOST=${VADDY_FQDN}
}

get_os_type_bit() {
	if [ `uname` = "Darwin" ]; then
		OS_type='macosx'
		OS_bit="64bit"
	elif [ `uname` = "Linux" ]; then
		OS_type='linux'
		if [ `uname -m` = "i686" ]; then
			OS_bit="32bit"
		else
			OS_bit="64bit"
		fi
	else
        	echo -e "Error: Sorry, this tool supports only Mac and Linux now. Please wait for next update." 1>&2
        	exit 1
	fi
}

load_config
get_os_type_bit

#####################################################################



#command args
ACTION=$1
CRAWL_OPT_KEY=$2
CRAWL_OPT_VAL=$3


VADDY_AGENT_OPTIONS="-u ${VADDY_USER} -h ${VADDY_FQDN} -w ${VADDY_YOUR_LOCAL_IP}:${VADDY_YOUR_LOCAL_PORT}"


usage_exit() {
        echo -e "Usage: $0 action [-crawl crawl_id or crawl_label] \n  Action list: connect, disconnect, scan, check\n" 1>&2
	echo -e "Example1(make connection): $0 connect" 1>&2
	echo -e "Example2(start scan with crawlID 1234): $0 scan -crawl 1234" 1>&2
        exit 1
}

get_govaddy_binary_name() {
	echo "vaddy-"${OS_type}-${OS_bit}
}

set_crawl_label() {
	if [ "${CRAWL_OPT_KEY}" = '-crawl' ]  && [ "${CRAWL_OPT_VAL}" != "" ]; then
		echo -e "Set crawl Label: ${CRAWL_OPT_VAL}"
		export VADDY_CRAWL=${CRAWL_OPT_VAL}
	fi
}

connect() {
	echo -e "\n=== Connect ===\n"
	java -cp ./bin/vaddy_agent.jar vaddy.pip.Vaddy ${VADDY_AGENT_OPTIONS} -a start
	CONNECT_EXIT=$?
	echo -e "Connect Status: $CONNECT_EXIT"

	if [ $CONNECT_EXIT -ne 0 ]; then
        	echo -e "Error: Can not make connection for private net." 1>&2
		exit 1;
	fi
}

disconnect() {
	echo -e "\n=== Disconnect ===\n"
	java -cp ./bin/vaddy_agent.jar vaddy.pip.Vaddy ${VADDY_AGENT_OPTIONS} -a stop
}

start_scan() {
	echo -e "\n=== Start scan ===\n"

	if [ -e ${BEFORE_SCAN_HOOK_PATH} ]; then
		source ${BEFORE_SCAN_HOOK_PATH}
	fi

	GOVADDY_CLI=`get_govaddy_binary_name`
	if [ `which tee` ]; then
		${GOVADDY_BIN_DIR}${GOVADDY_CLI} | tee ${VADDY_SCAN_LOG_PATH}

		# with tee command, it can not get go-vaddy command exit status.
		SCAN_RESULT_GREP=`grep "Scan Success" ${VADDY_SCAN_LOG_PATH}`
		if [ "${SCAN_RESULT_GREP}" != "" ]; then
			GOVADDY_EXIT=0
		else
			GOVADDY_EXIT=1
		fi
	else
		${GOVADDY_BIN_DIR}${GOVADDY_CLI}
		GOVADDY_EXIT=$?
	fi

	echo -e "GoVAddy Status: $GOVADDY_EXIT"


	if [ -e ${AFTER_SCAN_HOOK_PATH} ]; then
		source ${AFTER_SCAN_HOOK_PATH}
	fi
}

health_check() {
	echo -e "\n=== Health Check ===\n"
}



case $1 in
	connect)
		connect
		exit
		;;
	disconnect)
		disconnect
		exit
		;;
	scan)
		connect
		set_crawl_label
		start_scan
		disconnect
		
		echo -e "Exit: ${GOVADDY_EXIT}"
		exit $GOVADDY_EXIT
		;;
	*)
		usage_exit
		exit
		;;
esac
