Installation
============

The tests were made using an Amazon EC2 ms1.small instance, with 24GB of disk. The disk
space is an issue here, because the MongoDB test enviroment accumulates data really fast,
so we should always erasing this temporary data (the best approach would change the tests
to erase the database).

1. Go

```
wget http://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz
tar -xzf go1.2.1.linux-amd64.tar.gz
sudo mv go /usr/local/
```

2. MongoDB

```
sudo apt-get install mongodb
```

3. Git

```
sudo apt-get install git
```

4. Jenkins

```
wget -q -O - http://pkg.jenkins-ci.org/debian/jenkins-ci.org.key | sudo apt-key add -
deb http://pkg.jenkins-ci.org/debian binary/
sudo apt-get update
sudo apt-get install jenkins
```

5. Jenkins plugins (via Jenkins web interface)

* Build Pipeline Plugin
* Environment Injector Plugin
* GIT plugin

6. Import jobs into Jenkins (via Jenkins web interface)