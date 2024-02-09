#!/bin/bash
set -xe

#SRP intallation
sudo rm -rf /srp-tools
sudo mkdir /srp-tools
sudo wget --tries=5 --timeout=60 --quiet --output-document /srp-tools/srp https://artifactory.eng.vmware.com/artifactory/srp-tools-generic-local/srpcli/0.13.30-20230919002607-45c259c-218/linux-amd64/srp
sudo chmod +x /srp-tools/srp
sudo /srp-tools/srp --version
sudo /srp-tools/srp update --yes
sudo /srp-tools/srp --version
sudo wget --tries=5 --timeout=60 --quiet --output-document /tmp/linux-observer-2.0.0.tar.gz https://artifactory.eng.vmware.com/artifactory/srp-tools-generic-local/observer/2.0.0/linux-observer-2.0.0.tar.gz
sudo mkdir /srp-tools/observer
cd /srp-tools/observer
sudo tar zxf /tmp/linux-observer-2.0.0.tar.gz
sudo /srp-tools/observer/bin/observer_agent --version

SRP_WORKING_DIR=$WORKSPACE/provenance
#start observer
[ -d "$SRP_WORKING_DIR" ] && sudo rm -rf "$SRP_WORKING_DIR"
mkdir -p $SRP_WORKING_DIR
sudo /srp-tools/observer/bin/observer_agent -m start_observer --output_environment ${SRP_WORKING_DIR}/envs.sh --env_to_shell
source ${SRP_WORKING_DIR}/envs.sh
