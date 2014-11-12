# Example directory structure

```
/data/images
├── application1
│   ├── Dockerfile
│   ├── prepare.sh
│   └── provision
│       ├── deploy.sh
│       └── unpack
│           ├── etc
│           │   └── start.d
│           │       ├── 10-process
│           │       ├── 20-process
│           │       └── 30-process
│           └── opt
│               └── application
│                   └── server
│                       └── config.ini
├── application2
│   ├── Dockerfile
│   ├── prepare.sh
│   └── provision
│       ├── some-config.ini
│       ├── deploy.sh
│       ├── MobileWeb
│       │   └── build
│       │       ├── classes
│       │       ├── config
│       │       ├── controllers
│       │       ├── modules
│       │       ├── views
│       │       └── webroot
│       ├── unpack
│       │   ├── etc
│       │   │   ├── confd
│       │   │   │   ├── conf.d
│       │   │   │   │   ├── some-config.toml
│       │   │   │   │   └── caching-config-json.toml
│       │   │   │   ├── confd.toml
│       │   │   │   └── templates
│       │   │   │       ├── some-config-ini.tmpl
│       │   │   │       └── caching-config-json.tmpl
│       │   │   ├── cron.d
│       │   │   │   └── application
│       │   │   └── start.d
│       │   │       ├── 10-crond
│       │   │       ├── 40-caching-writer
│       │   │       └── 99-php-fpm
│       │   └── opt
│       │       └── php
│       │           ├── etc
│       │           │   └── php-fpm.conf
│       │           └── lib
│       │               ├── php.ini
│       │               └── prepend.php
│       └── version
├── application3
│   ├── Dockerfile
│   └── provision
│       ├── deploy.sh
│       └── unpack
│           └── etc
│               ├── redis.conf
│               └── start.d
│                   └── 99-redis
├── application4
│   ├── Dockerfile
│   └── provision
│       ├── deploy.sh
│       └── unpack
│           ├── bin
│           │   └── nginx-conftest
│           ├── etc
│           │   ├── confd
│           │   │   ├── conf.d
│           │   │   │   └── nginx-upstreams.toml
│           │   │   └── templates
│           │   │       └── nginx-upstreams.tmpl
│           │   └── start.d
│           │       ├── 10-confd
│           │       └── 99-nginx
│           └── opt
│               └── nginx
│                   ├── certs
│                   │   ├── site1.com.au.crt
│                   │   ├── site1.com.au.key
│                   │   ├── site2.com.au.crt
│                   │   ├── site2.com.au.key
│                   └── conf
│                       ├── mime.types
│                       ├── nginx.conf
│                       ├── sites-enabled
│                       │   ├── mobile
│                       │   └── www
│                       └── upstreams.conf
└── wharfie
```
