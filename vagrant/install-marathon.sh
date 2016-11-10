echo "Installing JAVA 8"
sudo apt-get -y update
sudo apt-get install -y software-properties-common
sudo add-apt-repository -y ppa:openjdk-r/ppa
sudo apt-get update
sudo apt-get install openjdk-8-jdk -y
sudo update-alternatives --set java /usr/lib/jvm/java-8-openjdk-amd64/jre/bin/java

#sudo apt-key adv --keyserver keyserver.ubuntu.com --recv E56151BF
#sudo add-apt-repository -y ppa:webupd8team/java
#sudo apt-get update -y

#echo debconf shared/accepted-oracle-license-v1-1 select true | sudo debconf-set-selections
#echo debconf shared/accepted-oracle-license-v1-1 seen true | sudo debconf-set-selections
#sudo apt-get install -y --force-yes oracle-java9-installer oracle-java9-set-default

# sudo apt-get install marathon chronos
# sudo service marathon start
# sudo service chronos start
