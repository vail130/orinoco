FROM vail130/orinoco-base

ADD scripts/run-litmus.sh /usr/local/bin/run-litmus
RUN chmod +x /usr/local/bin/run-litmus

CMD ["--", "/usr/local/bin/run-litmus"]
