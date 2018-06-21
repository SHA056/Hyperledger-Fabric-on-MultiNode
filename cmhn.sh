#!/bin/bash

# Ask user for confirmation to proceed
function askProceed () {
  read -p "Continue (y/n)? " ans
  case "$ans" in
    y|Y )
      echo "proceeding ..."
    ;;
    n|N )
      echo "exiting..."
      exit 1
    ;;
    * )
      echo "invalid response"
      askProceed
    ;;
  esac
}


function clearAll (){
  echo
  echo "##########################################################"
  echo "######### Clearing docker network Check 01 ###############"
  echo "##########################################################"
  docker network prune
  echo
  echo "##########################################################"
  echo "########### Leaving Docker Swarm as Manager Node##########"
  echo "##########################################################"
  docker swarm leave -f
  echo
  echo "##########################################################"
  echo "######### Clearing docker network Check 02 ###############"
  echo "##########################################################"
  docker network prune
  echo
  echo "##########################################################"
  echo "############### Clearing docker Volume ###################"
  echo "##########################################################"
  docker volume prune
  echo
  echo "##########################################################"
  echo "######### Clearing docker network Check 02 ###############"
  echo "##########################################################"
  docker system prune
  echo
  echo "##########################################################"
  echo "######## Clearing running docker Containers ##############"
  echo "##########################################################"
  docker rm -f $(docker ps -q)
  echo
  echo "##########################################################"
  echo "######### Clearing all crypto-materials ##################"
  echo "##########################################################"
  sudo rm -rf crypto-config
  echo
  echo "##########################################################"
  echo "########### Clearing all channel artifacts ###############"
  echo "##########################################################"
  sudo rm -f channel-artifacts/*

}


# ask for confirmation to proceed
askProceed

# generate crypto-material
clearAll