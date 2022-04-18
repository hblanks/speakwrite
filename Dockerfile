FROM debian:bullseye-slim
COPY build/speakwrite /
CMD ["/speakwrite", "serve"]
