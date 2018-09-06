FROM alpine

COPY siserver /opt/
COPY siclient /opt/

CMD /opt/siserver
