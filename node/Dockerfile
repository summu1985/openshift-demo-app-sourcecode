FROM image-registry.openshift-image-registry.svc:5000/openshift/nodejs:18-minimal-ubi9
USER root
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN npm install express ws
RUN chown -R 1001:0 /app
EXPOSE 8080
USER 1001
CMD ["node", "web-socket-server.js"]
