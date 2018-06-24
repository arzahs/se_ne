FROM golang:latest
RUN mkdir -p /go/src/se_ne

WORKDIR /go/src/se_ne/
COPY . .
RUN curl https://glide.sh/get | sh && glide install
RUN go install -v se_ne
RUN cp -r /go/src/se_ne/assets /go/bin/assets && cp -r /go/src/se_ne/templates /go/bin/templates
RUN rm -rf ./*

ENV APP_HOST 0.0.0.0
ENV APP_PORT 8000
ENV APP_STATIC_PATH /go/bin/assets/
ENV APP_TEMPLATE_PATH /go/bin/templates/

EXPOSE 8000

ENTRYPOINT se_ne