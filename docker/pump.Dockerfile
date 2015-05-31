FROM vail130/orinoco-base

ADD scripts/run-pump.sh /usr/local/bin/run-pump
RUN chmod +x /usr/local/bin/run-pump

CMD ["--", "/usr/local/bin/run-pump"]
