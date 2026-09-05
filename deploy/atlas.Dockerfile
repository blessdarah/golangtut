ARG ATLAS_VERSION=1.3.0

FROM arigaio/atlas:${ATLAS_VERSION} AS atlas
FROM golang:1.27-alpine

COPY --from=atlas /atlas /usr/local/bin/atlas

ENTRYPOINT ["atlas"]
