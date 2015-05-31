FROM vail130/orinoco-base

ADD scripts/run-tap.sh /usr/local/bin/run-tap
RUN chmod +x /usr/local/bin/run-tap

CMD ["--", "/usr/local/bin/run-tap"]
