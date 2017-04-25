FROM scratch
COPY ./{{.CommandLine.AppName | ToLower }} /
EXPOSE 443
ENTRYPOINT ["./{{.CommandLine.AppName | ToLower}}", "serve"]
