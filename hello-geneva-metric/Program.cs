using System.Diagnostics.Metrics;
using System.Runtime.InteropServices;
using OpenTelemetry;
using OpenTelemetry.Exporter;
using OpenTelemetry.Metrics;
using OpenTelemetry.Resources;
 
public class Program
{
    private static readonly Meter MyMeter = new("FruitCompany.FruitSales", "1.0");
    private static readonly Counter<long> MyFruitCounter = MyMeter.CreateCounter<long>("FruitsSold");
 
    public static async Task Main()
    {
        var meterProvider = Sdk.CreateMeterProviderBuilder()
            .AddMeter(MyMeter.Name)
            .AddConsoleExporter()
            .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService("otlp-test").AddAttributes(new List<KeyValuePair<string, object>>
                {
                    new KeyValuePair<string, object>("_microsoft_metrics_account", "HuaLiangTest"),
                    new KeyValuePair<string, object>("_microsoft_metrics_namespace", "HuaLiangOTelNameSpace"),
                }))
          .AddOtlpExporter((opt, metricReaderOptions) =>
          {
              opt.Protocol = OtlpExportProtocol.Grpc;
              metricReaderOptions.TemporalityPreference = MetricReaderTemporalityPreference.Delta;
          })
          .Build();
 
        while (true)
        {
            MyFruitCounter.Add(1, new("name", "apple"), new("color", "red"));
            MyFruitCounter.Add(2, new("name", "lemon"), new("color", "yellow"));
            MyFruitCounter.Add(1, new("name", "lemon"), new("color", "yellow"));
            MyFruitCounter.Add(2, new("name", "apple"), new("color", "green"));
            MyFruitCounter.Add(5, new("name", "apple"), new("color", "red"));
            MyFruitCounter.Add(4, new("name", "lemon"), new("color", "yellow"));
            await Task.Delay(1000);
        }
 
        // Dispose meterProvider at the end of the application.
        meterProvider.Dispose();
    }
}