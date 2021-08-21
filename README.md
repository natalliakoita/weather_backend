# weather_backend

## build and push image to docker hub
docker build -t natalliakoita/weather_backend:v1 .
docker push natalliakoita/weather_backend:v1

    where v1 - tag as version

## run project with docker compose
 1. prepare .env file. Copy example.env as .env and set valid values 
 2. docker-compose up