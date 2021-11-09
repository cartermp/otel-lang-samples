open System
open System.Threading.Tasks

open Microsoft.AspNetCore.Builder
open Microsoft.Extensions.Hosting

open OpenTelemetry.Trace
open Honeycomb.OpenTelemetry

let builder = WebApplication.CreateBuilder(Environment.GetCommandLineArgs())
builder.Services.AddHoneycomb(builder.Configuration) |> ignore

let rootHandler (tracer: Tracer) =
    task {
        use span = tracer.StartActiveSpan("sleep span")
        span.SetAttribute("duration_ms", 100) |> ignore

        do! Task.Delay(100)

        return "Hello World!"
    }

let app = builder.Build()
app.MapGet("/", Func<Tracer, Task<string>>(rootHandler)) |> ignore

app.Run()
