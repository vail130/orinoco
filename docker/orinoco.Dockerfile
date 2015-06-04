FROM vail130/orinoco-base

ADD scripts/run-orinoco.sh /usr/local/bin/run-orinoco
RUN chmod +x /usr/local/bin/run-orinoco

CMD ["--", "/usr/local/bin/run-orinoco"]
