# wharfie

`wharfie` is a lightweight assistant for docker with a focus on removing large 'copy' layers from docker containers.

We found ourselves creating a directory structure as described in [example directory structure](docs/directory_structure.md)

The with the majority of our Dockerfile being

```
FROM registry.office/base:latest
COPY provision /tmp/provision
RUN /tmp/provision/deploy.sh
```

Which resulted in an unnecessary large layer in the middle of our image that was quite literally a duplicate of the data used elsewhere.

`wharfie` lets us avoid the COPY step by being a simple httpd with built in compression tied to our build structure (think, nginx but with rigid layout and built in packaging)

With the resulting Dockerfile becoming

```
FROM registry.office/base:latest
RUN /bin/echo -e "GET /application1 HTTP/1.0\r\nhost:172.17.42.1\r\n\r\n" | /bin/nc 172.17.42.1 2864 | /usr/bin/tail -n+6 | /bin/bash
```

This could equally be `curl http://127.17.42.1/application1 | /bin/bash` but the ubuntu base image doesn't have curl installed by default.

Everything under `/application1` will return a script that'll download and extract from `/application1/bundle` which is everything in `/application1/provisioning` compressed as a single tarball

##License

Copyright (c) 2014 Shannon Wynter, Ladbrokes Digital Australia Pty Ltd. Licensed under GPL2. See the [LICENSE.md](LICENSE.md) file for a copy of the license.