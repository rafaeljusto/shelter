Installation
============

The tests were made using an Amazon EC2 ms1.small instance, with 24GB of disk. The disk
space is an issue here, because the integration tests accumulates MongoDB data really fast
(the best approach would be changing the tests to erase the database). For now the
integration tests job is erasing the databases via shell script.

## Locale

```
sudo update-locale LANG=en_US.utf8 LC_MESSAGES=POSIX
```

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

## Karma

```
sudo apt-get install npm
sudo npm install -g karma
sudo npm install -g karma-cli
sudo npm install -g karma-jasmine
sudo npm install -g karma-phantomjs-launcher
sudo npm install -g karma-ng-html2js-preprocessor
```

## FPM

```
sudo apt-get install rubygems
sudo gem install fpm
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
* Git plugin
* Github OAuth Plugin

## Import jobs into Jenkins (via Jenkins web interface)

All jobs are found in `<shelter>/deploy/jenkins/jobs`. You should follow the pipeline:

1. shelter-retrieve
2. shelter-build
3. shelter-unit-test
4. shelter-integration-test
5. shelter-interface-test
6. shelter-publish

You can import each job using the command:

```
java -jar jenkins-cli.jar -s http://server create-job newmyjob < myjob.xml
```

## Update API key

Update API key from job shelter-deploy. This key can be found in https://bintray.com
account, under user settings.
