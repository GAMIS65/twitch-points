FROM postgres:16-alpine
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=giveaway
EXPOSE 5432
