FROM vail130/orinoco-base

ADD scripts/run-tests.sh /usr/local/bin/run-tests
RUN chmod +x /usr/local/bin/run-tests

CMD ["--", "/usr/local/bin/run-tests", "s3"]
