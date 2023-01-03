FROM public.ecr.aws/lambda/go:latest AS builder

# Copy lib (workaround to provide it as local package)
RUN mkdir /faas/ &&\
    mkdir /faas/lib
COPY lib/ /faas/lib

WORKDIR /faas/src
RUN yum install go -y
COPY . ./
RUN go mod download

# Build the binary
RUN chmod +x ./build_handlers.sh
RUN ./build_handlers.sh generated/

FROM public.ecr.aws/lambda/go:latest

COPY --from=builder /faas/src/bin /var/task

CMD [ "DeleteShop" ]