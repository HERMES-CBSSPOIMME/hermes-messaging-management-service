FROM scratch

WORKDIR /

ENV HERMES_CONFIG_FILE_PATH=/config.json

COPY ./management-service /

copy ./config.json /

EXPOSE 8085/tcp

CMD ["/management-service"]
