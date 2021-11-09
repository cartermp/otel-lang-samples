using OpenTelemetry.Trace;
using Honeycomb.OpenTelemetry;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddHoneycomb(builder.Configuration);

var app = builder.Build();
app.MapGet("/", async (Tracer tracer) =>
{
    using var span = tracer.StartActiveSpan("sleepy span");
    span.SetAttribute("duration_ms", 100);

    await Task.Delay(100);

    return "Hello World!";
});

app.Run();
