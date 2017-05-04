FROM scratch

COPY dist/mqtt-bq /mqtt-bq

ENV GOOGLE_APPLICATION_CREDENTIALS=/credentials/key.json

CMD ["/mqtt-bq"]