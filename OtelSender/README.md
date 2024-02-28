# RAIO Stub Service for Otel Experiment
This service stub is to do experiment on OpenTelemetry integration.

# Verify OpenTelemetry Exporter/Collector flow locally.
Start a local open telemetry collector. The */home/huliang/testdata/yamls/otel/config.debug.yaml* is the yaml with debug setting. The file can be found [here](https://github.com/microsoft/rai-orchestrator/wiki/OpenTelemetry-Emit-OTel-Logs-Spans-to-OTel-Collector#download-the-otel-collector-container-image-and-run-it-with-debug-configuration)
````
sudo docker run -it -p 4317:4317 -p 4318:4318 -v /home/huliang/testdata/yamls/otel/config.debug.yaml:/etc/otelcol-contrib/config.yaml otel/opentelemetry-collector-contrib:0.93.0
````

Export OTEL_EXPORTER_OTLP_ENDPOINT environment variable, this is necessary because OpenTelemetry SDK is using https://localhost:4318 by default.
````
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
````
Run the stub service in another terminal
````
go run .
````

The service is listening to *http://localhost:12345*, the 3 APIs are /annotate, /api/health and /readiness.
/annotate accpet POST http request and it will output the text in the body with a random number. The response will be written to a Otel Span.
/api/heath and /readiness both return 200 (OK)

Once the service is started and the Otel collector is running, the metrics will be sent to Otel collector and if you make http request to the stub service, spans will also be sent to the collector. You should be able to see them in the console output from the Otel collector.


# Verify OpenTelemetry Exporter/Collector flow with remote public collector in Azure.
Export this endpoint enviroinment variable to the docker container.
````
export OTEL_EXPORTER_OTLP_ENDPOINT=https://ca-otelcol-lgvgvhiuark32.nicefield-824a522d.westus3.azurecontainerapps.io
export AZURE_TENANT_ID = "33e01921-4d64-4f8c-a055-5bdaffd5e33d"
export AZURE_CLIENT_ID = "9c7ae59d-9323-4423-a0da-38ddce774875"
export AZURE_CLIENT_SECRET = "CANNOTSHOW"
````


# Build the docker container.
A Dockerfile is in the same directory. The command to build the docker image is:
````
go build -o otelsender .
docker build -t otel-sender-container .
docker run -p 12345:12345 otel-sender-container
````
