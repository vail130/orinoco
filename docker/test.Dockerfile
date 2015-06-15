FROM vail130/orinoco-base

RUN apt-get update && apt-get install -y python

ADD scripts/run-tests.sh /usr/local/bin/run-tests
ADD scripts/reflect.py /usr/local/bin/reflect
RUN chmod +x /usr/local/bin/run-tests && \
	chmod +x /usr/local/bin/reflect

CMD ["--", "/usr/local/bin/run-tests"]
