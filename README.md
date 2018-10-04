# orderbot

# Set up interactive components for localhost
```
ngrok http 3000
```
and set it on slack interavtive

# Compile for linux
```
dep ensure
env GOOS=linux GOARCH=386 go build -o main
```
# Start docker
```
docker-compose build
docker-compose up
```

# Feeling lazy (one liner)
```
env GOOS=linux GOARCH=386 go build -o main && docker-compose up --build
```

# Set up for production server
1. Install docker
```
# cd /etc/yum.repos.d/
# wget http://yum.oracle.com/public-yum-ol7.repo
# vi public-yum-ol7.repo

[ol7_latest]
name=Oracle Linux $releasever Latest ($basearch)
baseurl=http://yum.oracle.com/repo/OracleLinux/OL7/latest/$basearch/
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-oracle
gpgcheck=1
enabled=1

[ol7_UEKR4]
name=Latest Unbreakable Enterprise Kernel Release 4 for Oracle Linux $releasever ($basearch)
baseurl=http://yum.oracle.com/repo/OracleLinux/OL7/UEKR4/$basearch/
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-oracle
gpgcheck=1
enabled=1

[ol7_addons]
name=Oracle Linux $releasever Add ons ($basearch)
baseurl=http://yum.oracle.com/repo/OracleLinux/OL7/addons/$basearch/
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-oracle
gpgcheck=1
enabled=1
```
2. Set up git repo
```
git clone 
```
3. Set up cron