from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.exporter.otlp.proto.http.metric_exporter import OTLPMetricExporter
from opentelemetry.sdk.resources import Resource
 
trace.set_tracer_provider(TracerProvider())
tracer = trace.get_tracer(__name__)
 
span_processor = BatchSpanProcessor(OTLPSpanExporter(endpoint="http://127.0.0.1:4317"))
trace.get_tracer_provider().add_span_processor(span_processor)
 
#with tracer.start_as_current_span("HuaLiang"):
#    print("Hello world from OpenTelemetry Python!")


import requests  
import time  

from opentelemetry import metrics
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import (
    ConsoleMetricExporter,
    PeriodicExportingMetricReader,
)

# metric_reader = PeriodicExportingMetricReader(OTLPMetricExporter(endpoint="http://localhost:4318"))
# provider = MeterProvider(metric_readers=[metric_reader],
#             resource=Resource.create({
#             "service.name": "shoppingcart",
#             "service.instance.id": "instance-12",
#         }))

# # Sets the global default meter provider
# metrics.set_meter_provider(provider)

# # Creates a meter from the global meter provider
# meter = metrics.get_meter("huliang-test-meter-name")

# work_counter = meter.create_counter(
#     "work.counter", unit="1", description="Counts the amount of work done"
# )

# work_counter.add(1, {"work.type": "Nothing"})

while True:  
    try:  
        # response = requests.get("http://example.com")
        response = requests.get("http://localhost:12345/annotate")  
        status_code = response.status_code  
        print("Status code:", status_code)  
    except requests.exceptions.RequestException as e:  
        print("Error occurred:", e)

    time.sleep(10)  # Wait for 1 minute before making the next request  