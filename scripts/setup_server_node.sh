#!/bin/bash

#
# Function definitions for stuff that's being re-used
#
write_log_entry() {
    logText=$1
    echo "---- " `date` "----" | tee -a /var/log/setup_server_node.log
    echo "$logText" | tee -a /var/log/setup_server_node.log
}

help() {
    echo "This script installs GoLang and, copies the GO sample programs to the admins home directory and builds them."
    echo "Options:"
    echo "  -a      The name of the admin user to be able to copy the files to the home-directory of the admin!"
    echo "  -s      The ID of the subscription that should be set as a system-wide environment variable!"
    echo "  -r      The name of the resource group that should be set as a system-wide environment variable!"
}

#
# Main Script Execution
#
write_log_entry "Starting setup_server_node.sh..."
write_log_entry "Parsing options..."

while getopts ":a:s:r:" opt; do
    case $opt in
        a)
        adminName="$OPTARG"
        ;;

        s)
        subscriptionId="$OPTARG"
        ;;

        r)
        resGroup="$OPTARG"
        ;;

        \?)
        echo "Called with invalid options! Unknown option: -$OPTARG"
        write_log_entry "Called with invalid options! Unknown option: -$OPTARG"
        exit -10
    esac
done

#
# Verify the command line arguments
#
if [ ! $adminName ] || [ ! $subscriptionId ] || [ ! $resGroup ]; then
    write_log_entry "Missing parameters!"
    help
    exit -10
fi

#
# Next install all required bits
#

write_log_entry "Setting up server `hostname` for admin $adminName with subscription $subscriptionId and resource group $resGroup..."

write_log_entry "Installing GoLang 1.8.3..."
wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
tar xvf go1.8.3.linux-amd64.tar.gz
sudo chown -R root:root ./go
sudo mv go /usr/local

#
# Next compile the Go Application (assumes it's downloaded as assets with the ARM template custom script extension)
#
mkdir ./app
mv *.go ./app

export PATH="$PATH:/usr/local/go/bin"
export GOPATH="`realpath ./`/app"
export GOBIN="$GOPATH/bin"
go get ./app
go build -o msitests ./app

sudo mkdir /usr/local/msiandmeta
sudo cp ./msitests /usr/local/msiandmeta
sudo chown -R $adminName:$adminName /usr/local/msiandmeta

#
# Configure apache2 to use the Go application as a CGI script
#
cat ./template.msiandmeta.sh \
| awk -v USER="$adminName" '{gsub("__USER__", USER)}1' \
| awk -v APP_NAME="msitests" '{gsub("__APP_NAME__", APP_NAME)}1' \
| awk -v APP_PATH="/usr/local/msiandmeta" '{gsub("__APP_PATH__", APP_PATH)}1' \
| awk -v SUBS="$subscriptionId" '{gsub("__SUBSCRIPTION_ID__", SUBS)}1' \
| awk -v RGROUP="$resGroup" '{gsub("__RESOURCE_GROUP__", RGROUP)}1' \
>> msiandmeta.sh

#
# Now make sure the script is handled by the system for starting/stopping the service
#
sudo cp ./msiandmeta.sh /etc/init.d
sudo chmod +x /etc/init.d/msiandmeta.sh
sudo update-rc.d msiandmeta.sh defaults

#
# Finally, start the MSI Test GO Application as a service
#
sudo service msiandmeta start