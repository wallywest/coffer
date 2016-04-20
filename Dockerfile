FROM golang:1.6.1-alpine

ADD docker-wait /bin/docker-wait
ENV REPO_PATH https://gitlab.vailsys.com/vail-cloud-services/coffer.git
ENV GITREPO master

RUN set -x \
  && apk add --no-cache --virtual .build-deps \
    gcc \
    make \
    musl-dev \
    git \
  && apk add --no-cache --update bash \
  && chmod 755 /bin/docker-wait \
  && git clone "$REPO_PATH" "$GOPATH/src/gitlab.vailsys.com/vail-cloud-services/coffer" \
  && cd "$GOPATH/src/gitlab.vailsys.com/vail-cloud-services/coffer" \
  && git checkout -q "$GITREPO" \
  && make tools \
  && make build && mv bin/coffer /bin/coffer \
  && apk del .build-deps

EXPOSE 6000
