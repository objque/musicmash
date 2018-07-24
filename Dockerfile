FROM scratch

COPY bin/musicmash musicmash

CMD ["/musicmash"]