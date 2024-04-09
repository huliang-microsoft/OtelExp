from typing import Iterable
import requests  
import time
import opentelemetry.instrumentation.requests
from opentelemetry.sdk.resources import Resource


from opentelemetry.exporter.otlp.proto.http.metric_exporter import (
    OTLPMetricExporter,
)
from opentelemetry.metrics import (
    CallbackOptions,
    Observation,
    get_meter_provider,
    set_meter_provider,
)
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader

import os  
  
# Get a specific environment variable  
aad_token = os.environ.get("TOKEN")
headers = { "Authorization": "Bearer " + aad_token }
print(headers)

exporter = OTLPMetricExporter(headers=headers)
reader = PeriodicExportingMetricReader(exporter, export_interval_millis=1000)
provider = MeterProvider(metric_readers=[reader],
            resource=Resource.create({
            "service.name": "huliang-raio-python",
            "service.namespace": "huliang-raio-python-ns",
            "service.instance.id": "huliang-local",
        }))
set_meter_provider(provider)


def observable_counter_func(options: CallbackOptions) -> Iterable[Observation]:
    yield Observation(1, {})


def observable_up_down_counter_func(
    options: CallbackOptions,
) -> Iterable[Observation]:
    yield Observation(-10, {})


def observable_gauge_func(options: CallbackOptions) -> Iterable[Observation]:
    yield Observation(9, {})


meter = get_meter_provider().get_meter("getting-started", "0.1.2")



# # Async Counter
# observable_counter = meter.create_observable_counter(
#     "observable_counter",
#     [observable_counter_func],
# )

# # UpDownCounter
# updown_counter = meter.create_up_down_counter("updown_counter")
# updown_counter.add(1)
# updown_counter.add(-5)

# # Async UpDownCounter
# observable_updown_counter = meter.create_observable_up_down_counter(
#     "observable_updown_counter", [observable_up_down_counter_func]
# )

# # Histogram
# histogram = meter.create_histogram("histogram")
# histogram.record(99.9)

# # Async Gauge
# gauge = meter.create_observable_gauge("gauge", [observable_gauge_func])

# while True:  
#     try:  
#         response = requests.get("http://example.com")  
#         status_code = response.status_code  
#         print("Status code:", status_code)  
#     except requests.exceptions.RequestException as e:  
#         print("Error occurred:", e)

#     time.sleep(10)  # Wait for 1 minute before making the next request  

# Instrument the HTTP client
# Create a counter instrument for HTTP requests  
opentelemetry.instrumentation.requests.RequestsInstrumentor().instrument()
i = 0
while i < 10:
    response = requests.get(url="https://www.example.org/")
    print(response.status_code)
    i = i+1
    print(i)
    time.sleep(3)
  
# Manually trigger metric export  
#provider.force_flush()  
  
# Shutdown the exporter  
#provider.shutdown()

# try to update header
exporter._session.headers.update({ "Authorization": "Bearer " + "BABA" })
response = requests.get(url="https://www.example.org/")
print(response.status_code)

# Counter
counter = meter.create_counter("rai.annotate")
counter.add(1)