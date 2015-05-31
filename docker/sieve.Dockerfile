FROM vail130/orinoco-base

ADD scripts/run-sieve.sh /usr/local/bin/run-sieve
RUN chmod +x /usr/local/bin/run-sieve

CMD ["--", "/usr/local/bin/run-sieve"]
