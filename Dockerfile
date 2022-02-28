FROM debian:bullseye-slim
COPY build/blog-httpd /
CMD ["/blog-httpd", "serve"]
