FROM vail130/orinoco-base

RUN apt-get update && apt-get install -y python

CMD ["--", "bash", "/go/src/github.com/vail130/orinoco/docker/scripts/run-tests.sh"]
