#!/bin/bash
### BEGIN INIT INFO
# Provides:          msiandmeta
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: GoLang App using Azure MSI and Metadata
# Description:       Runs a Go Application which is a web server that demonstrates usage of Managed Service Identities and in-VM Instance Metadata
### END INIT INFO

appUserName=__USER__
appPath=__APP_PATH__
appName=__APP_NAME__

processIDFilename=$appPath/$appName.pid
logFilename=$appPath/$appName.log

#
# Starts the simple GO REST service
# 
start() {
    # Needed by the GO App to access subscription and resource group, correctly
    export SUBSCRIPTION_ID="__SUBSCRIPTION_ID__"
    export RESOURCE_GROUP="__RESOURCE_GROUP__"
    
    # Check if the service runs by looking at it's Process ID and Log Files
    if [ -f $processIDFilename ] && [ "`ps | grep -w $(cat $processIDFilename)`" ]; then
        echo 'Service already running' >&2
        return 1
    fi
    echo 'Starting service...' >&2
    su -c "start-stop-daemon -SbmCv -x /usr/bin/nohup -p \"$processIDFilename\" -d \"$appPath\" -- \"./$appName\" > \"$logFilename\"" $appUserName
    echo 'Service started' >&2
}

#
# Stops the simple GO REST service
#
stop() {
    if [ ! -f $processIDFilename ] && [ ! "`ps | grep -w $(cat $processIDFilename)`" ]; then
        echo "Service not running" >&2
        return 1
    fi
    echo "Stopping Service..." >&2
    start-stop-daemon -K -p "$processIDFilename"
    rm -f "$processIDFilename"
    echo "Service stopped!" >&2
}

#
# Main script execution
#

case $1 in

    start)
      start
      ;;

    stop)
      stop
      ;;

    restart)
      stop
      start
      ;;

    \?)
      echo "Usage: $0 start|stop|restart"
esac

