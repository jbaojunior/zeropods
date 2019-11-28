FROM alpine:3.10

# Image To build a application to do scale in a k8s deployment 
LABEL maintainer="jbaojunior@gmail.com"

RUN apk update && apk add bash && rm -rf /var/cache/apk/*

ADD camunda-autoscaler /usr/local/bin/
RUN chmod +x /usr/local/bin/camunda-autoscaler

ENTRYPOINT ["/usr/local/bin/camunda-autoscaler"]