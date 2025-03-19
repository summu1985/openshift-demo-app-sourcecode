#--Stage 1
FROM golang:alpine3.18 as build-stage
RUN mkdir -p /config /app /bin
COPY . /bin/. 
ENV ARCH=${ARCH:-amd64}
WORKDIR /bin/
RUN GOOS=linux GOARCH=${ARCH} go build -o /app/app-${ARCH} main.go
#--Stage 2
FROM golang:alpine3.18 
RUN mkdir -p /app/config /app/bin
ENV ARCH=${ARCH:-amd64}
ENV CONFIG_PATH=/app/config
COPY --from=build-stage /app /app/bin/.
COPY --from=build-stage /bin/todo-api.properties /app/config/.
RUN mv /app/bin/app-${ARCH} /app/bin/go-app
WORKDIR /app
ENTRYPOINT [ "./bin/go-app" ]
