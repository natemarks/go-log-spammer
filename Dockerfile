FROM ubuntu:24.04
LABEL maintainer="Nate Marks <npmarks@gmail.com>"
COPY /build/linux/amd64/log-spammer /
RUN chmod +x /log-spammer
CMD ["/log-spammer"]

