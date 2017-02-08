FROM onsdigital/dp-go

WORKDIR /app/

COPY ./build/dp-dd-search-api .

ENTRYPOINT ./dp-dd-search-api
