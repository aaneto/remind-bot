echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker build -t remind-bot .
docker tag remind-bot adilsinho/remind-bot:latest
docker tag remind-bot adilsinho/remind-bot:$TRAVIS_TAG
docker push adilsinho/remind-bot:latest
docker push adilsinho/remind-bot:$TRAVIS_TAG