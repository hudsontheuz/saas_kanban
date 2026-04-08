FROM migrate/migrate:v4.19.0

COPY . /migrations

ENTRYPOINT ["migrate", "-path", "/migrations"]