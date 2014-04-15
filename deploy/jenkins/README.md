Installation
============

The tests were made using an Amazon EC2 ms1.small instance, with 24GB of disk. The disk
space is an issue here, because the integration tests accumulates MongoDB data really fast
(the best approach would be changing the tests to erase the database). For now the
integration tests job is erasing the databases via shell script.

## Go

```
wget http://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz
tar -xzf go1.2.1.linux-amd64.tar.gz
sudo mv go /usr/local/
```

## MongoDB

```
sudo apt-get install mongodb
```

## Git

```
sudo apt-get install git
```

## Jenkins

```
wget -q -O - http://pkg.jenkins-ci.org/debian/jenkins-ci.org.key | sudo apt-key add -
deb http://pkg.jenkins-ci.org/debian binary/
sudo apt-get update
sudo apt-get install jenkins
```

## Jenkins plugins (via Jenkins web interface)

* Build Pipeline Plugin
* Environment Injector Plugin
* GIT plugin

## Import jobs into Jenkins (via Jenkins web interface)

All jobs are found in `<shelter>/deploy/jenkins/jobs`.