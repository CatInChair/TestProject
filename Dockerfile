FROM golang:bullseye
WORKDIR /project
RUN mkdir routes
RUN mkdir db
ADD go.sum go.mod ./
ADD main.go .
COPY db/* /project/db/
COPY routes/* /project/routes/
RUN CGO_ENABLED=1 GOOS=linux go build -o /project/test_project
EXPOSE 8080
CMD ["/project/test_project"]